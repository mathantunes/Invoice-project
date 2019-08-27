package invoice

import (
	"context"

	services "github.com/mathantunes/arex_project/services"
)

type invoiceUploaderServer struct {
}

func (sv *invoiceUploaderServer) CreateXMLInvoice(ctx context.Context, req *services.Invoice) (*services.Response, error) {
	// customer_id Retrieve ID of AREX Customer

	// invoice_number  Invoice Number from XML InvoiceNumber
	// currency from XML InvoiceTotalVatIncludedAmount AmountCurrencyIdentifier="EUR"
	// face_value from XML InvoiceTotalVatIncludedAmount * 100 so cents are part of the int
	// counterparty_vat ???
	// issue_date from XML InvoiceDate
	// due_date from XML InvoiceDueDate
	return nil, nil

}

func (sv *invoiceUploaderServer) UpdateInvoicePreview(srv services.InvoiceUploader_UpdateInvoicePreviewServer) error {
	return nil

}

func (sv *invoiceUploaderServer) UpdateAttachment(srv services.InvoiceUploader_UpdateAttachmentServer) error {
	return nil

}
