// internal/app/controller/controller.go
package controller

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"trainingFinder/internal/config"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Controller struct {
	grpcServer *grpc.Server
	httpMux    *runtime.ServeMux
	cfg        *config.Config
	wg         sync.WaitGroup
}

func New(cfg *config.Config, grpcServer *grpc.Server, httpMux *runtime.ServeMux) *Controller {
	return &Controller{cfg: cfg, grpcServer: grpcServer, httpMux: httpMux}
}

func (c *Controller) RegisterServices(registerFuncs ...func(*grpc.Server, *runtime.ServeMux) error) error {
	for _, f := range registerFuncs {
		if err := f(c.grpcServer, c.httpMux); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) RunGRPC(ctx context.Context) error {
	grpcPort := c.cfg.Server.GRPCPort
	go func() {
		lis, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatal("failure with listen:", err)
		}
		log.Printf("grpc server listening on %s", grpcPort)
		if err := c.grpcServer.Serve(lis); err != nil {
			log.Fatal("grpc server error:", err)
		}
	}()
	return nil
}

func (c *Controller) RunHTTP(ctx context.Context) error {
	srv := &http.Server{Addr: c.cfg.Server.HTTPPort, Handler: c.httpMux}
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		<-ctx.Done()
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctxShutdown)
	}()
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("http server error:", err)
		}
	}()
	return nil
}

func (c *Controller) Wait() {
	c.wg.Wait()
}

