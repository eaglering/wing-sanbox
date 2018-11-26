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

package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "wing_server/modules/sandbox/proto"
	"google.golang.org/grpc/credentials/oauth"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/credentials"
	"crypto/tls"
)

const (
	PORT = ":15747"
	TOKEN = "V9max5VOMkt3q="
)

// fetchToken simulates a token lookup and omits the details of proper token
// acquisition. For examples of how to acquire an OAuth2 token, see:
// https://godoc.org/golang.org/x/oauth2
func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: TOKEN,
	}
}

func main() {
	perRPC := oauth.NewOauthAccess(fetchToken())
	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(perRPC),
		grpc.WithTransportCredentials(
			credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}),
		),
	}
	// Set up a connection to the server.
	conn, err := grpc.Dial(PORT, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSandboxClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 120 * time.Second)
	defer cancel()
	r, err := c.Compile(ctx, &pb.Input{
		Language: "php",
		Data: "<?php phpinfo();",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	if r.Data == "" {
		log.Println("empty")
	}
	log.Printf("Greeting: %s", r)
}
