package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mathantunes/arex_project/services"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "arex"
)

func CreateInvoice(inv *services.InternalInvoice) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	hasCommited := false
	for ok := true; ok; ok = !hasCommited {
		tx, err := db.Begin()
		if err != nil {
			return
		}
		defer tx.Rollback()

		_, err = tx.Exec(`set transaction isolation level repeatable read`)
		if err != nil {
			return
		}

		if inv.GetType() == services.InvoiceType_AR {
			insertAR(tx, inv)
		} else if inv.GetType() == services.InvoiceType_AP {
			insertAP(tx, inv)
		}
	}
}

func insertAR(tx *sql.Tx, inv *services.InternalInvoice) {
	tx.Exec(`INSERT INTO "ar_invoices" 
	("customer_id", "invoice_number", "currency", "face_value", "counterparty_vat",
	"issue_date", "due_date")
	VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		inv.GetCustomerID(),
		inv.GetInvoiceNumber(),
		inv.GetCurrency(),
		inv.GetFaceValue(),
		inv.GetCounterPartyVAT(),
		inv.GetIssueDate(),
		inv.GetDueDate())
}
func insertAP(tx *sql.Tx, inv *services.InternalInvoice) {

}
