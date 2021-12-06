package api

import (
	"context"
	"errors"
	model "github.com/ozonmp/lic-license-api/internal/model/license"
	"github.com/ozonmp/lic-license-api/internal/pkg/logger"
	pb "github.com/ozonmp/lic-license-api/pkg/lic-license-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

func (a *licenseAPI) DescribeLicenseV1(
	ctx context.Context,
	req *pb.DescribeLicenseV1Request,
) (*pb.DescribeLicenseV1Response, error) {

	if err := req.Validate(); err != nil {
		logger.WarnKV(ctx, "DescribeLicenseV1 - invalid argument", "err", err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	license, err := a.licService.Get(ctx, req.LicenseId)
	if err != nil {
		if errors.Is(err, model.ErrLicenseNotFound) {
			logger.WarnKV(ctx, "License #"+
				strconv.FormatUint(req.GetLicenseId(), 10)+
				" not found")

			totalTemplateNotFound.Inc()
			return nil, status.Error(codes.NotFound, "License not found")
		}
		logger.WarnKV(ctx, "DescribeLicenseV1 -- failed", "err", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	if license == nil {
		logger.DebugKV(ctx, "License #"+
			strconv.FormatUint(req.GetLicenseId(), 10)+
			" not found")
		totalTemplateNotFound.Inc()

		return nil, status.Error(codes.NotFound, "License not found")
	}

	logger.DebugKV(ctx, "DescribeTemplateV1 - success", "err", err)

	return &pb.DescribeLicenseV1Response{
		License: convertServiceToPb(license),
	}, nil
}
