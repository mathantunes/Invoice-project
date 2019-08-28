package validator

import "encoding/xml"

// Validator Interface implemented by validators
type Validator interface {
	Validate(countryCode, vatNumber string) (ValidationResponse, error)
}

// ValidationVAT Response message for valid responses from VIES
type ValidationVAT struct {
	CountryCode string `xml:"countryCode"`
	VatNumber   string `xml:"vatNumber"`
	RequestDate string `xml:"requestDate"`
	Valid       bool   `xml:"valid"`
	Name        string `xml:"name"`
	Address     string `xml:"address"`
}

// ValidationFault Response message for invalid responses from VIES
type ValidationFault struct {
	FaultCode   string `xml:"faultcode"`
	FaultString string `xml:"faultstring"`
}

// ValidationBody XML Interpretation Body for SOAP. Should not be used outside of this package
type ValidationBody struct {
	XMLName  xml.Name
	CheckVat ValidationVAT   `xml:"checkVatResponse"`
	Fault    ValidationFault `xml:"Fault"`
}

// ValidationResponse XML Interpretation of SOAP Response. Should not be used outside of this package
type ValidationResponse struct {
	XMLName xml.Name
	Body    ValidationBody
}
