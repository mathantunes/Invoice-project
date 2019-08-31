package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/mathantunes/arex_project/services"
)

//The PSQL connection parameters where added as const values
var host = os.Getenv("PSQL_HOST")
var port = 5432
var user = os.Getenv("PSQL_USER")
var password = os.Getenv("PSQL_PASS")

const (
	dbname = "arex"

	// for simplicity, the status definitions were saved as const values
	pendingStatus = "PENDING"
	readyStatus   = "READY"
	invalidStatus = "INVALID"
)

// QueryInvoice Operates the connection and transaction management for queries.
// Takes a function for executing queries and commits all the changes on success
func QueryInvoice(inv *services.InternalInvoice, fun func(*sql.Tx, *services.InternalInvoice) error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	hasCommitted := false
	for ok := true; ok; ok = !hasCommitted {
		tx, err := db.Begin()
		if err != nil {
			log.Println(err)
			return
		}
		defer tx.Rollback()

		_, err = tx.Exec(`set transaction isolation level repeatable read`)
		if err != nil {
			log.Println(err)
			return
		}

		err = fun(tx, inv)

		if err != nil {
			if strings.Contains(err.Error(), "could not serialize access due to concurrent update") {
				log.Println(err)
				continue
			} else {
				log.Println(err)
				return
			}
		}
		tx.Commit()
		hasCommitted = true
	}
	return
}

// insertAR Inserts a new Invoice into AR
func insertAR(tx *sql.Tx, inv *services.InternalInvoice) error {
	_, err := tx.Exec(`INSERT INTO "ar_invoices"
	("customer_id", "invoice_number", "currency", "face_value", "counterparty_vat",
	"counterparty_country", "issue_date", "due_date", "status", "company_name")
	VALUES ($1,$2,$3,$4,$5,$6,$7, $8, $9, $10)`,
		inv.GetCustomerID(),
		inv.GetInvoiceNumber(),
		inv.GetCurrency(),
		inv.GetFaceValue(),
		inv.GetCounterPartyVAT(),
		inv.GetCounterPartyCountry(),
		inv.GetIssueDate(),
		inv.GetDueDate(),
		pendingStatus,
		inv.GetCompanyName(),
	)
	return err
}

// insertAP Inserts a new Invoice into AP
func insertAP(tx *sql.Tx, inv *services.InternalInvoice) error {
	_, err := tx.Exec(`INSERT INTO "ap_invoices"
	("customer_id", "invoice_number", "currency", "face_value", "counterparty_vat",
	"counterparty_country", "issue_date", "due_date", "status", "company_name")
	VALUES ($1,$2,$3,$4,$5,$6,$7, $8, $9, $10)`,
		inv.GetCustomerID(),
		inv.GetInvoiceNumber(),
		inv.GetCurrency(),
		inv.GetFaceValue(),
		inv.GetCounterPartyVAT(),
		inv.GetCounterPartyCountry(),
		inv.GetIssueDate(),
		inv.GetDueDate(),
		pendingStatus,
		inv.GetCompanyName())
	return err
}

// updateAP Queries to find the AP Invoice and updates the VAT and Status
func updateAP(tx *sql.Tx, inv *services.InternalInvoice) error {
	var status string
	var vat string
	err := tx.QueryRow(`SELECT status, counterparty_vat FROM "ap_invoices" 
		WHERE invoice_number = $1`, inv.GetInvoiceNumber()).Scan(&status, &vat)

	if err != nil {
		return err
	}

	if status == readyStatus {
		return fmt.Errorf("the invoice %v for VAT %v is ready", inv.GetInvoiceNumber(), inv.GetCounterPartyVAT())
	}

	var currentStatus string
	if inv.GetValidVAT() == true {
		currentStatus = readyStatus
	} else {
		currentStatus = invalidStatus
	}
	_, err = tx.Exec(`UPDATE ap_invoices SET status = $1, counterparty_vat = $2, company_name = $3
		WHERE invoice_number = $4`,
		currentStatus,
		inv.GetCounterPartyVAT(),
		inv.GetCompanyName(),
		inv.GetInvoiceNumber())
	return err
}

// updateAR Queries to find the AR Invoice and updates the VAT and Status
func updateAR(tx *sql.Tx, inv *services.InternalInvoice) error {
	var status string
	var vat string
	row := tx.QueryRow(`SELECT status, counterparty_vat FROM ar_invoices 
		WHERE invoice_number=$1`,
		inv.GetInvoiceNumber())
	err := row.Scan(&status, &vat)
	if err != nil {
		return err
	}

	if status == readyStatus {
		return fmt.Errorf("the invoice %v for VAT %v is ready", inv.GetInvoiceNumber(), inv.GetCounterPartyVAT())
	}

	var currentStatus string
	if inv.GetValidVAT() == true {
		currentStatus = readyStatus
	} else {
		currentStatus = invalidStatus
	}
	_, err = tx.Exec(`UPDATE ar_invoices SET status = $1, counterparty_vat = $2, company_name = $3
		WHERE invoice_number = $4`,
		currentStatus,
		inv.GetCounterPartyVAT(),
		inv.GetCompanyName(),
		inv.GetInvoiceNumber())
	return err
}
