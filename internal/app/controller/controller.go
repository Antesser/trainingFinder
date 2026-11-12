package controller

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"time"

	authMiddlewere "github.com/Antesser/trainingFinder/internal/app/controller/middleware/grpc"
	loggerMiddlewere "github.com/Antesser/trainingFinder/internal/app/controller/middleware/logger"
	"github.com/Antesser/trainingFinder/internal/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type ImplementationAdapter interface {
	RegisterHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
	RegisterServer(grpcServer *grpc.Server)
}

type Controller interface {
	ServeHTTP(ctx context.Context)
	ServeGRPC(ctx context.Context)
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
	c.ServeGRPC(ctx)
	c.ServeHTTP(ctx)
}

func (c *controller) ServeGRPC(ctx context.Context) {
	lc := net.ListenConfig{}
	lis, err := lc.Listen(ctx, "tcp", c.cfg.GRPCPort)
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
			UseProtoNames: true, EmitUnpopulated: true,
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

	srv := &http.Server{
		Addr:              c.cfg.HTTPPort,
		Handler:           cors.AllowAll().Handler(httpMux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("http addr: %s", c.cfg.HTTPPort)
		log.Printf("swagger addr: %s/docs", c.cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to serve http: %v", err)
		}
	}()
}
