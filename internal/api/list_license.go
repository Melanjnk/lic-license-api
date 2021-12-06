package api

import (
	"context"
	"github.com/ozonmp/lic-license-api/internal/pkg/logger"
	pb "github.com/ozonmp/lic-license-api/pkg/lic-license-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *licenseAPI) ListLicenseV1(
	ctx context.Context,
	req *pb.ListLicenseV1Request,
) (*pb.ListLicenseV1Response, error) {
	if err := req.Validate(); err != nil {
		logger.WarnKV(ctx, "ListLicenseV1 - invalid argument", "err", err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	licenses, err := a.licService.List(ctx, 0, 10)
	if err != nil {
		logger.ErrorKV(ctx, "ListLicenseV1 -- failed", "err", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.DebugKV(ctx, "ListLicenseV1 - success")

	return &pb.ListLicenseV1Response{
		Licenses: make([]*pb.License, len(licenses)),
	}, nil
}
