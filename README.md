# 机场接送订单匹配系统

## 1. 项目概述

机场接送订单匹配系统是一个后端服务，旨在高效地为机场接送场景匹配乘客和司机。系统会综合考虑机场位置、时间段、车辆类型和定价规则进行智能匹配。平台通过乘客出价与司机报价之间的差价获得收益。

## 2. 系统架构

系统采用领域驱动设计（DDD）原则，代码结构分为以下模块：
- **API 层**（`api/`）：处理 HTTP 和 gRPC 请求。
- **应用层**（`internal/app/`）：协调业务用例。
- **领域层**（`internal/domain/`）：包含核心业务逻辑和实体（订单、结算、用户、事件总线等）。
- **基础设施层**（`pkg/`）：提供数据库、事件总线、支付、Redis、HTTP 工具等集成。

这种分层结构有助于系统的可维护性和可扩展性。架构示意如下：

```
[API 层] <-> [应用层] <-> [领域层] <-> [基础设施层]
```

## 3. 技术栈

- **语言：** Go 1.22
- **框架：** Gin（HTTP）、Gorm (ORM)
- **数据库：** MySQL、Redis
- **消息队列：** Kafka
- **测试：** Go Test
- **容器化：** Docker Compose

## 4. 运行说明

```bash
# 克隆仓库
$ git clone xxx
$ cd airport-pickup

# 使用 Docker Compose 启动
$ docker compose up

# 或本地运行
$ go run cmd/server/main.go

# db migrate
docker exec -i airport-mysql mysql -uairport -pairport airport < db/migrations/001_init_schema.sql

# other
go mod tidy
go build ./...
go test -v ./...

# docker start
docker-compose pull
docker-compose up -d

# docker restart
docker-compose down
docker-compose up -d
```

## 5. API 文档

#### 1. 创建乘客
- **POST** `/passengers`
- **请求体：**
  ```json
  {
    "name": "Alice"
  }
  ```

#### 2. 创建司机
- **POST** `/drivers`
- **请求体：**
  ```json
  {
    "name": "Bob",
    "rating": 4.9
  }
  ```

#### 3. 创建接送请求
- **POST** `/pickup_requests`
- **请求体：**
  ```json
  {
    "passenger_id": "174b032d1244ea6320a77041c034bd8f",
    "airport_code": "SFO",
    "vehicle_type": "sedan",
    "desired_time": "2025-11-05T10:00:00Z",
    "max_price_per_km": 2.5,
    "prefer_high_rating": true
  }
  ```

#### 4. 创建司机报价
- **POST** `/driver_offers`
- **请求体：**
  ```json
  {
    "driver_id": "0bd803342d1661d5380c833f04929417",
    "airport_code": "SFO",
    "vehicle_type": "sedan",
    "available_from": "2025-11-05T09:00:00Z",
    "available_to": "2025-11-05T12:00:00Z",
    "price_per_km": 2.0
  }
  ```

#### 5. 查询订单
- **GET** `/bookings`

#### 6. 完成订单
- **POST** `/bookings?id=ed6c04d6777b4d782f312519623fdf18`

## 6. 领域模型 / 匹配逻辑

匹配算法流程如下：
1. 筛选出可用时间段与乘客请求重叠的司机。
2. 按车辆类型和（可选）评分过滤司机。
3. 选择报价不高于乘客最高出价且价格最低的司机进行匹配。

**伪代码：**
```
for each passenger_request:
    candidates = 查找司机，满足：
        司机.available_time 与乘客.desired_time 重叠
        且 司机.vehicle_type == 乘客.vehicle_type
        且 司机.price_per_km <= 乘客.max_price_per_km
    if candidates:
        选择报价最低的司机
        创建匹配
```

其他也需要考虑，如：1、取消接口 2、接送请求、司机报价漏匹配重试机制（添加定时任务检索，添加驱动消息）