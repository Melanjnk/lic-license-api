package producer

import (
	"context"
	"errors"

	"github.com/gammazero/workerpool"
	"github.com/golang/mock/gomock"
	"github.com/ozonmp/lic-license-api/internal/mocks"
	"github.com/ozonmp/lic-license-api/internal/model/license"
	"github.com/stretchr/testify/suite"
)

type ProducerTestSuite struct {
	suite.Suite
	mockCtrl       *gomock.Controller
	sender         *mocks.MockLicenseEventSender
	repo           *mocks.MockLicenseEventRepo
	producerCount  uint64
	events         chan license.LicenseEvent
	workerPoolSize int
	workerPool     *workerpool.WorkerPool
	producer       LicenseProducer
}

func (suite *ProducerTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.sender = mocks.NewMockLicenseEventSender(suite.mockCtrl)
	suite.repo = mocks.NewMockLicenseEventRepo(suite.mockCtrl)
	suite.producerCount = uint64(1)
	suite.events = make(chan license.LicenseEvent)
	suite.workerPoolSize = 1
	suite.workerPool = workerpool.New(suite.workerPoolSize)
	suite.producer = NewKafkaLicenseProducer(
		suite.producerCount,
		suite.sender,
		suite.events,
		suite.workerPool,
		suite.repo,
	)
}

func (suite *ProducerTestSuite) TearDownTest() {
	suite.producer.Close()
	suite.workerPool.StopWait()
	close(suite.events)
	suite.mockCtrl.Finish()
}

func (suite *ProducerTestSuite) TestStart() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	suite.producer.Start(ctx)
}

func (suite *ProducerTestSuite) TestEventChanRead() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	suite.producer.Start(ctx)

	// clean if ok
	suite.sender.EXPECT().Send(gomock.Any()).Return(nil).Times(1)
	suite.repo.EXPECT().Remove(gomock.Any()).Times(1)

	// update if error
	suite.sender.EXPECT().Send(gomock.Any()).Return(errors.New("test error")).Times(1)
	suite.repo.EXPECT().Unlock(gomock.Any())

	testEvent := license.LicenseEvent{
		ID:     1,
		Type:   license.Created,
		Status: license.Deferred,
		Entity: nil,
	}
	suite.events <- testEvent
	suite.events <- testEvent
}
