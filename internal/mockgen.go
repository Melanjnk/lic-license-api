package internal

//go:generate mockgen -destination=./mocks/repo_license_mock.go -package=mocks github.com/ozonmp/lic-license-api/internal/app/repo LicenseEventRepo
//go:generate mockgen -destination=./mocks/repo_lic_mock.go -package=mocks github.com/ozonmp/lic-license-api/internal/repo Repo
//go:generate mockgen -destination=./mocks/sender_license_mock.go -package=mocks github.com/ozonmp/lic-license-api/internal/app/sender LicenseEventSender
//go:generate mockgen -destination=./mocks/worker_pool_mock.go -package=mocks github.com/ozonmp/lic-license-api/internal/app/worker_pool WorkerLicPool
