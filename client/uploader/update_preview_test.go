package main

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/mathantunes/arex_project/services"
	"google.golang.org/grpc"
)

func TestUpdateInvoicePreview(t *testing.T) {
	/* INPUTS */
	fileBytes := readFile("./testdata/current_invoice_preview.pdf")
	invoiceNumber := int64(10000)

	addr := ":5000"
	chunckSize := 1000
	/* INPUTS END */
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Error(err)
		return
	}

	client := services.NewInvoiceUploaderClient(conn)

	stream, err := client.UpdateInvoicePreview(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	inputBuffer := bytes.NewReader(fileBytes)
	outBufferForStream := make([]byte, chunckSize)
	writing := true

	for writing {
		n, err := inputBuffer.Read(outBufferForStream)
		if err != nil {
			if err == io.EOF {
				response, err := stream.CloseAndRecv()
				if err != nil {
					t.Error(err)
				}
				t.Log(response)
				writing = false
			}
			t.Log(err)
		}
		stream.Send(&services.InvoicePreview{
			Preview:       outBufferForStream[:n],
			InvoiceNumber: invoiceNumber,
		})
	}
}
