package invoice

import (
	"encoding/xml"

	services "github.com/mathantunes/arex_project/services"
)

type Invoice struct {
	Type            services.InvoiceType
	CustomerID      uint16
	InvoiceNumber   int
	Currency        string
	FaceValue       int
	CounterPartyVAT string
	IssueDate       string
	DueDate         string
}

// invoiceXML definition for xml standards receiver
type invoiceXML struct {
	XMLName xml.Name       `xml:"Finvoice"`
	Details invoiceDetails `xml:"InvoiceDetails"`
}

// invoiceDetails represents the InvoiceDetails Node from standards
type invoiceDetails struct {
	InvoiceNumber  int                 `xml:"InvoiceNumber"`
	IssueDate      string              `xml:"InvoiceDate"`
	TotalVAT       vatIncluded         `xml:"InvoiceTotalVatAmount"`
	PaymentDetails paymentTermsDetails `xml:"PaymentTermsDetails"`
	FaceValue      string              `xml:"InvoiceTotalVatIncludedAmount"`
}

// vatIncluded represents the VATIncluded nodes, which contains the Currency Attribute
type vatIncluded struct {
	Currency string `xml:"AmountCurrencyIdentifier,attr"`
}

// paymentTermsDetails represents the PaymentTermsDetails Node from standards
type paymentTermsDetails struct {
	DueDate string `xml:"InvoiceDueDate"`
}
