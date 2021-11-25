package api

import (
	"context"
	"github.com/ozonmp/lic-license-api/internal/pkg/logger"
	pb "github.com/ozonmp/lic-license-api/pkg/lic-license-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *licenseAPI) RemoveLicenseV1(
	ctx context.Context,
	req *pb.RemoveLicenseV1Request,
) (*pb.RemoveLicenseV1Response, error) {
	if err := req.Validate(); err != nil {
		logger.WarnKV(ctx, "RemoveLicenseV1 - invalid argument", "err", err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ok, err := a.licService.Remove(ctx, req.LicenseId)
	if err != nil {
		logger.ErrorKV(ctx, "RemoveLicenseV1 -- failed", "err", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Msg("RemoveLicenseV1 - success")

	return &pb.RemoveLicenseV1Response{
		Found: ok,
	}, nil
}
