package api

import (
	"context"
	model "github.com/ozonmp/lic-license-api/internal/model/license"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/ozonmp/lic-license-api/pkg/lic-license-api"
)

var (
	totalTemplateNotFound = promauto.NewCounter(prometheus.CounterOpts{
		Name: "lic_license_api_license_not_found_total",
		Help: "Total number of licenses that were not found",
	})
)

type licenseService interface {
	Get(tx context.Context, subdomainID uint64) (*model.License, error)
	Add(ctx context.Context, service *model.License) (uint64, error)
	List(ctx context.Context, offset uint64, limit uint64) ([]*model.License, error)
	Remove(ctx context.Context, serviceID uint64) (bool, error)
}

type licenseAPI struct {
	pb.UnimplementedLicLicenseApiServiceServer
	licService licenseService
}

// NewLicenseAPI returns api of lic-license-api service
func NewLicenseAPI(srv licenseService) pb.LicLicenseApiServiceServer {
	return &licenseAPI{licService: srv}
}

func convertServiceToPb(license *model.License) *pb.License {
	return &pb.License{
		LicenseId: license.ID,
		Title:     license.Title,
		CreatedAt: timestamppb.New(license.CreatedAt),
		UpdatedAt: timestamppb.New(license.UpdatedAt),
	}
}
