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
		want    bool
		wantErr bool
	}{
		{
			name:    "Valid Request from FINLAND",
			args:    args{"FI", "25160553"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Invalid VAT Request from FINLAND",
			args:    args{"FI", "000000"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "Invalid Parameters on Request",
			args:    args{"", "0"},
			want:    false,
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt)
			}
		})
	}
}
