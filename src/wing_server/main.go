/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto

package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"os"
	bGateway "wing_server/bootstrap/gateway"
	bGRpc "wing_server/bootstrap/grpc"
	"google.golang.org/grpc/testdata"
	"runtime"
	"time"
)

var (
	serverPem = os.Getenv("SERVER_PEM")
	serverKey = os.Getenv("SERVER_KEY")
)

func main() {
	flag.Parse()

	pCtx := context.Background()
	ctx, cancel :=  context.WithCancel(pCtx)

	if runtime.GOOS == "windows" {
		serverPem = testdata.Path("server1.pem")
		serverKey = testdata.Path("server1.key")
	}
	cert, err := tls.LoadX509KeyPair(serverPem, serverKey)
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}

	go func() {
		defer cancel()
		bGRpc.Run(ctx, cert)
	}()
	go func() {
		defer cancel()
		bGateway.Run(ctx, cert)
	}()
	<-ctx.Done()
	time.Sleep(time.Second)
}
