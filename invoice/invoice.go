package invoice

import (
	"context"
	"fmt"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/mathantunes/arex_project/queuer"
	services "github.com/mathantunes/arex_project/services"
	"github.com/mathantunes/arex_project/validator"
)

// UploaderServer Holds the GRPC server implementation
type UploaderServer struct {
	validator.Validator
	queuer.QueueManager
}

const (
	// CreateInvoiceQueue Queue Name
	CreateInvoiceQueue = "create_invoice"
	// UpdateInvoiceQueue Queue Name
	UpdateInvoiceQueue = "update_invoice"
)

// CreateXMLInvoice Accepts an Invoice Structure and parses all needed fields
// Calls VAT validator and sends Protobuf payload to SQS
func (sv *UploaderServer) CreateXMLInvoice(ctx context.Context, req *services.Invoice) (*services.Response, error) {

	//Get invoice data from XML
	invoice, err := parseInvoiceInfo(req)
	if err != nil {
		return nil, err
	}

	//WaitGroup allows Validation and Queuing to happen
	//concurrently whilst waiting for both to finish before ending the service
	var wg sync.WaitGroup
	wg.Add(2)

	// //Initialize SQS Connection
	// err = sv.Init()
	// if err != nil {
	// 	return nil, err
	// }

	// Helper function to queue invoice concurrently
	queueInvoice := func(wg *sync.WaitGroup, queueName string) error {
		//Get Queue URL for Creating Invoice
		url, err := sv.GetQueueURL(queueName)
		if err != nil {
			return err
		}

		//Marshal Payload as Protobuf serialization
		payload, err := proto.Marshal(invoice)
		err = sv.WriteToQueue(url, payload)
		if err != nil {
			return err
		}
		wg.Done()
		return nil
	}

	//Start Creating Invoice on Queue
	go func() {
		err = queueInvoice(&wg, CreateInvoiceQueue)
		if err != nil {
			//Allow for retry
			fmt.Println(err)
		}
	}()

	//Start Validation
	validation, err := sv.Validate(invoice.CounterPartyCountry, invoice.CounterPartyVAT)
	if err != nil {
		return nil, err
	}
	invoice.ValidVAT = validation
	//Start Creating Invoice Update on Queue
	go func() {
		err = queueInvoice(&wg, UpdateInvoiceQueue)
		if err != nil {
			//Allow for retry
			fmt.Println(err)
		}
	}()
	wg.Wait()

	return nil, nil
}

// UpdateInvoicePreview Queues the Update
func (sv *UploaderServer) UpdateInvoicePreview(srv services.InvoiceUploader_UpdateInvoicePreviewServer) error {
	return nil

}

// UpdateAttachment Queues the Update
func (sv *UploaderServer) UpdateAttachment(srv services.InvoiceUploader_UpdateAttachmentServer) error {
	return nil

}
