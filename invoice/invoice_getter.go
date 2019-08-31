package invoice

import (
	"errors"
	"fmt"
	"io"

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
	fmt.Println("Received Invoice Preview Request for Invoice: %v", req.GetInvoiceNumber())
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
func (sv *GetterServer) GetAttachments(req *services.QueryInvoice, stream services.InvoiceGetter_GetAttachmentsServer) error {
	if req.GetInvoiceNumber() == 0 {
		return errors.New("Got Invoice Number equal to zero")
	}
	fmt.Println("Received Attachments Request for Invoice: %v", req.GetInvoiceNumber())
	fileManager := filestore.New()
	filenames, err := fileManager.ListItems(AttachmentsBucket, fmt.Sprintf("%v", req.GetInvoiceNumber()))
	if err != nil {
		return err
	}
	for _, filename := range filenames {
		readerBytes, err := fileManager.Download(AttachmentsBucket, fmt.Sprintf("%v.pdf", filename))
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
				InvoiceNumber: req.GetInvoiceNumber(),
				Data:          buffer[:n],
				Filename:      filename,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GetAttachment Downloads a single attachment
func (sv *GetterServer) GetAttachment(req *services.QueryAttachment, stream services.InvoiceGetter_GetAttachmentServer) error {
	if req.GetFilename() == "" {
		return errors.New("Got Filename empty")
	}
	fmt.Println("Received Attachment Request for File: %v", req.GetFilename())
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
