# JCQ 消费者

独立的 JCQ 消费进程，默认从 `configs/config.yaml` 的 `jcq` 节读取配置。

## 配置
在 `server/configs/config.yaml` 增加或修改：

```yaml
jcq:
  access_key: "your-ak"
  secret_key: "your-sk"
  topic: "your-topic"
  consumer_group_id: "your-consumer-group"
  http_url: "https://jcq-shared-004.cn-north-1.jdcloud.com"
  auto_ack: false
  poll_size: 5
  poll_interval_seconds: 1
```

> 请使用你自己的 AK/SK/Topic/ConsumerGroup/HTTP 接入点。

## 构建与运行

```bash
cd server
go build ./cmd/jcqconsumer
./jcqconsumer
```

## 业务处理

在 `main.go` 的 `handleMessage` 中接入你的业务逻辑（消息解析、订单状态更新、对账等）。

