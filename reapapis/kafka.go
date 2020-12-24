package reapapis

import (
	"context"
	"github.com/ethereum/go-ethereum/log"
	"github.com/segmentio/kafka-go"
	"os"
	"time"
)

const (
	FromBegining = -2
	ToEnd = -1
)

type kafkaClient struct {
	config   *kafka.ReaderConfig
	consumer *kafka.Reader
	//dataChannel chan kafka.Message
	callback func(message kafka.Message) bool
}

// 변수 초기화.
func NewKafkaClient(brokers []string, topic string, procCallBack func(message kafka.Message) bool) *kafkaClient {
	return &kafkaClient{
		// MaxWait값은, MinBytes 값보다 적게 채워지고, 수신이 안되나, MaxWait 지나면 자동수신
		config: &kafka.ReaderConfig{
			Brokers:       brokers,
			GroupID:       "",
			Topic:         topic,
			Partition:     0,
			Dialer:        nil,
			QueueCapacity: 0,
			MinBytes:      1,
			MaxBytes:      10e6,
			MaxWait:       10 * time.Millisecond,
		},
		consumer: nil,
		callback: procCallBack,
	}
}

// kafka broker 에 연결시도, 자동 firstoffset으로 셋팅됨
func (k *kafkaClient) connect() {
	k.consumer = kafka.NewReader(*k.config)
}

// ReadOffset 이 값 기준으로, kafka 서버의 Offset 시작위치부터
// ReadOffset 값이 ToEnd이면 마지막 메시지부터
// ReadOffset 값이 FromBeginning이면 처음 메시지부터
// kafka로부터 메시지를 계속 수신함. blocking mode
// 등록된 callback 함수가, 처리 후, 리턴값이 true 인 경우만, kafka의 offset commit
func (k *kafkaClient) ReadBackground(ctx context.Context, ReadOffset int64) {
	k.connect()

	if err := k.consumer.SetOffset(ReadOffset); err != nil {
		k.consumer.Close()
		log.Error("kafka, failed to set message index", "kafka", err)
		os.Exit(-1)
	}
	go func() {
		for {
			m, err := k.consumer.FetchMessage(ctx)
			if err != nil {
				log.Error("FetchMessage error :", "kafka", err)
				os.Exit(-1)
			}
			if k.callback(m) == true {
				k.commit(ctx, m)
				log.Debug("Commit", "kafka Message index", m.Offset)
			} else {
				log.Error("Not Commit", "kafka Message index", m.Offset)
			}
		}
	}()
}

//kafka에 message를 수신 후, commit 한다.
func (k *kafkaClient) commit(ctx context.Context, message ...kafka.Message) error {
	return k.consumer.CommitMessages(ctx, message...)
}
