package worker

import (
	"context"
	"fmt"
	"github.com/emirpasic/gods/trees/redblacktree"
	evt "github.com/gavin/airport-pickup/internal/domain/eventbus"
	order "github.com/gavin/airport-pickup/internal/domain/order"
	orderentity "github.com/gavin/airport-pickup/internal/domain/order/entity"
	"github.com/gavin/airport-pickup/internal/domain/order/service"
	"github.com/gavin/airport-pickup/pkg/redisstore"
	"github.com/gavin/airport-pickup/pkg/util"
	"sync"
)

// 轻量级有序容器，模拟红黑树必要接口
// 为避免外部依赖，这里使用有序切片实现。

type rbItem interface {
	Key() int64
	Equal(rbItem) bool
}

type rbTree struct {
	tree *redblacktree.Tree // key:int64, value:[]rbItem
}

func int64Comparator(a, b interface{}) int {
	ai := a.(int64)
	bi := b.(int64)
	switch {
	case ai < bi:
		return -1
	case ai > bi:
		return 1
	default:
		return 0
	}
}

func newRbTree() *rbTree {
	return &rbTree{tree: redblacktree.NewWith(int64Comparator)}
}

func (t *rbTree) ReplaceOrInsert(it rbItem) {
	if t == nil {
		return
	}
	key := it.Key()
	if val, found := t.tree.Get(key); found {
		lst := val.([]rbItem)
		lst = append(lst, it)
		t.tree.Put(key, lst)
	} else {
		t.tree.Put(key, []rbItem{it})
	}
}

func (t *rbTree) Delete(it rbItem) {
	if t == nil {
		return
	}
	key := it.Key()
	val, found := t.tree.Get(key)
	if !found {
		return
	}
	lst := val.([]rbItem)
	idx := -1
	for i, item := range lst {
		if item.Equal(it) {
			idx = i
			break
		}
	}
	if idx == -1 {
		return
	}
	lst = append(lst[:idx], lst[idx+1:]...)
	if len(lst) == 0 {
		t.tree.Remove(key)
	} else {
		t.tree.Put(key, lst)
	}
}

// OrderWorkerService 串联 Redis、内存订单簿与领域撮合服务。
// 线程安全：使用全局互斥锁保护内存结构。
type OrderWorkerService struct {
	orderRepo order.OrderRepository
	matching  service.MatchingService
	bus       evt.EventBus
	redis     *redisstore.Client

	mu           sync.RWMutex
	requestBooks map[string]*rbTree // key: airport:vehicle -> requests tree
	offerBooks   map[string]*rbTree // key: airport:vehicle -> offers tree
}

func NewOrderWorkerService(orderRepo order.OrderRepository, matching service.MatchingService, bus evt.EventBus, redis *redisstore.Client) *OrderWorkerService {
	return &OrderWorkerService{
		orderRepo:    orderRepo,
		matching:     matching,
		bus:          bus,
		redis:        redis,
		requestBooks: make(map[string]*rbTree),
		offerBooks:   make(map[string]*rbTree),
	}
}

// —— 内存订单簿条目 ——

type requestItem struct{ v *orderentity.PickupRequest }

func (a requestItem) Key() int64 {
	return a.v.DesiredTime.Unix()
}

func (a requestItem) Equal(b rbItem) bool {
	bb, ok := b.(requestItem)
	if !ok {
		return false
	}
	return a.v.ID == bb.v.ID
}

type offerItem struct{ v *orderentity.DriverOffer }

func (a offerItem) Key() int64 {
	return int64(a.v.PricePerKm * 100)
}

func (a offerItem) Equal(b rbItem) bool {
	bb, ok := b.(offerItem)
	if !ok {
		return false
	}
	return a.v.ID == bb.v.ID
}

func bookKey(airport, vehicle string) string { return fmt.Sprintf("%s:%s", airport, vehicle) }

func (s *OrderWorkerService) getOrCreateTrees(key string) (reqTree, offerTree *rbTree) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if t, ok := s.requestBooks[key]; ok {
		reqTree = t
	} else {
		reqTree = newRbTree()
		s.requestBooks[key] = reqTree
	}
	if t, ok := s.offerBooks[key]; ok {
		offerTree = t
	} else {
		offerTree = newRbTree()
		s.offerBooks[key] = offerTree
	}
	return
}

