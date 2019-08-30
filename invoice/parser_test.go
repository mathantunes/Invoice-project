package invoice

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	services "github.com/mathantunes/arex_project/services"
)

func Test_parseInvoiceInfo(t *testing.T) {
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
		inv      *services.Invoice
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    *services.InternalInvoice
		wantErr bool
		ignoreFile bool
	}{
		{
			name: "Successful Parse",
			args: args{&services.Invoice{
				IssuerId: "123",
				Type:     services.InvoiceType_AR,
			}, "./testdata/invoice.xml"},
			want: &services.InternalInvoice{
				Type:                services.InvoiceType_AR,
				CustomerID:          "123",
				InvoiceNumber:       10000,
				Currency:            "EUR",
				FaceValue:           110261,
				CounterPartyVAT:     "13078237",
				CounterPartyCountry: "FI",
				IssueDate:           "20190813",
				DueDate:             "20190823",
			},
			wantErr: false,
			ignoreFile: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileBytes := readFile(tt.args.filePath)
			tt.args.inv.Data = fileBytes
			got, err := parseInvoiceInfo(tt.args.inv)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseInvoiceInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.ignoreFile {
				got.InvoiceFile = nil
				tt.want.InvoiceFile = nil
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseInvoiceInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
