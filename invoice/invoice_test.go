package invoice

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func Test_parseInvoiceXML(t *testing.T) {

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
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    invoiceXML
		wantErr bool
	}{
		{
			name: "Parse Valid Invoice XML",
			args: args{"./testdata/invoice.xml"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileBytes := readFile(tt.args.filePath)
			got, err := parseInvoiceXML(fileBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseInvoiceXML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseInvoiceXML() = %v, want %v", got, tt.want)
			}
		})
	}
}
