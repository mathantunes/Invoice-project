package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/mathantunes/arex_project/filestore"
	"github.com/mathantunes/arex_project/invoice"
	"github.com/mathantunes/arex_project/queuer"
	"github.com/mathantunes/arex_project/services"
	"github.com/mathantunes/arex_project/validator"
	"google.golang.org/grpc"
)

const (
	InvoicePreviewBucket = "preview_bucket"
	AttachmentsBucket    = "attachment_bucket"
)

func main() {

	fileManager := filestore.New()
	//Create buckets
	err := fileManager.CreateBucket(InvoicePreviewBucket)
	if err != nil {
		log.Println(err)
	}
	err = fileManager.CreateBucket(AttachmentsBucket)
	if err != nil {
		log.Println(err)
	}

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
		queuer.New(),
	})

	services.RegisterInvoiceGetterServer(grpcServer, &invoice.GetterServer{})
	// pb.RegisterRouteGuideServer(grpcServer, &routeGuideServer{})
	grpcServer.Serve(lis)
}
