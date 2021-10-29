package retranslator

import (
	"context"
	"github.com/gammazero/workerpool"
	"github.com/ozonmp/lic-license-api/internal/app/consumer"
	"github.com/ozonmp/lic-license-api/internal/app/producer"
	"github.com/ozonmp/lic-license-api/internal/app/repo"
	"github.com/ozonmp/lic-license-api/internal/app/sender"
	"github.com/ozonmp/lic-license-api/internal/model/license"
	"time"
)

type Retranslator interface {
	Start(ctx context.Context)
	Close()
}

type Config struct {
	ChannelSize uint64

	ConsumerCount  uint64
	ConsumeSize    uint64
	ConsumeTimeout time.Duration

	ProducerCount uint64
	WorkerCount   int

	Repo   repo.LicenseEventRepo
	Sender sender.LicenseEventSender
}

type retranslator struct {
	events     chan license.LicenseEvent
	consumer   consumer.LicenseConsumer
	producer   producer.LicenseProducer
	workerPool *workerpool.WorkerPool
	//cancel   context.CancelFunc
}

func NewRetranslator(cfg Config) Retranslator {
	events := make(chan license.LicenseEvent, cfg.ChannelSize)
	workerPool := workerpool.New(cfg.WorkerCount)

	consumer := consumer.NewLicenseDbConsumer(
		cfg.ConsumerCount,
		cfg.ConsumeSize,
		cfg.ConsumeTimeout,
		cfg.Repo,
		events)
	producer := producer.NewKafkaLicenseProducer(
		cfg.ProducerCount,
		cfg.Sender,
		events,
		workerPool,
		cfg.Repo,
	)

	return &retranslator{
		events:     events,
		consumer:   consumer,
		producer:   producer,
		workerPool: workerPool,
	}
}

func (r *retranslator) Start(ctx context.Context) {
	r.producer.Start(ctx)
	r.consumer.Start(ctx)
}

func (r *retranslator) Close() {
	r.consumer.Close()
	r.producer.Close()
}
