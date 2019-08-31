package invoice

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/mathantunes/arex_project/filestore"
	services "github.com/mathantunes/arex_project/services"
)

type GetterServer struct{}

const (
	ChunckSize = 1000
)

// GetInvoicePreview Downloads file from storage and sends a stream of bytes to the client
func (sv *GetterServer) GetInvoicePreview(req *services.QueryInvoice, stream services.InvoiceGetter_GetInvoicePreviewServer) error {
	if req.GetInvoiceNumber() == 0 {
		return errors.New("Got Invoice Number equal to zero")
	}
	log.Println("Received Invoice Preview Request for Invoice: %v", req.GetInvoiceNumber())
	fileManager := filestore.New()
	readerBytes, err := fileManager.Download(InvoicePreviewBucket, fmt.Sprintf("%v.pdf", req.GetInvoiceNumber()))
	if err != nil {
		return err
	}
	buffer := make([]byte, ChunckSize)
	writing := true
	for writing {
		n, err := readerBytes.Read(buffer)
		if err != nil {
			if err == io.EOF {
				writing = false
			}
			return err
		}
		err = stream.Send(&services.InvoicePreview{
			InvoiceNumber: req.GetInvoiceNumber(),
			Preview:       buffer[:n],
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// GetAttachments downloads all attachments for an Invoice
func (sv *GetterServer) GetAttachments(ctx context.Context, req *services.QueryInvoice) (*services.AttachmentsResponse, error) {
	// GetAttachments(req *services.QueryInvoice, stream services.InvoiceGetter_GetAttachmentsServer) error {
	if req.GetInvoiceNumber() == 0 {
		return &services.AttachmentsResponse{}, errors.New("Got Invoice Number equal to zero")
	}
	log.Println("Received Attachments Request for Invoice: %v", req.GetInvoiceNumber())
	fileManager := filestore.New()
	filenames, err := fileManager.ListItems(AttachmentsBucket, fmt.Sprintf("%v", req.GetInvoiceNumber()))
	if err != nil {
		return &services.AttachmentsResponse{}, err
	}
	return &services.AttachmentsResponse{
		InvoiceNumber: req.GetInvoiceNumber(),
		Filenames:     filenames,
	}, nil
}

// GetAttachment Downloads a single attachment
func (sv *GetterServer) GetAttachment(req *services.QueryAttachment, stream services.InvoiceGetter_GetAttachmentServer) error {
	if req.GetFilename() == "" {
		return errors.New("Got Filename empty")
	}
	log.Println("Received Attachment Request for File: %v", req.GetFilename())
	fileManager := filestore.New()
	readerBytes, err := fileManager.Download(AttachmentsBucket, req.GetFilename())
	if err != nil {
		return err
	}
	buffer := make([]byte, ChunckSize)
	writing := true
	for writing {
		n, err := readerBytes.Read(buffer)
		if err != nil {
			if err == io.EOF {
				writing = false
			}
			return err
		}
		err = stream.Send(&services.AttachmentsResponse{
			Data:     buffer[:n],
			Filename: req.GetFilename(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
