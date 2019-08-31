package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/mathantunes/arex_project/services"
	"google.golang.org/grpc"
)

func TestGetAttachments(t *testing.T) {
	/* INPUTS */
	invoiceNumber := int64(90000)
	addr := ":6020"
	/* INPUTS */

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Error(err)
	}

	client := services.NewInvoiceGetterClient(conn)
	resp, err := client.GetAttachments(context.Background(), &services.QueryInvoice{
		InvoiceNumber: invoiceNumber,
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp)
}
