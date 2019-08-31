package main

import (
	"context"
	"testing"

	"github.com/mathantunes/arex_project/services"
	"google.golang.org/grpc"
)

func TestUpdateVAT(t *testing.T) {
	/* INPUTS */
	invoiceType := services.InvoiceType_AP //services.InvoiceType_AR
	invoiceNumber := int64(10000)
	newVAT := "A123123"
	newCountry := "FI"
	addr := ":5000"
	/* INPUTS END */
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Error(err)
		return
	}

	client := services.NewInvoiceUploaderClient(conn)

	response, err := client.UpdateCounterPartyVAT(context.Background(), &services.CounterPartyVAT{
		VAT: newVAT,
		Country: newCountry,
		InvoiceNumber: invoiceNumber,
		Type: invoiceType,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(response)	
}
