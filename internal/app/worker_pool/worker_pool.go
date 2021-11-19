package workerpool

import (
	"github.com/gammazero/workerpool"
	"github.com/ozonmp/lic-license-api/internal/app/repo"
	"github.com/ozonmp/lic-license-api/internal/model"
)

type WorkerLicPool interface {
	Clean(licenseEvent model.LicenseEvent) error
	Update(licenseEvent model.LicenseEvent) error
	Stop()
}

type submitter interface {
	Submit(task func())
	StopWait()
}

func NewWorkerLicPool(workerCount int, repo repo.LicenseEventRepo) WorkerLicPool {
	return RetranslatorWorkerLicPool{
		submitter: workerpool.New(workerCount),
		repo:      repo,
	}
}

type RetranslatorWorkerLicPool struct {
	submitter submitter
	repo      repo.LicenseEventRepo
}

func (wp RetranslatorWorkerLicPool) Clean(event model.LicenseEvent) error {
	wp.submitter.Submit(func() {
		wp.clean(event)
	})
	return nil
}

func (wp RetranslatorWorkerLicPool) clean(event model.LicenseEvent) {
	if event.Type == model.Created {
		wp.repo.Remove([]uint64{event.ID})
	}
}

func (wp RetranslatorWorkerLicPool) Update(event model.LicenseEvent) error {
	wp.submitter.Submit(func() {
		wp.update(event)
	})
	return nil
}

func (wp RetranslatorWorkerLicPool) update(event model.LicenseEvent) {
	if event.Type == model.Created {
		wp.repo.Unlock([]uint64{event.ID})
	}
}

func (wp RetranslatorWorkerLicPool) Stop() {
	wp.submitter.StopWait()
}
