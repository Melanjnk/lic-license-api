package retranslator

import (
	"context"
	"github.com/ozonmp/lic-license-api/internal/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestStart(t *testing.T) {

	ctrl := gomock.NewController(t)
	repo := mocks.MockLicenseEventRepo(ctrl)
	sender := mocks.NewMockLicenseEventSender(ctrl)

	repo.EXPECT().Lock(gomock.Any()).AnyTimes()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := Config{
		ChannelSize:    512,
		ConsumerCount:  2,
		ConsumeSize:    10,
		ConsumeTimeout: 10 * time.Second,
		ProducerCount:  2,
		WorkerCount:    2,
		Repo:           repo,
		Sender:         sender,
	}

	retranslator := NewRetranslator(cfg)
	retranslator.Start(ctx)
	retranslator.Close()
}
