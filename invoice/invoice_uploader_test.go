package invoice

import (
	"context"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	services "github.com/mathantunes/arex_project/services"
	"github.com/mathantunes/arex_project/validator"
)

var initEnvResult = os.Setenv("S3_ENDPOINT", "localhost:4572")

var readFile = func(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil
	}
	return b
}

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
					IssuerId: "A0EEBC99-9C0B-4EF8-BB6D-6BB9BD380A10",
					Type:     services.InvoiceType_AP,
				},
				filePath: "./testdata/invoice.xml"},
			want: &services.Response{
				Status: services.EStatus_Ok,
			},
			wantErr: false,
		},
		// {
		// 	name: "Successful Call to real",
		// 	sv:   &UploaderServer{&ValidationMock{}, queuer.New()},
		// 	args: args{ctx: context.Background(),
		// 		req: &services.Invoice{
		// 			IssuerId: "A0EEBC99-9C0B-4EF8-BB6D-6BB9BD380A21",
		// 			Type:     services.InvoiceType_AP,
		// 		},
		// 		filePath: "./testdata/invoice.xml"},
		// 	want:    nil,
		// 	wantErr: false,
		// },
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

func (v *ValidationMock) Validate(countryCode, vatNumber string) (validator.InternalResponse, error) {
	return validator.InternalResponse{}, nil
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

func TestUploaderServer_UpdateCounterPartyVAT(t *testing.T) {
	type args struct {
		ctx context.Context
		req *services.CounterPartyVAT
	}
	tests := []struct {
		name    string
		sv      *UploaderServer
		args    args
		want    *services.Response
		wantErr bool
	}{
		{
			name: "UpdateCounterPartyVAT call on MOCK",
			sv:   &UploaderServer{&ValidationMock{}, &QueueMock{}},
			args: args{ctx: context.Background(),
				req: &services.CounterPartyVAT{
					Country:       "FI",
					VAT:           "25160554",
					InvoiceNumber: 1000,
					Type:          services.InvoiceType_AR,
				}},
			want: &services.Response{
				Status: services.EStatus_Ok,
			},
		},
		// {
		// 	name: "UpdateCounterPartyVAT call on Real QUEUE",
		// 	sv:   &UploaderServer{&ValidationMock{}, queuer.New()},
		// 	args: args{ctx: context.Background(),
		// 		req: &services.CounterPartyVAT{
		// 			Country:       "FI",
		// 			VAT:           "25160553",
		// 			InvoiceNumber: 10000,
		// 			Type:          services.InvoiceType_AR,
		// 		}},
		// 	want: nil,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.sv.UpdateCounterPartyVAT(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploaderServer.UpdateCounterPartyVAT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UploaderServer.UpdateCounterPartyVAT() = %v, want %v", got, tt.want)
			}
		})
	}
}
