package internal

//go:generate mockgen -destination=./mocks/repo_license_mock.go -package=mocks github.com/ozonmp/lic-license-api/internal/repo/license/service LicenseRepo
//go:generate mockgen -destination=./mocks/event_repo_mock.go -package=mocks github.com/ozonmp/lic-license-api/internal/repo/license/service EventRepo
//go:generate mockgen -destination=./mocks/repo_lic_mock.go -package=mocks github.com/ozonmp/lic-license-api/internal/repo/license/service Repo
//go:generate mockgen -destination=./mocks/transactional_session_mock.go -package=mocks github.com/ozonmp/lic-license-api/internal/repo TransactionalSession
//go:generate mockgen -destination=./mocks/sender_license_mock.go -package=mocks github.com/ozonmp/lic-license-api/internal/app/sender LicenseEventSender
//go:generate mockgen -destination=./mocks/worker_pool_mock.go -package=mocks github.com/ozonmp/lic-license-api/internal/app/worker_pool WorkerLicPool
