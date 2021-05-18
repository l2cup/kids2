package network

import (
	"net"

	"github.com/l2cup/kids2/internal/errors"
	"google.golang.org/grpc"
)

func DialGRPC(ni *Info, opts ...grpc.DialOption) (*grpc.ClientConn, errors.Error) {
	opts = append(opts, grpc.WithInsecure(), grpc.WithBlock())
	conn, err := grpc.Dial(ni.Address(), opts...)
	if err != nil {
		return nil, errors.NewInternalError(
			"couldn't dial gRPC connection",
			errors.InternalServerError,
			err,
		)
	}
	return conn, errors.Nil()
}

func ListenAndServeGRPC(ni *Info, srv *grpc.Server) errors.Error {
	listen, err := net.Listen("tcp", ni.Address())
	if err != nil {
		return errors.NewInternalError(
			"couldn't listen grpc",
			errors.InternalServerError,
			err,
			"addr", ni.Address(),
			"node_id", ni.ID,
		)
	}

	if err := srv.Serve(listen); err != nil {
		return errors.NewInternalError(
			"failed to serve grpc",
			errors.InternalServerError,
			err,
		)
	}

	return errors.Nil()
}
