package retranslator

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/ozonmp/lic-license-api/internal/mocks"
	"github.com/ozonmp/lic-license-api/internal/model/license"
	"github.com/stretchr/testify/suite"
	"time"
)

type RetranslatorTestSuite struct {
	suite.Suite
	mockCtrl     *gomock.Controller
	repo         *mocks.MockLicenseEventRepo
	sender       *mocks.MockLicenseEventSender
	cfg          Config
	retranslator Retranslator
}

func (suite *RetranslatorTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.repo = mocks.NewMockLicenseEventRepo(suite.mockCtrl)
	suite.sender = mocks.NewMockLicenseEventSender(suite.mockCtrl)
	suite.cfg = Config{
		ChannelSize:    512,
		ConsumerCount:  1,
		ConsumeSize:    10,
		ConsumeTimeout: time.Millisecond * 500,
		ProducerCount:  1,
		WorkerCount:    1,
		Repo:           suite.repo,
		Sender:         suite.sender,
	}
	suite.retranslator = NewRetranslator(suite.cfg)
}

func (suite *RetranslatorTestSuite) TearDownTest() {
	suite.retranslator.Close()
	suite.mockCtrl.Finish()
}

func (suite *RetranslatorTestSuite) TestStart() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	suite.repo.EXPECT().Lock(gomock.Any()).AnyTimes()
	suite.retranslator.Start(ctx)
}

func (suite *RetranslatorTestSuite) TestPubSub() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	lockEvents := []license.LicenseEvent{
		{
			ID:     1,
			Type:   license.Created,
			Status: license.Deferred,
			Entity: nil,
		},
	}
	// clean
	suite.repo.EXPECT().Lock(gomock.Any()).Return(lockEvents, nil).Times(1)
	suite.sender.EXPECT().Send(gomock.Any()).Return(nil).Times(1)
	suite.repo.EXPECT().Remove(gomock.Any()).Times(1)

	// update
	suite.repo.EXPECT().Lock(gomock.Any()).Return(lockEvents, nil).Times(1)
	suite.sender.EXPECT().Send(gomock.Any()).Return(errors.New("test error")).Times(1)
	suite.repo.EXPECT().Unlock(gomock.Any()).Times(1)

	suite.retranslator.Start(ctx)
}
