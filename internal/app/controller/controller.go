package controller

import (
	"context"
	"log"
	"net"
	"net/http"

	authMiddlewere "github.com/Antesser/trainingFinder/internal/app/controller/middleware/grpc"

	loggerMiddlewere "github.com/Antesser/trainingFinder/internal/app/controller/middleware/logger"
	"github.com/Antesser/trainingFinder/internal/config"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

type ImplementationAdapter interface {
	RegisterHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
	RegisterServer(grpcServer *grpc.Server)
}

type Controller interface {
	ServeHTTP(ctx context.Context)
	ServeGRPC()
	Run(ctx context.Context)
}

type controller struct {
	cfg             config.ServerConfig
	authCfg         config.AuthConfig
	implementations []ImplementationAdapter // можно запихнуть все серверы в этот интерфейс и сделать по красоте
}

func New(cfg config.ServerConfig, authCfg config.AuthConfig, implementations ...ImplementationAdapter) Controller {
	return &controller{
		cfg:             cfg,
		authCfg:         authCfg,
		implementations: implementations,
	}
}

func (c *controller) Run(ctx context.Context) {
	c.ServeGRPC()
	c.ServeHTTP(ctx)
}

func (c *controller) ServeGRPC() {
	lis, err := net.Listen("tcp", c.cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed with error %v to listen grpc port: %s", err, c.cfg.GRPCPort)
	}
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			authMiddlewere.WithAuth(c.cfg.Secret, c.authCfg),
			loggerMiddlewere.WithLogging(),
		),
	)

	for _, imp := range c.implementations {
		imp.RegisterServer(s)
	}

	go func() {
		reflection.Register(s)
		log.Printf("grpc addr: %s", lis.Addr())
		if err = s.Serve(lis); err != nil {
			log.Fatal(err, "failed to serve grpc")
		}
	}()
}

func (c *controller) ServeHTTP(ctx context.Context) {
	runtimeMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{UseProtoNames: true, EmitUnpopulated: true},
		}),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	var err error
	for _, imp := range c.implementations {
		err = multierr.Append(err, imp.RegisterHandlerFromEndpoint(ctx, runtimeMux, c.cfg.GRPCPort, opts))
	}

	if err != nil {
		log.Fatal(err, "failed to register gateway")
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/", runtimeMux)
	prefix := "/docs/"

	fs := http.FileServer(http.Dir("./swagger/"))
	httpMux.Handle(prefix, http.StripPrefix(prefix, fs))

	// main http
	go func() {
		log.Printf("http addr: %s", c.cfg.HTTPPort)
		log.Printf("swagger addr: %s/docs", c.cfg.HTTPPort)
		if err = http.ListenAndServe(c.cfg.HTTPPort,
			cors.AllowAll().Handler(httpMux)); err != nil {
			log.Fatal(err, "failed to serve http")
		}
	}()
}
