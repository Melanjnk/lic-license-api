package consumer

import (
	"context"
	"github.com/ozonmp/lic-license-api/internal/model/license"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ozonmp/lic-license-api/internal/mocks"
	"github.com/stretchr/testify/suite"
)

type ConsumerTestSuite struct {
	suite.Suite
	repo          *mocks.MockLicenseEventRepo
	mockCtrl      *gomock.Controller
	events        chan license.LicenseEvent
	consumerCount uint64
	batchSize     uint64
	timeout       time.Duration
	consumer      LicenseConsumer
}

func (suite *ConsumerTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.repo = mocks.NewMockLicenseEventRepo(suite.mockCtrl)
	suite.events = make(chan license.LicenseEvent)
	suite.consumerCount = uint64(1)
	suite.batchSize = uint64(10)
	suite.timeout = time.Millisecond
	suite.consumer = NewLicenseDbConsumer(
		suite.consumerCount,
		suite.batchSize,
		suite.timeout,
		suite.repo,
		suite.events,
	)
}

func (suite *ConsumerTestSuite) TearDownTest() {
	suite.consumer.Close()
	close(suite.events)
	suite.mockCtrl.Finish()
}

func (suite *ConsumerTestSuite) TestStart() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	suite.repo.EXPECT().Add(gomock.Any()).AnyTimes()
	suite.consumer.Start(ctx)
}
