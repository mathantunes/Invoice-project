package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/mathantunes/arex_project/invoice"
	"github.com/mathantunes/arex_project/services"
	"github.com/mathantunes/arex_project/validator"
	"google.golang.org/grpc"
)

/*

 */

func main() {
	port := flag.String("port", "5000", "grpc port for listener")
	flag.Parse()
	runServer(fmt.Sprintf(":%v", *port))
}

func runServer(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	services.RegisterInvoiceUploaderServer(grpcServer, &invoice.UploaderServer{
		&validator.VIESValidator{},
	})
	// pb.RegisterRouteGuideServer(grpcServer, &routeGuideServer{})
	grpcServer.Serve(lis)
}
