package api

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/ozonmp/lic-license-api/internal/repo"

	pb "github.com/ozonmp/lic-license-api/pkg/lic-license-api"
)

var (
	totalTemplateNotFound = promauto.NewCounter(prometheus.CounterOpts{
		Name: "lic_license_api_license_not_found_total",
		Help: "Total number of licenses that were not found",
	})
)

type licenseAPI struct {
	pb.UnimplementedLicLicenseApiServiceServer
	repo repo.Repo
}

// NewLicenseAPI returns api of lic-license-api service
func NewLicenseAPI(r repo.Repo) pb.LicLicenseApiServiceServer {
	return &licenseAPI{repo: r}
}
