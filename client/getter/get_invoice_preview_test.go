package main

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/mathantunes/arex_project/services"
	"google.golang.org/grpc"
)

func TestGetInvoicePreview(t *testing.T) {
	/* INPUTS */
	invoiceNumber := int64(10000)
	addr := ":5050"

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Error(err)
	}

	client := services.NewInvoiceGetterClient(conn)
	stream, err := client.GetInvoicePreview(context.Background(), &services.QueryInvoice{
		InvoiceNumber: invoiceNumber,
	})
	if err != nil {
		t.Error(err)
		return
	}
	outputBytes := make([]byte, 100000)
	for {
		preview, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				goto END
			}
			t.Error(err)
		}
		outputBytes = append(outputBytes, preview.GetPreview()...)
	}
END:
	f, err := os.Create("./testdata/current_get_invoice_preview.pdf")
	f.Write(outputBytes)
}
