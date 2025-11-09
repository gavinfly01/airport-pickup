package eventbus

import (
	"context"
	"encoding/json"
	"log"
	"runtime/debug"
	"sync"
	"time"

	"github.com/IBM/sarama"
	evt "github.com/gavin/airport-pickup/internal/domain/eventbus"
)

const headerEventName = "event-name"

// KafkaEventBus 基于 Kafka 的事件总线实现。
type KafkaEventBus struct {
	producer sarama.SyncProducer
	group    sarama.ConsumerGroup
	topic    string

	mu        sync.RWMutex
	handlers  map[string][]func(evt.Event)
	ctx       context.Context
	cancel    context.CancelFunc
	startOnce sync.Once
	closed    chan struct{}
}

// NewKafkaEventBus 初始化生产者和消费组。
func NewKafkaEventBus(brokers []string, topic, groupID string) (*KafkaEventBus, error) {
	version, err := sarama.ParseKafkaVersion("2.8.0")
	if err != nil {
		return nil, err
	}

	pcfg := sarama.NewConfig()
	pcfg.Version = version
	pcfg.Producer.RequiredAcks = sarama.WaitForAll
	pcfg.Producer.Return.Successes = true
	pcfg.Producer.Idempotent = true
	pcfg.Net.MaxOpenRequests = 1
	pcfg.Producer.Retry.Max = 6
	pcfg.Producer.Retry.Backoff = 100 * time.Millisecond

	ccfg := sarama.NewConfig()
	ccfg.Version = version
	ccfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	ccfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	ccfg.Consumer.Return.Errors = true

	prod, err := sarama.NewSyncProducer(brokers, pcfg)
	if err != nil {
		return nil, err
	}
	group, err := sarama.NewConsumerGroup(brokers, groupID, ccfg)
	if err != nil {
		_ = prod.Close()
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &KafkaEventBus{
		producer: prod,
		group:    group,
		topic:    topic,
		handlers: make(map[string][]func(evt.Event)),
		ctx:      ctx,
		cancel:   cancel,
		closed:   make(chan struct{}),
	}, nil
}

// Publish: 将事件名写入 header，并序列化事件为 JSON。
func (k *KafkaEventBus) Publish(e evt.Event) {
	log.Printf("[eventbus] publish event: %s, value: %+v", e.Name(), e)
	b, err := json.Marshal(e)
	if err != nil {
		log.Printf("[eventbus] marshal event %s error: %v", e.Name(), err)
		return
	}
	msg := &sarama.ProducerMessage{
		Topic: k.topic,
		Key:   sarama.StringEncoder(e.Name()),
		Value: sarama.ByteEncoder(b),
		Headers: []sarama.RecordHeader{
			{Key: []byte(headerEventName), Value: []byte(e.Name())},
		},
	}
	partition, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		log.Printf("[eventbus] send event %s error: %v", e.Name(), err)
	} else {
		log.Printf("[eventbus] event %s sent successfully, partition=%d, offset=%d", e.Name(), partition, offset)
	}
}

// Subscribe: 注册处理器，并在首次调用时启动消费循环。
func (k *KafkaEventBus) Subscribe(eventName string, handler func(evt.Event)) {
	log.Printf("[eventbus] subscribe event: %s", eventName)
	k.mu.Lock()
	k.handlers[eventName] = append(k.handlers[eventName], handler)
	k.mu.Unlock()

	k.startOnce.Do(func() { go k.consumeLoop() })
}

// Start: 显式启动消费循环（可选）。
func (k *KafkaEventBus) Start() { k.startOnce.Do(func() { go k.consumeLoop() }) }

// Close: 停止消费并关闭连接。
func (k *KafkaEventBus) Close() error {
	k.cancel()
	<-k.closed
	var first error
	if err := k.group.Close(); err != nil {
		first = err
	}
	if err := k.producer.Close(); err != nil && first == nil {
		first = err
	}
	return first
}

func (k *KafkaEventBus) consumeLoop() {
	defer close(k.closed)
	handler := &cgHandler{bus: k}
	for {
		if err := k.group.Consume(k.ctx, []string{k.topic}, handler); err != nil {
			if k.ctx.Err() != nil {
				return
			}
			log.Printf("[kafka] consume error: %v", err)
			time.Sleep(time.Second)
		}
		if k.ctx.Err() != nil {
			return
		}
	}
}

type cgHandler struct{ bus *KafkaEventBus }

func (h *cgHandler) Setup(s sarama.ConsumerGroupSession) error   { return nil }
func (h *cgHandler) Cleanup(s sarama.ConsumerGroupSession) error { return nil }

func (h *cgHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		name := ""
		for _, hd := range msg.Headers {
			if string(hd.Key) == headerEventName {
				name = string(hd.Value)
				break
			}
		}
		if name == "" { // fallback to key
			if msg.Key != nil {
				name = string(msg.Key)
			}
		}
		if name == "" {
			sess.MarkMessage(msg, "")
			continue
		}

		log.Printf("[eventbus] received event: %s, partition=%d, offset=%d, value=%s", name, msg.Partition, msg.Offset, string(msg.Value))

		ev := decodeEvent(name, msg.Value)

		h.bus.mu.RLock()
		handlers := append([]func(evt.Event){}, h.bus.handlers[name]...)
		h.bus.mu.RUnlock()

		for _, cb := range handlers {
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("[kafka] handler panic: %v\n%s", r, debug.Stack())
					}
				}()
				cb(ev)
			}()
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}

// decodeEvent 根据事件名反序列化为具体领域事件类型。
func decodeEvent(name string, payload []byte) evt.Event {
	switch name {
	case evt.EventOrderMatched:
		var v evt.OrderMatched
		if err := json.Unmarshal(payload, &v); err == nil {
			return v
		}
	case evt.EventOrderCompleted:
		var v evt.OrderCompleted
		if err := json.Unmarshal(payload, &v); err == nil {
			return v
		}
	case evt.EventPaymentSucceeded:
		var v evt.PaymentSucceeded
		if err := json.Unmarshal(payload, &v); err == nil {
			return v
		}
	case evt.EventSettlementCreated:
		var v evt.SettlementCreated
		if err := json.Unmarshal(payload, &v); err == nil {
			return v
		}
	case evt.EventRevenueUpdated:
		var v evt.RevenueUpdated
		if err := json.Unmarshal(payload, &v); err == nil {
			return v
		}
	case evt.EventPickupRequestCreated:
		var v evt.PickupRequestCreated
		if err := json.Unmarshal(payload, &v); err == nil {
			return v
		}
	case evt.EventDriverOfferCreated:
		var v evt.DriverOfferCreated
		if err := json.Unmarshal(payload, &v); err == nil {
			return v
		}
	}
	// 默认返回一个仅带名称的事件，避免丢失。
	return rawEvent(name)
}

type rawEvent string

func (r rawEvent) Name() string { return string(r) }
