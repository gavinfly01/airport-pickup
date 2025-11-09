package eventbus

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/IBM/sarama"
	evt "github.com/gavin/airport-pickup/internal/domain/eventbus"
)

// mockSyncProducer 实现 sarama.SyncProducer 接口
// 只记录发送的消息，不实际发送到 Kafka

type mockSyncProducer struct {
	msgs []*sarama.ProducerMessage
	fail bool
}

func (m *mockSyncProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	if m.fail {
		return 0, 0, errors.New("mock send error")
	}
	m.msgs = append(m.msgs, msg)
	return 1, 100, nil
}
func (m *mockSyncProducer) SendMessages(_ []*sarama.ProducerMessage) error { return nil }
func (m *mockSyncProducer) Close() error                                   { return nil }
func (m *mockSyncProducer) TxnStatus() sarama.ProducerTxnStatusFlag {
	return sarama.ProducerTxnStatusFlag(0)
}
func (m *mockSyncProducer) IsTransactional() bool { return false }

// 事务相关方法补齐
func (m *mockSyncProducer) BeginTxn() error {
	if m.fail {
		return errors.New("mock begin txn error")
	}
	return nil
}
func (m *mockSyncProducer) CommitTxn() error {
	if m.fail {
		return errors.New("mock commit txn error")
	}
	return nil
}
func (m *mockSyncProducer) AbortTxn() error {
	if m.fail {
		return errors.New("mock abort txn error")
	}
	return nil
}
func (m *mockSyncProducer) EndTxn() error { return nil }
func (m *mockSyncProducer) AddMessageToTxn(_ *sarama.ConsumerMessage, _ string, _ *string) error {
	return nil
}
func (m *mockSyncProducer) AddOffsetsToTxn(_ map[string][]*sarama.PartitionOffsetMetadata, _ string) error {
	return nil
}

// mockConsumerGroup 实现 sarama.ConsumerGroup 接口
// 只记录 Consume 调用，不实际消费 Kafka 消消息

type mockConsumerGroup struct {
	consumeCalled bool
	ctx           context.Context
}

func (m *mockConsumerGroup) Consume(ctx context.Context, _ []string, _ sarama.ConsumerGroupHandler) error {
	m.consumeCalled = true
	m.ctx = ctx
	<-ctx.Done() // 修复：等待 context 被 cancel，模拟真实 sarama 行为
	return ctx.Err()
}
func (m *mockConsumerGroup) Errors() <-chan error        { return nil }
func (m *mockConsumerGroup) Close() error                { return nil }
func (m *mockConsumerGroup) Pause(_ map[string][]int32)  {}
func (m *mockConsumerGroup) Resume(_ map[string][]int32) {}
func (m *mockConsumerGroup) PauseAll()                   {}
func (m *mockConsumerGroup) ResumeAll()                  {}

func TestKafkaEventBus_Publish(t *testing.T) {
	prod := &mockSyncProducer{}
	group := &mockConsumerGroup{}
	bus := &KafkaEventBus{
		producer: prod,
		group:    group,
		topic:    "test-topic",
		handlers: make(map[string][]func(evt.Event)),
		ctx:      context.Background(),
		cancel:   func() {},
		closed:   make(chan struct{}),
	}
	e := evt.OrderMatched{BookingID: "bkid", RequestID: "rid", DriverOfferID: "doid"}
	bus.Publish(e)
	if len(prod.msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(prod.msgs))
	}
	msg := prod.msgs[0]
	keyBytes, _ := msg.Key.Encode()
	if string(keyBytes) != e.Name() {
		t.Errorf("expected key %s, got %s", e.Name(), string(keyBytes))
	}
	if msg.Topic != "test-topic" {
		t.Errorf("expected topic 'test-topic', got %s", msg.Topic)
	}
	if len(msg.Headers) == 0 || string(msg.Headers[0].Key) != "event-name" || string(msg.Headers[0].Value) != e.Name() {
		t.Errorf("expected header 'event-name' with value %s, got %+v", e.Name(), msg.Headers)
	}
	valBytes, _ := msg.Value.Encode()
	var got evt.OrderMatched
	if err := json.Unmarshal(valBytes, &got); err != nil {
		t.Errorf("unmarshal error: %v", err)
	}
	if got.BookingID != "bkid" || got.RequestID != "rid" || got.DriverOfferID != "doid" {
		t.Errorf("event fields not match: %+v", got)
	}
}

func TestKafkaEventBus_Subscribe_Start_Close(t *testing.T) {
	prod := &mockSyncProducer{}
	group := &mockConsumerGroup{}
	ctx, cancel := context.WithCancel(context.Background()) // 修复：使用可取消的 context
	bus := &KafkaEventBus{
		producer: prod,
		group:    group,
		topic:    "test-topic",
		handlers: make(map[string][]func(evt.Event)),
		ctx:      ctx,
		cancel:   cancel, // 修复：赋值 cancel 方法
		closed:   make(chan struct{}),
	}
	bus.Subscribe("OrderMatched", func(e evt.Event) {}) // Subscribe 已启动消费循环
	bus.Start()                                         // 再次调用不会重复启动

	// 等待 goroutine 设置 consumeCalled，避免竞态
	deadline := time.Now().Add(200 * time.Millisecond)
	for !group.consumeCalled && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	if !group.consumeCalled {
		t.Errorf("expected consume called, got false (after wait)")
	}
	if err := bus.Close(); err != nil {
		t.Errorf("close error: %v", err)
	}
}

func TestDecodeEvent(t *testing.T) {
	e := evt.OrderMatched{BookingID: "bkid", RequestID: "rid", DriverOfferID: "doid"}
	b, _ := json.Marshal(e)
	res := decodeEvent(evt.EventOrderMatched, b)
	om, ok := res.(evt.OrderMatched)
	if !ok {
		t.Fatalf("expected OrderMatched type, got %T", res)
	}
	if om.BookingID != "bkid" || om.RequestID != "rid" || om.DriverOfferID != "doid" {
		t.Errorf("event fields not match: %+v", om)
	}

	// 测试未知事件名
	res2 := decodeEvent("UnknownEvent", []byte(`{"foo":"bar"}`))
	if res2.Name() != "UnknownEvent" {
		t.Errorf("expected rawEvent name 'UnknownEvent', got %s", res2.Name())
	}
}
