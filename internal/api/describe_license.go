package api

import (
	"context"
	"errors"
	model "github.com/ozonmp/lic-license-api/internal/model/license"
	pb "github.com/ozonmp/lic-license-api/pkg/lic-license-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *licenseAPI) DescribeLicenseV1(
	ctx context.Context,
	req *pb.DescribeLicenseV1Request,
) (*pb.DescribeLicenseV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("DescribeLicenseV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	license, err := a.licService.DescribeLicense(ctx, req.LicenseId)
	if err != nil {
		if errors.Is(err, model.ErrLicenseNotFound) {
			log.Debug().Uint64("licenseId", req.GetLicenseId()).Msg("license not found")
			totalTemplateNotFound.Inc()
			return nil, status.Error(codes.NotFound, "license not found")
		}
		log.Error().Err(err).Msg("DescribeLicenseV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if license == nil {
		log.Debug().Uint64("licenseId", req.LicenseId).Msg("license not found")
		totalTemplateNotFound.Inc()

		return nil, status.Error(codes.NotFound, "license not found")
	}

	log.Debug().Msg("DescribeTemplateV1 - success")

	return &pb.DescribeLicenseV1Response{
		License: convertServiceToPb(license),
	}, nil
}
