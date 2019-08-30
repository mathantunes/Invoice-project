package validator

import (
	"reflect"
	"testing"
)

func TestVIESValidator_Validate(t *testing.T) {
	type args struct {
		countryCode string
		vatNumber   string
	}
	tests := []struct {
		name    string
		v       *VIESValidator
		args    args
		want    InternalResponse
		wantErr bool
	}{
		{
			name: "Valid Request from FINLAND",
			v:    &VIESValidator{},
			args: args{"FI", "25160553"},
			want: InternalResponse{
				Valid:       true,
				CompanyName: "Comtower Finland Oy",
			},
			wantErr: false,
		},
		{
			name:    "Invalid VAT Request from FINLAND",
			args:    args{"FI", "000000"},
			want:    InternalResponse{Valid: false, CompanyName: "---"},
			wantErr: false,
		},
		{
			name:    "Invalid Parameters on Request",
			args:    args{"", "0"},
			want:    InternalResponse{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.v.Validate(tt.args.countryCode, tt.args.vatNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("VIESValidator.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VIESValidator.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
