package grpc

import (
	"errors"
	"runtime/debug"

	govalidator "buf.build/go/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func WithValidation() grpc.UnaryServerInterceptor {
	v, err := govalidator.New()
	if err != nil {
		log.Fatal("failed to initialize validator ", err)
	}

	return protovalidate.UnaryServerInterceptor(v)
}

func WithRecovery() grpc.UnaryServerInterceptor {
	grpcPanicRecoveryHandler := func(p any) (err error) {
		log.Error(errors.New("panic"), "msg recovered from panic: ", p, "stack: ", string(debug.Stack()))
		return status.Errorf(codes.Internal, "%s", p)
	}

	return recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler))
}
