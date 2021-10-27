package workerpool

import (
	"github.com/gammazero/workerpool"
	"github.com/ozonmp/lic-license-api/internal/app/repo"
	"github.com/ozonmp/lic-license-api/internal/model/license"
)

type WorkerLicPool interface {
	Clean(licenseEvent license.LicenseEvent) error
	Update(licenseEvent license.LicenseEvent) error
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

func (wp RetranslatorWorkerLicPool) Clean(event license.LicenseEvent) error {
	wp.submitter.Submit(func() {
		wp.clean(event)
	})
	return nil
}

func (wp RetranslatorWorkerLicPool) clean(event license.LicenseEvent) {
	if event.Type == license.Created {
		wp.repo.Remove([]uint64{event.ID})
	}
}

func (wp RetranslatorWorkerLicPool) Update(event license.LicenseEvent) error {
	wp.submitter.Submit(func() {
		wp.update(event)
	})
	return nil
}

func (wp RetranslatorWorkerLicPool) update(event license.LicenseEvent) {
	if event.Type == license.Created {
		wp.repo.Unlock([]uint64{event.ID})
	}
}

func (wp RetranslatorWorkerLicPool) Stop() {
	wp.submitter.StopWait()
}
