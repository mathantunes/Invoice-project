package invoice

import (
	"context"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	services "github.com/mathantunes/arex_project/services"
)

func TestUploaderServer_CreateXMLInvoice(t *testing.T) {
	readFile := func(filename string) []byte {
		file, err := os.Open(filename)
		if err != nil {
			t.Error(err)
		}

		b, err := ioutil.ReadAll(file)
		if err != nil {
			t.Error(err)
		}
		return b
	}

	type args struct {
		ctx      context.Context
		req      *services.Invoice
		filePath string
	}
	tests := []struct {
		name    string
		sv      *UploaderServer
		args    args
		want    *services.Response
		wantErr bool
	}{
		{
			name: "Successful Call",
			sv:   &UploaderServer{&ValidationMock{}, &QueueMock{}},
			args: args{ctx: context.Background(),
				req: &services.Invoice{
					IssuerId: 123,
					Type:     services.InvoiceType_AR,
				},
				filePath: "./testdata/invoice.xml"},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileBytes := readFile(tt.args.filePath)
			tt.args.req.Data = fileBytes
			got, err := tt.sv.CreateXMLInvoice(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploaderServer.CreateXMLInvoice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UploaderServer.CreateXMLInvoice() = %v, want %v", got, tt.want)
			}
		})
	}
}

type ValidationMock struct{}

func (v *ValidationMock) Validate(countryCode, vatNumber string) (bool, error) {
	return true, nil
}

type QueueMock struct{}

func (q *QueueMock) GetQueueURL(queueName string) (string, error) {
	return "YES", nil
}
func (q *QueueMock) WriteToQueue(queueURL string, body []byte) error {
	return nil
}

func (q *QueueMock) ReadFromQueue(queueURL string) (string, error) {
	return "", nil
}

func (q *QueueMock) CreateQueue(queueURL string) error {
	return nil
}

func (q *QueueMock) Init() error { return nil }
