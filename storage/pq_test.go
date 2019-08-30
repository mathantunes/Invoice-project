package main

import (
	"database/sql"
	"fmt"
	"log"
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
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			QueryInvoice(tt.args.inv, tt.args.fun)
		})
	}
}

func Test_insertAR(t *testing.T) {

	type args struct {
		tx  *sql.Tx
		inv *services.InternalInvoice
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success Insert AR",
			args: args{initDB(), &services.InternalInvoice{
				Type:                services.InvoiceType_AR,
				CustomerID:          "B0EEBC99-9C0B-4EF8-BB6D-6BB9BD370A10",
				InvoiceNumber:       1,
				Currency:            "EUR",
				FaceValue:           200,
				CounterPartyVAT:     "123123",
				CounterPartyCountry: "FI",
				IssueDate:           "20190101",
				DueDate:             "20190102",
				ValidVAT:            true,
				CompanyName:         "TestCompany",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := insertAR(tt.args.tx, tt.args.inv); (err != nil) != tt.wantErr {
				t.Errorf("insertAR() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_insertAP(t *testing.T) {
	type args struct {
		tx  *sql.Tx
		inv *services.InternalInvoice
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success Insert AP",
			args: args{initDB(), &services.InternalInvoice{
				Type:                services.InvoiceType_AP,
				CustomerID:          "B0EEBC99-9C0B-4EF8-BB6D-6BB9BD370A10",
				InvoiceNumber:       1,
				Currency:            "EUR",
				FaceValue:           200,
				CounterPartyVAT:     "123123",
				CounterPartyCountry: "FI",
				IssueDate:           "20190101",
				DueDate:             "20190102",
				ValidVAT:            true,
				CompanyName:         "TestCompany",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := insertAP(tt.args.tx, tt.args.inv); (err != nil) != tt.wantErr {
				t.Errorf("insertAP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_updateAP(t *testing.T) {
	type args struct {
		tx  *sql.Tx
		inv *services.InternalInvoice
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success Update AP",
			args: args{initDB(), &services.InternalInvoice{
				Type:            services.InvoiceType_AP,
				InvoiceNumber:   1,
				CounterPartyVAT: "123123",
				ValidVAT:        true,
				CompanyName:     "TestCompany",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := updateAP(tt.args.tx, tt.args.inv); (err != nil) != tt.wantErr {
				t.Errorf("updateAP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_updateAR(t *testing.T) {
	type args struct {
		tx  *sql.Tx
		inv *services.InternalInvoice
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success Update AR",
			args: args{initDB(), &services.InternalInvoice{
				Type:            services.InvoiceType_AR,
				CustomerID:      "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				InvoiceNumber:   10000,
				CounterPartyVAT: "13078237",
				ValidVAT:        true,
				CompanyName:     "TestCompany",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := updateAR(tt.args.tx, tt.args.inv); (err != nil) != tt.wantErr {
				t.Errorf("updateAR() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func initDB() *sql.Tx {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer db.Close()

	tx, _ := db.Begin()
	return tx
}
