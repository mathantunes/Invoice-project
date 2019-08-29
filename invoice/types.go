package invoice

import (
	"encoding/xml"
)

// invoiceXML definition for xml standards receiver
type invoiceXML struct {
	XMLName            xml.Name           `xml:"Finvoice"`
	Details            invoiceDetails     `xml:"InvoiceDetails"`
	SellerPartyDetails sellerPartyDetails `xml:"SellerPartyDetails"`
	BuyerPartyDetails  buyerPartyDetails  `xml:"BuyerPartyDetails"`
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

type buyerPartyDetails struct {
	BuyerOrganisationName    string `xml:"BuyerOrganisationName"`
	BuyerOrganisationTaxCode string `xml:"BuyerOrganisationTaxCode"`
}

type sellerPartyDetails struct {
	SellerOrganisationName    string `xml:"SellerOrganisationName"`
	SellerOrganisationTaxCode string `xml:"SellerOrganisationTaxCode"`
}
