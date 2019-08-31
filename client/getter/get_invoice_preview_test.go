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

func TestGetInvoicePreview(t *testing.T) {
	/* INPUTS */
	invoiceNumber := int64(90000)
	addr := ":6020"
	/* INPUTS */

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
	outputBytes := make([]byte, 1000000)
	for {
		preview, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Printf("EOF")
				goto END
			}
			goto END
		}
		outputBytes = append(outputBytes, preview.GetPreview()...)
	}
END:
	fmt.Printf("END")
	f, err := os.Create("./testdata/current_get_invoice_preview.pdf")
	if err != nil {
		t.Error(err)
	}
	_, err = f.Write(outputBytes)
	if err != nil {
		t.Error(err)
	}
	return
}
