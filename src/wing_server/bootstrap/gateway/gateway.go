package gateway

import (
	"context"
	"crypto/tls"
	"flag"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"log"
	"net/http"
	bGRpc "wing_server/bootstrap/grpc"
	"github.com/golang/glog"
	"fmt"
	"google.golang.org/grpc/connectivity"
	pb "wing_server/modules/sandbox/proto"
)

var (
	Address = flag.String("endpoint", ":15747", "endpoint of the gRPC service")
)

// fetchToken simulates a token lookup and omits the details of proper token
// acquisition. For examples of how to acquire an OAuth2 token, see:
// https://godoc.org/golang.org/x/oauth2
func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: bGRpc.Token,
	}
}

// healthz returns a simple health handler which returns ok.
func healthz(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if s := conn.GetState(); s != connectivity.Ready {
			http.Error(w, fmt.Sprintf("grpc server is %s", s), http.StatusBadGateway)
			return
		}
		fmt.Fprintln(w, "ok")
	}
}

func authorized(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Authorization"); origin != bGRpc.Token {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func Run(ctx context.Context, cert tls.Certificate) {
	perRPC := oauth.NewOauthAccess(fetchToken())
	conn, err := grpc.DialContext(
		ctx,
		*bGRpc.Address,
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
	mux.HandleFunc("/healthz", healthz(conn))

	gwMux := gwruntime.NewServeMux()
	if err := pb.RegisterSandboxHandler(ctx, gwMux, conn); err != nil {
		log.Fatalf("Register sandbox handler fail, %v", err)
	}

	mux.Handle("/", gwMux)

	s1 := &http.Server{
		Addr:    *Address,
		Handler: authorized(mux),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}
	go func() {
		<-ctx.Done()
		log.Println("Shutting down the http server")
		if err := s1.Shutdown(ctx); err != nil {
			log.Fatalf("Failed to shutdown http server: %v", err)
		}
	}()

	log.Printf("Starting gateway listening at %v", Address)
	if err := s1.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
		log.Fatalf("Failed to listen and serve: %v", err)
	}
}
