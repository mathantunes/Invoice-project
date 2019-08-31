package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/mathantunes/arex_project/services"
	"google.golang.org/grpc"
)

var readFile = func(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return b
}

func TestCreateXMLInvoice(t *testing.T) {
	/* INPUTS */
	fileBytes := readFile("./testdata/current_invoice.xml")
	addr := ":5000"

	//INVOICE INPUTS
	issuerID := "AAEEBC99-9C0B-4EF8-BB6D-6BB9BD380A10"
	invoiceType := services.InvoiceType_AP //services.InvoiceType_AR
	/* INPUTS END */

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Error(err)
	}

	client := services.NewInvoiceUploaderClient(conn)

	response, err := client.CreateXMLInvoice(context.Background(), &services.Invoice{
		Data:     fileBytes,
		IssuerId: issuerID,
		Type:     invoiceType,
	},
	)
	if err != nil {
		t.Error(err)
	}
	t.Log(response)
}
