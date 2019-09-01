package main

/*
	Initializes all necessary buckets to store files.

	Starts the GRPC Server with all services attached.
*/

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

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

	//Remove this in production with AWS.
	//This is done so the container initialization waits for elasticmq and localstack to initialize properly
	time.Sleep(20 * time.Second)

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

	var port = new(string)
	*port = "5000"
	*port = os.Getenv("GRPC_PORT")
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

	log.Println("Running GRPC Server on", addr)

	grpcServer.Serve(lis)
}