// 只读获取现有树，避免无谓创建
func (s *OrderWorkerService) getTrees(key string) (reqTree, offerTree *rbTree) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.requestBooks[key], s.offerBooks[key]
}

// OnPickupRequestCreated 处理接机请求创建。
func (s *OrderWorkerService) OnPickupRequestCreated(e evt.PickupRequestCreated) error {
	key := bookKey(e.AirportCode, e.VehicleType)
	// 1. 更新 Redis 请求订单簿
	if s.redis != nil {
		_ = s.redis.AddPickupRequest(context.Background(), e.AirportCode, e.VehicleType, e, e.MaxPricePerKm)
	}
	// 2. 更新内存请求订单簿（红黑树）
	req := &orderentity.PickupRequest{ID: e.RequestID, PassengerID: e.PassengerID, AirportCode: e.AirportCode, VehicleType: e.VehicleType,
		DesiredTime: e.DesiredTime, MaxPricePerKm: e.MaxPricePerKm, PreferHighRating: e.PreferHighRating, Status: e.Status}
	reqTree, offerTree := s.getOrCreateTrees(key)
	s.mu.Lock()
	reqTree.ReplaceOrInsert(requestItem{v: req})
	s.mu.Unlock()
	// 3. 获取内存中的司机报价单进行匹配（只收集可能匹配的报价单）
	candidates := s.collectOffers(offerTree, req)
	if len(candidates) == 0 {
		return nil
	}
	offer, err := s.matching.MatchFromCandidates(req, candidates)
	if err != nil {
		return nil // 未匹配到，保持订单簿中的记录
	}
	// 4. 匹配成功：保存订单、发布事件
	if err := s.onMatched(req, offer); err != nil {
		return err
	}
	// 5. 清除内存中的请求、司机报价订单
	s.removeRequest(reqTree, req)
	s.removeOffer(offerTree, offer)
	return nil
}

// OnDriverOfferCreated 处理司机报价创建。
func (s *OrderWorkerService) OnDriverOfferCreated(e evt.DriverOfferCreated) error {
	key := bookKey(e.AirportCode, e.VehicleType)
	// 1. 更新 Redis 的司机报价订单簿
	if s.redis != nil {
		_ = s.redis.AddDriverOffer(context.Background(), e.AirportCode, e.VehicleType, e, e.PricePerKm)
	}
	// 2. 更新内存司机报价订单簿（红黑树）
	offer := &orderentity.DriverOffer{ID: e.OfferID, DriverID: e.DriverID, AirportCode: e.AirportCode, VehicleType: e.VehicleType,
		AvailableFrom: e.AvailableFrom, AvailableTo: e.AvailableTo, PricePerKm: e.PricePerKm, Rating: e.Rating, Status: e.Status}
	reqTree, offerTree := s.getOrCreateTrees(key)
	s.mu.Lock()
	offerTree.ReplaceOrInsert(offerItem{v: offer})
	s.mu.Unlock()
	// 3. 获取内存中的请求订单进行匹配（只收集可能匹配的请求单）
	requests := s.collectRequests(reqTree, offer)
	if len(requests) == 0 {
		return nil
	}
	// 只用新offer撮合，不再全量遍历所有报价
	for _, req := range requests {
		of, err := s.matching.MatchFromCandidates(req, []*orderentity.DriverOffer{offer})
		if err == nil && of != nil {
			if e.AirportCode == req.AirportCode && e.VehicleType == req.VehicleType {
				if err := s.onMatched(req, of); err != nil {
					return err
				}
				// 5. 清除内存中的请求、司机报价订单
				s.removeRequest(reqTree, req)
				s.removeOffer(offerTree, of)
				break
			}
		}
	}
	return nil
}

// collectOffers 根据请求初步过滤报价单，提升撮合效率
func (s *OrderWorkerService) collectOffers(tree *rbTree, req *orderentity.PickupRequest) []*orderentity.DriverOffer {
	res := make([]*orderentity.DriverOffer, 0)
	s.mu.RLock()
	defer s.mu.RUnlock()
	if tree == nil {
		return res
	}
	// 只遍历价格 <= req.MaxPricePerKm 的报价单
	maxPriceKey := int64(req.MaxPricePerKm * 100)
	it := tree.tree.Iterator()
	for it.Next() {
		k := it.Key().(int64)
		if k > maxPriceKey {
			break // 已超出最大价格，无需再遍历
		}
		lst := it.Value().([]rbItem)
		for _, item := range lst {
			offer := item.(offerItem).v
			// 机场、车型、时间区间初步过滤
			if offer.AirportCode != req.AirportCode {
				continue
			}
			if offer.VehicleType != req.VehicleType {
				continue
			}
			if req.DesiredTime.Before(offer.AvailableFrom) || req.DesiredTime.After(offer.AvailableTo) {
				continue
			}
			res = append(res, offer)
		}
	}
	return res
}

