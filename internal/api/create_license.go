package api

import (
	"context"
	"github.com/ozonmp/lic-license-api/internal/model"
	pb "github.com/ozonmp/lic-license-api/pkg/lic-license-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *licenseAPI) CreateLicenseV1(
	ctx context.Context,
	req *pb.CreateLicenseV1Request,
) (*pb.CreateLicenseV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("CreateLicenseV1Request - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res := model.License{
		ID:    1,
		Title: "Lic 1",
	}
	licenseID, err := a.repo.CreateLicense(ctx, &res)
	if err != nil {
		log.Err(err).Msg("CreateLicenseV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Msg("CreateLicenseV1 - success")

	return &pb.CreateLicenseV1Response{
		LicenseId: licenseID,
	}, nil
}
