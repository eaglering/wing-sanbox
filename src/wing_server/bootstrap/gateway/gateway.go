package gateway

import (
	"context"
	"crypto/tls"
	"flag"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/examples/gateway"
	"github.com/grpc-ecosystem/grpc-gateway/examples/proto/examplepb"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"log"
	"net/http"
	bGRpc "wing_server/bootstrap/grpc"
)

var (
	endpoint = flag.String("endpoint", ":15747", "endpoint of the gRPC service")
	network  = flag.String("network", "tcp", `one of "tcp" or "unix". Must be consistent to -endpoint`)
)

// fetchToken simulates a token lookup and omits the details of proper token
// acquisition. For examples of how to acquire an OAuth2 token, see:
// https://godoc.org/golang.org/x/oauth2
func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: bGRpc.Token,
	}
}

func Run(ctx context.Context, cert tls.Certificate) {
	opts := gateway.Options{
		Addr: ":8080",
		GRPCServer: gateway.Endpoint{
			Network: *network,
			Addr:    *endpoint,
		},
	}

	perRPC := oauth.NewOauthAccess(fetchToken())
	conn, err := grpc.DialContext(
		ctx,
		*endpoint,
		grpc.WithPerRPCCredentials(perRPC),
		grpc.WithTransportCredentials(
			credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}),
		),
	)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	go func() {
		<-ctx.Done()
		if err := conn.Close(); err != nil {
			glog.Errorf("Failed to close a client connection to the gRPC server: %v", err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthzServer(conn))

	gwMux := gwruntime.NewServeMux(opts.Mux...)

	for _, f := range []func(context.Context, *gwruntime.ServeMux, *grpc.ClientConn) error{
		examplepb.RegisterEchoServiceHandler,
		examplepb.RegisterStreamServiceHandler,
		examplepb.RegisterABitOfEverythingServiceHandler,
		examplepb.RegisterFlowCombinationHandler,
		examplepb.RegisterResponseBodyServiceHandler,
	} {
		if err := f(ctx, gwMux, conn); err != nil {
			return nil, err
		}
	}

	mux.Handle("/", gw)

	s1 := &http.Server{
		Addr:    *endpoint,
		Handler: allowCORS(mux),
	}
	go func() {
		<-ctx.Done()
		glog.Infof("Shutting down the http server")
		if err := s1.Shutdown(context.Background()); err != nil {
			glog.Errorf("Failed to shutdown http server: %v", err)
		}
	}()

	log.Printf("Starting listening at %v", endpoint)
	if err := s1.ListenAndServe(); err != http.ErrServerClosed {
		glog.Errorf("Failed to listen and serve: %v", err)
		return err
	}
}
