package invoice

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	services "github.com/mathantunes/arex_project/services"
	"golang.org/x/net/html/charset"
)

func parseInvoiceInfo(inv *services.Invoice) (*services.InternalInvoice, error) {
	//Byte payload for Invoice xml standards
	invoiceData := inv.GetData()

	if invoiceData == nil {
		return nil, fmt.Errorf("Invoice received for customer %v is invalid", inv.GetIssuerId())
	}

	//Translate byte payload to XML standards
	invoiceXML, err := parseInvoiceXML(invoiceData)
	if err != nil {
		return nil, err
	}

	//Translate to a simple Invoice Structure
	invoice := &services.InternalInvoice{
		Type:          inv.GetType(),
		CustomerID:    inv.GetIssuerId(),
		InvoiceNumber: int64(invoiceXML.Details.InvoiceNumber),
		Currency:      invoiceXML.Details.TotalVAT.Currency,
		IssueDate:     invoiceXML.Details.IssueDate,
		DueDate:       invoiceXML.Details.PaymentDetails.DueDate,
		InvoiceFile:   invoiceData,
	}
	//Refactor Invoice value: Remove comma from number and multiply by 100 to have integer number
	intFaceValue, err := strconv.Atoi(strings.Replace(invoiceXML.Details.FaceValue, ",", "", -1))
	if err != nil {
		return nil, err
	}

	invoice.FaceValue = int64(intFaceValue)
	switch inv.GetType() {
	case services.InvoiceType_AP:
		invoice.CounterPartyCountry, invoice.CounterPartyVAT = separateCountryAndVAT(invoiceXML.SellerPartyDetails.SellerOrganisationTaxCode)
	case services.InvoiceType_AR:
		invoice.CounterPartyCountry, invoice.CounterPartyVAT = separateCountryAndVAT(invoiceXML.BuyerPartyDetails.BuyerOrganisationTaxCode)
	}

	return invoice, nil
}

func parseInvoiceXML(data []byte) (invoiceXML, error) {
	invoice := new(invoiceXML)

	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = charset.NewReaderLabel
	err := decoder.Decode(invoice)

	return *invoice, err
}

func separateCountryAndVAT(full string) (country, vat string) {
	var countryRune, vatRune []rune
	for _, c := range full {
		switch {
		case c >= 'A' && c <= 'Z':
			countryRune = append(countryRune, c)
		case c >= '0' && c <= '9':
			vatRune = append(vatRune, c)
		}
	}
	country = string(countryRune)
	vat = string(vatRune)
	return country, vat
}
