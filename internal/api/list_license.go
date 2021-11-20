package api

import (
	"context"
	pb "github.com/ozonmp/lic-license-api/pkg/lic-license-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *licenseAPI) ListLicenseV1(
	ctx context.Context,
	req *pb.ListLicenseV1Request,
) (*pb.ListLicenseV1Response, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("DescribeLicenseV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	licenses, err := a.licService.ListLicense(ctx, 0, 10)
	if err != nil {
		log.Error().Err(err).Msg("DescribeLicenseV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Msg("DescribeLicenseV1 - success")

	return &pb.ListLicenseV1Response{
		Licenses: make([]*pb.License, len(licenses)),
	}, nil
}
