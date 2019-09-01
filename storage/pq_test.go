package main

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/mathantunes/arex_project/services"
)

func TestQueryInvoice(t *testing.T) {
	type args struct {
		inv *services.InternalInvoice
		fun func(*sql.Tx, *services.InternalInvoice) error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Insert AP Invoice 3",
			args: args{
				&services.InternalInvoice{
					Type:                services.InvoiceType_AP,
					CustomerID:          "B0EEBC99-9C0B-4EF8-BB6D-6BB9BD370A10",
					InvoiceNumber:       3,
					Currency:            "EUR",
					FaceValue:           200,
					CounterPartyVAT:     "123123",
					CounterPartyCountry: "FI",
					IssueDate:           "20190101",

					DueDate:     "20190102",
					ValidVAT:    false,
					CompanyName: "TestCompany",
				},
				insertAP,
			},
		},
		{
			name: "Update AP Invoice 3 to INVALID",
			args: args{
				&services.InternalInvoice{
					Type:            services.InvoiceType_AP,
					InvoiceNumber:   3,
					CounterPartyVAT: "123123",
					ValidVAT:        false,
					CompanyName:     "TestCompany",
				},
				updateAP,
			},
		},
		{
			name: "Update AP Invoice 3 to READY",
			args: args{
				&services.InternalInvoice{
					Type:            services.InvoiceType_AP,
					InvoiceNumber:   3,
					CounterPartyVAT: "123123",
					ValidVAT:        true,
					CompanyName:     "TestCompany",
				},
				updateAP,
			},
		},
		{
			name: "Insert AR Invoice 2 to PENDING",
			args: args{
				&services.InternalInvoice{
					Type:                services.InvoiceType_AR,
					CustomerID:          "B0EEBC99-9C0B-4EF8-BB6D-6BB9BD370A10",
					InvoiceNumber:       2,
					Currency:            "EUR",
					FaceValue:           200,
					CounterPartyVAT:     "123123",
					CounterPartyCountry: "FI",
					IssueDate:           "20190101",
					DueDate:             "20190102",
					ValidVAT:            false,
					CompanyName:         "TestCompany",
				},
				insertAR,
			},
		},
		{
			name: "Update AR Invoice 2 to INVALID",
			args: args{
				&services.InternalInvoice{
					Type:            services.InvoiceType_AR,
					CustomerID:      "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
					InvoiceNumber:   2,
					CounterPartyVAT: "13078237",
					ValidVAT:        false,
					CompanyName:     "TestCompany",
				},
				updateAR,
			},
		},
		{
			name: "Update AR Invoice 2 to READY",
			args: args{
				&services.InternalInvoice{
					Type:            services.InvoiceType_AR,
					CustomerID:      "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
					InvoiceNumber:   2,
					CounterPartyVAT: "13078237",
					ValidVAT:        true,
					CompanyName:     "TestCompany",
				},
				updateAR,
			},
		},
	}
	for _, tt := range tests {
		host = "localhost"
		port = "5432"
		user = "postgres"
		password = "alstom"
		t.Run(tt.name, func(t *testing.T) {
			QueryInvoice(tt.args.inv, tt.args.fun)
		})
	}
}
