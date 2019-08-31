package invoice

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/mathantunes/arex_project/filestore"
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

	// InvoicePreviewBucket Bucket Name
	InvoicePreviewBucket = "preview_bucket"
	// AttachmentsBucket Bucket Name
	AttachmentsBucket = "attachment_bucket"
)

// CreateXMLInvoice Accepts an Invoice Structure and parses all needed fields
// Calls VAT validator and sends Protobuf payload to SQS
func (sv *UploaderServer) CreateXMLInvoice(ctx context.Context, req *services.Invoice) (*services.Response, error) {
	log.Println("Received New XML Invoice for Customer: %v", req.GetIssuerId())

	//Get invoice data from XML
	invoice, err := parseInvoiceInfo(req)
	if err != nil {
		return &services.Response{}, err
	}

	//WaitGroup allows Validation and Queuing to happen
	//concurrently whilst waiting for both to finish before ending the service
	var wg sync.WaitGroup
	wg.Add(2)

	// Helper function to queue invoice concurrently
	queueInvoice := func(wg *sync.WaitGroup, queueName string, inv services.InternalInvoice) error {
		//Get Queue URL for Creating Invoice
		url, err := sv.GetQueueURL(queueName)
		if err != nil {
			return err
		}

		//Marshal Payload as Protobuf serialization
		payload, err := proto.Marshal(&inv)
		encoded := base64.StdEncoding.EncodeToString(payload)
		err = sv.WriteToQueue(url, []byte(encoded))
		if err != nil {
			return err
		}
		wg.Done()
		return nil
	}

	//Start Creating Invoice on Queue
	go func() {
		err = queueInvoice(&wg, CreateInvoiceQueue, *invoice)
		if err != nil {
			//Allow for retry
			log.Println(err)
		}
	}()

	//Start Validation
	validationResponse, err := sv.Validate(invoice.CounterPartyCountry, invoice.CounterPartyVAT)
	if err != nil {
		return &services.Response{}, err
	}
	invoice.ValidVAT = validationResponse.Valid
	invoice.CompanyName = validationResponse.CompanyName
	//Start Creating Invoice Update on Queue
	go func() {
		err = queueInvoice(&wg, UpdateInvoiceQueue, *invoice)
		if err != nil {
			//Allow for retry
			fmt.Println(err)
		}
	}()
	wg.Wait()

	return &services.Response{
		Status: services.EStatus_Ok,
	}, nil
}

// UpdateInvoicePreview receives stream of data for file and store file
func (sv *UploaderServer) UpdateInvoicePreview(stream services.InvoiceUploader_UpdateInvoicePreviewServer) error {
	log.Println("Received Update Invoice Preview")
	var fileBytes []byte
	var invoiceNumber int64
	for {
		invoice, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				//Finished TRANSMISSION, use GOTO key word to go to END
				goto END
			}
			return err
		}
		//Received nil byte chunk
		if invoice.GetPreview() == nil {
			return errors.New("Got Invoice Preview nil")
		}
		//Save InvoiceNumber
		invoiceNumber = invoice.GetInvoiceNumber()

		//Append Byte Chunk received
		fileBytes = append(fileBytes, invoice.GetPreview()...)
	}
	//END OF TRANSMISSION, upload file to s3 now
END:
	//Initialize File Storage
	fileManager := filestore.New()
	//Upload file
	err := fileManager.Upload(InvoicePreviewBucket, fmt.Sprintf("%v.pdf", invoiceNumber), bytes.NewReader(fileBytes))
	return err
}

// UpdateAttachment receives stream of data for file and store file
func (sv *UploaderServer) UpdateAttachment(stream services.InvoiceUploader_UpdateAttachmentServer) error {
	log.Println("Received Update Invoice Attachment")
	var fileBytes []byte
	var invoiceNumber int64
	//Loop to receive incoming stream data
	for {
		invoice, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				//Finished TRANSMISSION, use GOTO key word to go to END
				goto END
			}
			return err
		}
		//Received nil byte chunk
		if invoice.GetContent() == nil {
			return errors.New("Got Invoice Preview nil")
		}
		//Save InvoiceNumber
		invoiceNumber = invoice.GetInvoiceNumber()

		//Append Byte Chunk received
		fileBytes = append(fileBytes, invoice.GetContent()...)
	}
	//END OF TRANSMISSION, upload file to s3 now
END:
	//Initialize File Storage
	fileManager := filestore.New()
	//Upload file
	err := fileManager.Upload(AttachmentsBucket, fmt.Sprintf("%v/%v.pdf", invoiceNumber, time.Now().Unix()), bytes.NewReader(fileBytes))
	return err
}

// UpdateCounterPartyVAT Validates the Incoming VAT and Sends to the Queue for storing
func (sv *UploaderServer) UpdateCounterPartyVAT(ctx context.Context, req *services.CounterPartyVAT) (*services.Response, error) {
	//Verify Inputs
	if req.GetVAT() == "" || req.GetCountry() == "" {
		return nil, errors.New("UpdateCounterPartyVAT called with empty parameters")
	}

	log.Println("Received Update CounterParty for InvoiceNumber: %v", req.GetInvoiceNumber())
	//Call validator
	validationResponse, err := sv.Validate(req.GetCountry(), req.GetVAT())
	if err != nil {
		return nil, err
	}

	//Create Invoice to send to QUEUE
	invoice := services.InternalInvoice{
		ValidVAT:            validationResponse.Valid,
		CompanyName:         validationResponse.CompanyName,
		CounterPartyCountry: req.GetCountry(),
		CounterPartyVAT:     req.GetVAT(),
		InvoiceNumber:       req.GetInvoiceNumber(),
		Type:                req.GetType(),
	}

	url, err := sv.GetQueueURL(UpdateInvoiceQueue)
	if err != nil {
		return nil, err
	}

	//Marshal Payload as Protobuf serialization
	payload, err := proto.Marshal(&invoice)
	//Encode to base64 because SQS only allows ASCII Encoding
	encoded := base64.StdEncoding.EncodeToString(payload)
	//Send to QUEUE
	err = sv.WriteToQueue(url, []byte(encoded))

	return &services.Response{
		Status: services.EStatus_Ok,
	}, err
}
