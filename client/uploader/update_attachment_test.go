package main

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/mathantunes/arex_project/services"
	"google.golang.org/grpc"
)

// message Attachment {
//     int64 InvoiceNumber = 1;
//     bytes Content = 2;
//     string Issuer_id = 3;
//     string CounterpartyVat = 4;
// }
func TestUpdateAttachment(t *testing.T) {
	/* INPUTS */
	fileBytes := readFile("./testdata/current_invoice_attachment.pdf")
	invoiceNumber := int64(10000)

	addr := ":6020"
	chunckSize := 1000
	/* INPUTS END */
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Error(err)
		return
	}

	client := services.NewInvoiceUploaderClient(conn)

	stream, err := client.UpdateAttachment(context.Background())
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
				response, _ := stream.CloseAndRecv()
				t.Log(response)
				writing = false
			}
			t.Log(err)
		}
		stream.Send(&services.Attachment{
			Content:       outBufferForStream[:n],
			InvoiceNumber: invoiceNumber,
		})
	}
}
