package validator

import (
	"reflect"
	"testing"
)

func TestValidate(t *testing.T) {
	type args struct {
		countryCode string
		vatNumber   string
	}
	tests := []struct {
		name    string
		args    args
		want    ValidationResponse
		wantErr bool
	}{
		{
			name: "Valid Request from FINLAND",
			args: args{"FI", "25160553"},
			want: ValidationResponse{
				Body: ValidationBody{
					CheckVat: ValidationVAT{
						CountryCode: "FI",
						VatNumber:   "25160553",
						RequestDate: "2019-08-28+02:00",
						Valid:       true,
						Name:        "Comtower Finland Oy",
						Address:     `Sibeliuksenkatu 3 08100 LOHJA`,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid VAT Request from FINLAND",
			args: args{"FI", "000000"},
			want: ValidationResponse{
				Body: ValidationBody{
					CheckVat: ValidationVAT{
						Valid: false,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid Parameters on Request",
			args: args{"", "0"},
			want: ValidationResponse{
				Body: ValidationBody{
					Fault: ValidationFault{
						FaultString: "INVALID_INPUT",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &VIESValidator{}
			got, err := v.Validate(tt.args.countryCode, tt.args.vatNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Body.CheckVat.Valid, tt.want.Body.CheckVat.Valid) {
				t.Errorf("Validate() = %v, want %v", got.Body.CheckVat.Valid, tt.want.Body.CheckVat.Valid)
			}
			if !reflect.DeepEqual(got.Body.Fault.FaultString, tt.want.Body.Fault.FaultString) {
				t.Errorf("Validate() = %v, want %v", got.Body.Fault.FaultString, tt.want.Body.Fault.FaultString)
			}
		})
	}
}
