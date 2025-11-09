package worker

import (
	evt "github.com/gavin/airport-pickup/internal/domain/eventbus"
	"log"
)

type SettlementOrchestrator interface {
	OnOrderCompleted(bookingID string) error
}

type Consumer struct {
	worker *OrderWorkerService
}

func NewEventConsumer(bus evt.EventBus, settlement SettlementOrchestrator, worker *OrderWorkerService) *Consumer {
	c := &Consumer{worker: worker}
	// 结算编排：订单完成
	bus.Subscribe(evt.EventOrderCompleted, func(e evt.Event) {
		log.Printf("[event_consumer] handle event: %s, value: %+v", e.Name(), e)
		if oc, ok := e.(evt.OrderCompleted); ok {
			err := settlement.OnOrderCompleted(oc.BookingID)
			if err != nil {
				log.Printf("[event_consumer] OnOrderCompleted failed: %v", err)
			} else {
				log.Printf("[event_consumer] OnOrderCompleted success, bookingID=%s", oc.BookingID)
			}
		} else {
			log.Printf("[event_consumer] event type assertion failed: %T", e)
		}
	})
	// 撮合：接机请求创建
	bus.Subscribe(evt.EventPickupRequestCreated, func(e evt.Event) {
		log.Printf("[event_consumer] handle event: %s, value: %+v", e.Name(), e)
		if ev, ok := e.(evt.PickupRequestCreated); ok {
			if c.worker != nil {
				err := c.worker.OnPickupRequestCreated(ev)
				if err != nil {
					log.Printf("[event_consumer] OnPickupRequestCreated failed: %v", err)
				} else {
					log.Printf("[event_consumer] OnPickupRequestCreated success, requestID=%s", ev.RequestID)
				}
			}
		} else {
			log.Printf("[event_consumer] event type assertion failed: %T", e)
		}
	})
	// 撮合：司机报价创建
	bus.Subscribe(evt.EventDriverOfferCreated, func(e evt.Event) {
		log.Printf("[event_consumer] handle event: %s, value: %+v", e.Name(), e)
		if ev, ok := e.(evt.DriverOfferCreated); ok {
			if c.worker != nil {
				err := c.worker.OnDriverOfferCreated(ev)
				if err != nil {
					log.Printf("[event_consumer] OnDriverOfferCreated failed: %v", err)
				} else {
					log.Printf("[event_consumer] OnDriverOfferCreated success, offerID=%s", ev.OfferID)
				}
			}
		} else {
			log.Printf("[event_consumer] event type assertion failed: %T", e)
		}
	})
	// 撮合：订单匹配完成，清理内存与 Redis
	bus.Subscribe(evt.EventOrderMatched, func(e evt.Event) {
		log.Printf("[event_consumer] handle event: %s, value: %+v", e.Name(), e)
		if ev, ok := e.(evt.OrderMatched); ok {
			if c.worker != nil {
				err := c.worker.OnOrderMatched(ev)
				if err != nil {
					log.Printf("[event_consumer] OnOrderMatched failed: %v", err)
				} else {
					log.Printf("[event_consumer] OnOrderMatched success, bookingID=%s", ev.BookingID)
				}
			}
		} else {
			log.Printf("[event_consumer] event type assertion failed: %T", e)
		}
	})
	return c
}
