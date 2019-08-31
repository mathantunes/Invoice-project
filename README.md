
# arex_project

Testing the API

Initialize the containers by running docker-compose up
Run creatingTables bat to initialize the PostgreSQL database

Inside the Client folder, there are test files for each functionality:

1) Uploader: From the uploader folder
    Run "go test -run TestCreateXMLInvoice" to test the invoice creation.
    This test utilizes the "current_invoice.xml" file to upload.
    The Extra parameters are defined in the Beginning of the function

    Run "go test -run TestUpdateInvoicePreview" to test the Invoice Preview updload.
    This test utilizes the "current_invoice_preview.pdf" file to upload.
    The Extra parameters are defined in the Beginning of the function

    Run "go test -run TestUpdateAttachment" to test the Invoice Attachment upload.
    This test utilizes the "current_invoice_attachment.pdf" file to upload
    The Extra parameters are defined in the Beginning of the function

2) Getter: from the getter folder
    Run "go test -run TestGetInvoicePreview" to test the Invoice Preview Download.
    This test outputs the file to "testdata" folder on "current_get_invoice_preview.pdf" file.
    The Extra parameters are defined in the Beginning of the function

    Run "go test -run TestGetAttachments" to test the Invoice Attachments Listing.
    This test outputs the filenames for a certain InvoiceNumber.
    The Extra parameters are defined in the Beginning of the function

    Run "go test -run TestGetAttachment" to test the Invoice Attachment Download.
    This test outputs the files to "testdata" folder on "get_attachment.pdf" file.
    The Extra parameters are defined in the Beginning of the function


The application utilizes ElastiqMQ to mock SQS using the AWS SDK.
The application utilizes Localstack to mock s3 using the AWS SDK.

The application is made by:
1) (invoice_creator) GRPC Server application container that receives incoming messages and dispatch either to SQS, s3 and VIES API.
2) (storage) PSQL application container that receives SQS messages and stores on PostgreSQL.
3) PostgreSQL container
4) Localstack container
5) ElastiqMQ container
