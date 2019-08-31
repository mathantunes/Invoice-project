package main

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/mathantunes/arex_project/services"
	"google.golang.org/grpc"
)

func TestGetAttachment(t *testing.T) {
	/* INPUTS */
	filename := "1000"
	addr := ":5050"

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Error(err)
	}

	client := services.NewInvoiceGetterClient(conn)
	stream, err := client.GetAttachment(context.Background(), &services.QueryAttachment{
		Filename: filename,
	})
	if err != nil {
		t.Error(err)
		return
	}
	outputBytes := make([]byte, 100000)
	for {
		attachment, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				goto END
			}
			t.Error(err)
		}

		outputBytes = append(outputBytes, attachment.GetData()...)
	}
END:
		f, err := os.Create("./testdata/get_attachment.pdf", )
		if err != nil {
			t.Error(err)
		}
		f.Write(outputBytes)
}