// collectRequests 根据司机报价初步过滤请求单，提升撮合效率
func (s *OrderWorkerService) collectRequests(tree *rbTree, offer *orderentity.DriverOffer) []*orderentity.PickupRequest {
	res := make([]*orderentity.PickupRequest, 0)
	s.mu.RLock()
	defer s.mu.RUnlock()
	if tree == nil || offer == nil {
		return res
	}
	fromKey := offer.AvailableFrom.Unix()
	toKey := offer.AvailableTo.Unix()
	it := tree.tree.Iterator()
	for it.Next() {
		key := it.Key().(int64)
		if key < fromKey {
			continue // 还没到可用时间区间
		}
		if key > toKey {
			break // 超出可用时间区间，直接跳出
		}
		lst := it.Value().([]rbItem)
		for _, item := range lst {
			req := item.(requestItem).v
			// 机场、车型过滤
			if req.AirportCode != offer.AirportCode {
				continue
			}
			if req.VehicleType != offer.VehicleType {
				continue
			}
			// 时间区间已由key过滤
			res = append(res, req)
		}
	}
	return res
}

func (s *OrderWorkerService) removeRequest(tree *rbTree, req *orderentity.PickupRequest) {
	if tree == nil || req == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	tree.Delete(requestItem{v: req})
}

func (s *OrderWorkerService) removeOffer(tree *rbTree, offer *orderentity.DriverOffer) {
	if tree == nil || offer == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	tree.Delete(offerItem{v: offer})
}

// onMatched 保存 Booking、更新请求并发布事件
func (s *OrderWorkerService) onMatched(req *orderentity.PickupRequest, offer *orderentity.DriverOffer) error {
	b := s.matching.CreateBooking(req, offer, util.NewID)
	// 先变更领域对象状态
	if err := req.MarkMatched(); err != nil {
		return err
	}
	if err := offer.MarkMatched(); err != nil {
		return err
	}
	// 用事务保存三对象
	if err := s.orderRepo.UpdateAllInTransaction(b, req, offer); err != nil {
		return err
	}
	// 发送匹配成功消息
	s.bus.Publish(evt.OrderMatched{BookingID: b.ID, RequestID: req.ID, DriverOfferID: offer.ID})
	return nil
}

// OnOrderMatched 订阅回调：清理内存与 Redis 中的请求与司机报价
func (s *OrderWorkerService) OnOrderMatched(e evt.OrderMatched) error {
	// 读取仓库信息以定位键与构造删除对象
	var (
		req   *orderentity.PickupRequest
		offer *orderentity.DriverOffer
	)
	if s.orderRepo != nil && e.RequestID != "" {
		if r, err := s.orderRepo.GetPickupRequestByID(e.RequestID); err == nil && r != nil {
			req = r
		}
	}
	if s.orderRepo != nil && e.DriverOfferID != "" {
		if o, err := s.orderRepo.GetDriverOfferByID(e.DriverOfferID); err == nil && o != nil {
			offer = o
		}
	}
	// 确定 airport/vehicle
	airport, vehicle := "", ""
	if req != nil {
		airport, vehicle = req.AirportCode, req.VehicleType
	} else if offer != nil {
		airport, vehicle = offer.AirportCode, offer.VehicleType
	}
	key := ""
	if airport != "" || vehicle != "" {
		key = bookKey(airport, vehicle)
	}
	// 内存清理（如存在）
	if key != "" {
		reqTree, offerTree := s.getTrees(key)
		if reqTree != nil && req != nil {
			s.removeRequest(reqTree, req)
		}
		if offerTree != nil && offer != nil {
			s.removeOffer(offerTree, offer)
		}
	}
	// Redis 清理（幂等）
	if s.redis != nil && airport != "" && vehicle != "" {
		_ = s.redis.RemovePickupRequest(context.Background(), airport, vehicle, e.RequestID)
		_ = s.redis.RemoveDriverOffer(context.Background(), airport, vehicle, e.DriverOfferID)
	}
	return nil
}
