CREATE DATABASE arex;

-- SQL Script for Creating AP and AR invoice tables
CREATE TABLE IF NOT EXISTS ap_invoices (
    customer_id UUID NOT NULL,
    invoice_number INTEGER NOT NULL,
    currency TEXT NOT NULL,
    face_value INTEGER NOT NULL,
    counterparty_vat TEXT NOT NULL,
    counterparty_country TEXT NOT NULL,
    issue_date TEXT NOT NULL,
    due_date TEXT NOT NULL,
    status TEXT NOT NULL,
    company_name TEXT NOT NULL,

    PRIMARY KEY (counterparty_vat, invoice_number)
);

CREATE TABLE IF NOT EXISTS ar_invoices (
    customer_id UUID NOT NULL, 
    invoice_number INTEGER NOT NULL,
    currency TEXT NOT NULL,
    face_value INTEGER NOT NULL,
    counterparty_vat TEXT NOT NULL,
    counterparty_country TEXT NOT NULL,
    issue_date TEXT NOT NULL,
    due_date TEXT NOT NULL,
    status TEXT NOT NULL,
    company_name TEXT NOT NULL,
    
    PRIMARY KEY (customer_id, invoice_number)
);