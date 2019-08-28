package invoice

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	services "github.com/mathantunes/arex_project/services"
	"golang.org/x/net/html/charset"
)

// UploaderServer Holds the GRPC server implementation
type UploaderServer struct {
}

func (sv *UploaderServer) CreateXMLInvoice(ctx context.Context, req *services.Invoice) (*services.Response, error) {

	return nil, nil

}

func (sv *UploaderServer) UpdateInvoicePreview(srv services.InvoiceUploader_UpdateInvoicePreviewServer) error {
	return nil

}

func (sv *UploaderServer) UpdateAttachment(srv services.InvoiceUploader_UpdateAttachmentServer) error {
	return nil

}

func parseInvoiceInfo(inv *services.Invoice) (*Invoice, error) {
	invoiceData := inv.GetData()

	if invoiceData == nil {
		return nil, fmt.Errorf("Invoice received for customer %v is invalid", inv.GetIssuerId())
	}

	invoiceXML, err := parseInvoiceXML(invoiceData)
	if err != nil {
		return nil, err
	}

	invoice := Invoice{
		Currency:      invoiceXML.Details.TotalVAT.Currency,
		CustomerID:    uint16(inv.GetIssuerId()),
		DueDate:       invoiceXML.Details.PaymentDetails.DueDate,
		InvoiceNumber: invoiceXML.Details.InvoiceNumber,
		IssueDate:     invoiceXML.Details.IssueDate,
	}
	intFaceValue, err := strconv.Atoi(strings.Replace(invoiceXML.Details.FaceValue, ",", "", -1))
	if err != nil {
		return nil, err
	}

	invoice.FaceValue = intFaceValue * 100

	return &Invoice{}, nil
}

func parseInvoiceXML(data []byte) (invoiceXML, error) {
	invoice := new(invoiceXML)

	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = charset.NewReaderLabel
	err := decoder.Decode(invoice)

	return *invoice, err
}
