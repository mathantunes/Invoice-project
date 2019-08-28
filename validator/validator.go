package validator

/*
	Validator V0.0.1
	27.08.2019

	The package communicates with VIES API as a SOAP service

*/

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
)

const (
	// VIESEndpoint EndPoint for Vat Service Verification
	VIESEndpoint = "http://ec.europa.eu/taxation_customs/vies/services/checkVatService"
)

type VIESValidator struct{}

// Validate Communicates with VIES API
func (v *VIESValidator) Validate(countryCode, vatNumber string) (ValidationResponse, error) {

	//Generate request from input parameters
	requestPayload := createRequest(countryCode, vatNumber)

	//Add SOAP Action - Required for SOAP Requests
	soapAction := "urn:checkVatService"

	//Start HTTP Requester
	req, err := http.NewRequest(http.MethodPost, VIESEndpoint, bytes.NewReader(requestPayload))
	if err != nil {
		return ValidationResponse{}, err
	}
	req.Header.Set("Content-type", "text/xml")
	req.Header.Set("SOAPAction", soapAction)

	//Disable TLS Config for HTTP Client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	//Make the HTTP Request
	res, err := client.Do(req)
	if err != nil {
		return ValidationResponse{}, err
	}

	//Parse response XML into ValidationResponse Struct
	validationResponse := new(ValidationResponse)
	err = xml.NewDecoder(res.Body).Decode(validationResponse)
	if err != nil {
		return ValidationResponse{}, err
	}

	return *validationResponse, nil
}

// createRequest Generates the Request Message Payload
// Utilizes the format for checkVatService from VIES API
func createRequest(country, vat string) []byte {

	stringPayload := fmt.Sprintf(`<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
    <Body>
        <checkVat xmlns="urn:ec.europa.eu:taxud:vies:services:checkVat:types">
            <countryCode>%s</countryCode>
            <vatNumber>%s</vatNumber>
        </checkVat>
    </Body>
</Envelope>`, country, vat)

	return []byte(strings.TrimSpace(stringPayload))
}
