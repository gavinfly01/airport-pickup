package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"

	httpapi "github.com/gavin/airport-pickup/api/http"
	"github.com/gavin/airport-pickup/internal/app"
	"github.com/gavin/airport-pickup/internal/config"
	evt "github.com/gavin/airport-pickup/internal/domain/eventbus"
	"github.com/gavin/airport-pickup/internal/domain/order"
	"github.com/gavin/airport-pickup/internal/domain/order/service"
	"github.com/gavin/airport-pickup/internal/domain/settlement"
	"github.com/gavin/airport-pickup/internal/domain/user"
	"github.com/gavin/airport-pickup/internal/worker"
	kbus "github.com/gavin/airport-pickup/pkg/eventbus"
	"github.com/gavin/airport-pickup/pkg/payments"
	"github.com/gavin/airport-pickup/pkg/redisstore"
	mysqlrepo "github.com/gavin/airport-pickup/pkg/repository/mysql"
)

func buildRepos(cfg *config.Config) (user.PassengerRepository, user.DriverRepository, order.OrderRepository, settlement.SettlementRepository, error) {
	if dsn := cfg.Database.DSN; dsn != "" {
		db, err := mysqlrepo.NewDB(dsn)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		if cfg.Database.AutoMigrate {
			if err := mysqlrepo.AutoMigrate(db); err != nil {
				return nil, nil, nil, nil, err
			}
		}
		log.Println("using MySQL repositories")
		return mysqlrepo.NewPassengerRepository(db), mysqlrepo.NewDriverRepository(db), mysqlrepo.NewOrderRepository(db), mysqlrepo.NewSettlementRepository(db), nil
	}
	return nil, nil, nil, nil, errors.New("connect mysql failed: empty DSN")
}

func main() {
	// 读取配置文件路径
	cfgPath := flag.String("config", "config/dev.yaml", "path to config yaml")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	// Event bus: Kafka 优先，其次回退至内存
	var bus evt.EventBus

	brokers := cfg.Kafka.Brokers
	topic := cfg.Kafka.Topic
	groupID := cfg.Kafka.GroupID

	var kafkaBus *kbus.KafkaEventBus
	if len(brokers) > 0 && topic != "" && groupID != "" {
		kb, err := kbus.NewKafkaEventBus(brokers, topic, groupID)
		if err != nil {
			log.Printf("init KafkaEventBus failed, fallback to memory bus: %v", err)
		} else {
			kafkaBus = kb
			bus = kb
			log.Printf("using Kafka event bus, brokers=%v, topic=%s, group=%s", brokers, topic, groupID)
		}
	}
	if bus == nil {
		log.Fatalf("event bus init failed: no available event bus (kafka and memorybus both unavailable)")
	}

	// Redis 初始化
	var rds *redisstore.Client
	{
		opt := redisstore.Options{Addr: cfg.Redis.Addr, Password: cfg.Redis.Password, DB: cfg.Redis.DB}
		rds = redisstore.New(opt)
		if err := rds.Ping(context.Background()); err != nil {
			log.Printf("redis ping failed (will continue without redis): %v", err)
			rds = nil
		} else {
			log.Printf("redis connected: addr=%s db=%d", opt.Addr, opt.DB)
		}
	}

	// Payment client
	pay := payments.NewWalletClient()

	// Repositories
	passRepo, driverRepo, orderRepo, settlementRepo, err := buildRepos(cfg)
	if err != nil {
		log.Fatalf("repository init failed: %v", err)
	}

	// Domain services
	matching := service.NewMatchingService(orderRepo, driverRepo)

	// App services
	orderApp := app.NewOrderAppService(orderRepo, passRepo, driverRepo, matching, bus)
	settlementApp := app.NewSettlementAppService(settlementRepo, orderRepo, pay, bus)

	// Worker service for matching
	orderWorker := worker.NewOrderWorkerService(orderRepo, matching, bus, rds)

	// Workers: subscribe to events（首次订阅将启动 Kafka 消费循环）
	_ = worker.NewEventConsumer(bus, settlementApp, orderWorker)

	// 优雅关闭（Kafka 模式）
	if kafkaBus != nil {
		defer func() {
			if err := kafkaBus.Close(); err != nil {
				log.Printf("close Kafka bus error: %v", err)
			}
		}()
	}

	// HTTP router
	r := httpapi.NewRouter(orderApp, settlementApp)

	log.Printf("server listening on %s", cfg.Server.Addr)
	if err := http.ListenAndServe(cfg.Server.Addr, r); err != nil {
		log.Fatal(err)
	}
}
