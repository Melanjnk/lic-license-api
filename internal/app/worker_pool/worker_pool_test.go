package workerpool

import (
	"github.com/ozonmp/lic-license-api/internal/app/repo"
	"github.com/ozonmp/lic-license-api/internal/mocks"
	"github.com/ozonmp/lic-license-api/internal/model"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestWorkerPool(t *testing.T) {
	workersCount := 5
	var repo repo.LicenseEventRepo
	NewWorkerLicPool(workersCount, repo)
}

func TestCreateRetranslatorWorkerLicPool(t *testing.T) {
	workersCount := 5
	var repo repo.LicenseEventRepo

	wpI := NewWorkerLicPool(workersCount, repo)
	_, ok := wpI.(RetranslatorWorkerLicPool)
	assert.True(t, ok)
}

type testSubmitter struct {
	submitIsCalled bool
	task           func()
}

func (s *testSubmitter) Submit(task func()) {
	s.submitIsCalled = true
	s.task()
}

func (s *testSubmitter) StopWait() {}
func TestCleanCreatedType(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockLicenseEventRepo(ctrl)

	event := model.LicenseEvent{
		ID:     1,
		Type:   model.Created,
		Status: model.Processed,
		Entity: &model.License{
			ID: 1,
		},
	}
	repo.EXPECT().Remove([]uint64{event.ID}).Return(nil).Times(1)

	wp := RetranslatorWorkerLicPool{
		repo: repo,
	}
	task := func() {
		wp.clean(event)
	}
	s := &testSubmitter{
		task: task,
	}
	wp.submitter = s

	err := wp.Clean(event)
	assert.Nil(t, err)
	assert.True(t, s.submitIsCalled)
}

func TestCleanNotCreatedType(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockLicenseEventRepo(ctrl)

	checkCleanNotCreatedType(repo, model.LicenseEvent{
		ID:     1,
		Type:   model.Updated,
		Status: model.Processed,
		Entity: &model.License{
			ID: 1,
		},
	}, t)
	checkCleanNotCreatedType(repo, model.LicenseEvent{
		ID:     1,
		Type:   model.Removed,
		Status: model.Processed,
		Entity: &model.License{
			ID: 1,
		},
	}, t)
}

func checkCleanNotCreatedType(repo *mocks.MockLicenseEventRepo, event model.LicenseEvent, t *testing.T) {
	repo.EXPECT().Remove([]uint64{event.ID}).Return(nil).Times(0)

	wp := RetranslatorWorkerLicPool{
		submitter: nil,
		repo:      repo,
	}
	task := func() {
		wp.clean(event)
	}
	s := &testSubmitter{
		task: task,
	}
	wp.submitter = s

	err := wp.Clean(event)
	assert.Nil(t, err)
	assert.True(t, s.submitIsCalled)
}

func TestUpdateCreatedType(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockLicenseEventRepo(ctrl)

	event := model.LicenseEvent{
		ID:     1,
		Type:   model.Created,
		Status: model.Processed,
		Entity: &model.License{
			ID: 1,
		},
	}
	repo.EXPECT().Unlock([]uint64{event.ID}).Return(nil).Times(1)

	wp := RetranslatorWorkerLicPool{
		repo: repo,
	}
	task := func() {
		wp.update(event)
	}
	s := &testSubmitter{
		submitIsCalled: false,
		task:           task,
	}
	wp.submitter = s

	err := wp.Update(event)
	assert.Nil(t, err)
	assert.True(t, s.submitIsCalled)
}

func TestUpdateNotCreatedType(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockLicenseEventRepo(ctrl)

	checkUpdateNotCreatedType(repo, model.LicenseEvent{
		ID:     1,
		Type:   model.Updated,
		Status: model.Processed,
		Entity: &model.License{
			ID: 1,
		},
	}, t)
	checkUpdateNotCreatedType(repo, model.LicenseEvent{
		ID:     1,
		Type:   model.Removed,
		Status: model.Processed,
		Entity: &model.License{
			ID: 1,
		},
	}, t)
}

func checkUpdateNotCreatedType(repo *mocks.MockLicenseEventRepo, event model.LicenseEvent, t *testing.T) {
	repo.EXPECT().Unlock([]uint64{event.ID}).Return(nil).Times(0)

	wp := RetranslatorWorkerLicPool{
		repo: repo,
	}
	task := func() {
		wp.update(event)
	}
	s := &testSubmitter{
		submitIsCalled: false,
		task:           task,
	}
	wp.submitter = s

	err := wp.Update(event)
	assert.Nil(t, err)
	assert.True(t, s.submitIsCalled)
}
