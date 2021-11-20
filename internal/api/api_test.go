package api

import (
	"context"
	"github.com/golang/mock/gomock"
	mocks "github.com/ozonmp/lic-license-api/internal/mocks"
	"github.com/ozonmp/lic-license-api/internal/service/license"
	pb "github.com/ozonmp/lic-license-api/pkg/lic-license-api"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

func dialer(t *testing.T) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	ctrl := gomock.NewController(t)
	repo := mocks.NewMockRepo(ctrl)
	eventRepo := mocks.NewMockLicenseEventRepo(ctrl)
	tsx := mocks.NewMockTransactionalSession(ctrl)
	pb.RegisterLicLicenseApiServiceServer(server, NewLicenseAPI(license.NewLicenseService(repo, eventRepo, tsx)))

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(ctx context.Context, s string) (net.Conn, error) {
		return listener.Dial()
	}
}

func prepareClient(ctx context.Context, t *testing.T) (client pb.LicLicenseApiServiceClient, closeClient func()) {
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(t)))
	if err != nil {
		log.Fatal(err)
	}
	closeCl := func() {
		err := conn.Close()
		if err != nil {
			log.Panicln(err)
		}
	}
	return pb.NewLicLicenseApiServiceClient(conn), closeCl
}

func TestLicenseAPI_CreateLicenseV1(t *testing.T) {
	ctx := context.Background()

	client, closeCl := prepareClient(ctx, t)
	defer closeCl()

	requests := []*pb.CreateLicenseV1Request{
		{},
		{Title: ""},
	}
	for _, request := range requests {
		response, err := client.CreateLicenseV1(ctx, request)

		assert.Nil(t, response)
		assert.NotNil(t, err)

		er, _ := status.FromError(err)

		assert.Equal(t, codes.InvalidArgument, er.Code())
		assert.Equal(t, "invalid CreateLicenseV1Request.LicenseId: value must be greater than 0", er.Message())
	}
}

func TestLicenseAPI_DescribeLicenseV1(t *testing.T) {
	ctx := context.Background()

	client, closeCl := prepareClient(ctx, t)
	defer closeCl()

	requests := []*pb.DescribeLicenseV1Request{
		{},
		{LicenseId: 0},
	}
	for _, request := range requests {
		response, err := client.DescribeLicenseV1(ctx, request)

		assert.Nil(t, response)
		assert.NotNil(t, err)

		er, _ := status.FromError(err)

		assert.Equal(t, codes.InvalidArgument, er.Code())
		assert.Equal(t, "invalid DescribeLicenseV1Request.LicenseId: value must be greater than 0", er.Message())
	}
}

func TestLicenseAPI_ListLicenseV1(t *testing.T) {
	ctx := context.Background()

	client, closeCl := prepareClient(ctx, t)
	defer closeCl()

	requests := []*pb.ListLicenseV1Request{
		{},
		{},
	}
	t.Skip()
	for _, request := range requests {
		response, err := client.ListLicenseV1(ctx, request)

		assert.Nil(t, response)
		assert.NotNil(t, err)

		er, _ := status.FromError(err)

		assert.Equal(t, codes.InvalidArgument, er.Code())
		assert.Equal(t, "invalid ListLicenseV1Request.LicenseId: value must be greater than 0", er.Message())
	}
}

func TestLicenseAPI_RemoveLicenseV1(t *testing.T) {
	ctx := context.Background()

	client, closeCl := prepareClient(ctx, t)
	defer closeCl()

	requests := []*pb.RemoveLicenseV1Request{
		{},
		{LicenseId: 0},
	}
	for _, request := range requests {
		response, err := client.RemoveLicenseV1(ctx, request)

		assert.Nil(t, response)
		assert.NotNil(t, err)

		er, _ := status.FromError(err)

		assert.Equal(t, codes.InvalidArgument, er.Code())
		assert.Equal(t, "invalid RemoveLicenseV1Request.LicenseId: value must be greater than 0", er.Message())
	}
}
