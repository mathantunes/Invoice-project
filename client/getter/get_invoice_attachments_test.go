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
	invoiceNumber := int64(10000)
	addr := ":6020"
	/* INPUTS */

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Error(err)
	}

	client := services.NewInvoiceGetterClient(conn)
	stream, err := client.GetAttachments(context.Background(), &services.QueryInvoice{
		InvoiceNumber: invoiceNumber,
	})
	if err != nil {
		t.Error(err)
		return
	}
	outputBytes := make([]byte, 1000000)
	files := make(map[string][]byte)
	for {
		attachments, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				goto END
			}
			fmt.Println(err)
			goto END
		}

		fmt.Printf("append")
		outputBytes = append(outputBytes, attachments.GetData()...)
		files[attachments.GetFilename()] = outputBytes
	}
END:
	fmt.Printf("END")
	for key, value := range files {
		f, err := os.Create(fmt.Sprintf("./testdata/%v.pdf", key))
		if err != nil {
			t.Error(err)
		}
		_, err = f.Write(value)
		if err != nil {
			t.Error(err)
		}
	}
}
