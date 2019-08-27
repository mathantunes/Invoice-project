package validator

import (
	"io/ioutil"
	"net/http"
)

// go get github.com/tiaguinho/gosoap

const (
	FinlandCode = "FIN"
)

type ValidationResponse struct {
}

func Validate(vatNumber int) (ValidationResponse, error) {
	resp, err := http.NewRequest(http.MethodPost, "http://ec.europa.eu/taxation_customs/vies/checkVatService", nil)
	if err != nil {
		return ValidationResponse{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ValidationResponse{}, err
	}

	//Parse xml received what ever format it is
	return ValidationResponse{}, nil
}
