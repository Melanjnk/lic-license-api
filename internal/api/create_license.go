package api

import (
	"context"
	model "github.com/ozonmp/lic-license-api/internal/model/license"
	"github.com/ozonmp/lic-license-api/internal/pkg/logger"
	pb "github.com/ozonmp/lic-license-api/pkg/lic-license-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *licenseAPI) CreateLicenseV1(
	ctx context.Context,
	req *pb.CreateLicenseV1Request,
) (*pb.CreateLicenseV1Response, error) {

	if err := req.Validate(); err != nil {
		logger.WarnKV(ctx, "CreateLicenseV1Request - invalid argument", "err", err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res := model.License{
		ID:    1,
		Title: "Lic 1",
	}
	licenseID, err := a.licService.Add(ctx, &res)
	if err != nil {
		logger.ErrorKV(ctx, "CreateLicenseV1 -- failed", "err", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.DebugKV(ctx, "CreateLicenseV1 - success")

	return &pb.CreateLicenseV1Response{
		LicenseId: licenseID,
	}, nil
}
