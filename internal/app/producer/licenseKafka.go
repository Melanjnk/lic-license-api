package producer

import (
	"context"
	"github.com/gammazero/workerpool"
	"github.com/ozonmp/lic-license-api/internal/app/repo"
	"github.com/ozonmp/lic-license-api/internal/app/sender"

	//workerpool "github.com/ozonmp/lic-license-api/internal/app/worker_pool"
	"github.com/ozonmp/lic-license-api/internal/model/license"
	"sync"
)

type LicenseProducer interface {
	Start(ctx context.Context)
	Close()
}

type licenseProducer struct {
	workerCount uint64
	//timeout     time.Duration
	sender     sender.LicenseEventSender
	workerPool *workerpool.WorkerPool //WorkerLicPool
	repo       repo.LicenseEventRepo
	events     <-chan license.LicenseEvent
	wg         *sync.WaitGroup
	done       chan bool
}

func NewKafkaLicenseProducer(
	workerCount uint64,
	sender sender.LicenseEventSender,
	events <-chan license.LicenseEvent,
	workerPool *workerpool.WorkerPool,
	repo repo.LicenseEventRepo,
) LicenseProducer {

	wg := &sync.WaitGroup{}
	done := make(chan bool)

	return &licenseProducer{
		workerCount: workerCount,
		sender:      sender,
		events:      events,
		workerPool:  workerPool,
		wg:          wg,
		done:        done,
		repo:        repo,
	}
}

func (p *licenseProducer) produceEvent(event license.LicenseEvent) {
	if err := p.sender.Send(&event); err != nil {
		// update
		p.workerPool.Submit(func() {
			p.repo.Unlock([]uint64{event.ID})
		})
	} else {
		// clean
		p.workerPool.Submit(func() {
			p.repo.Remove([]uint64{event.ID})
		})
	}
}

func (p *licenseProducer) consumeEvents(ctx context.Context) {
	defer p.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-p.events:
			p.produceEvent(event)
		}
	}
}

func (p *licenseProducer) Start(ctx context.Context) {
	//ctx context.Context
	for i := uint64(0); i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.consumeEvents(ctx)
	}
}

func (p *licenseProducer) Close() {
	close(p.done)
	p.wg.Wait()
}
