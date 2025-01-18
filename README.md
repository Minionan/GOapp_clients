# GOapp_clients

This is a simple client-handling app written in GO.
It is a good start when creating your web app where clients functionality is required.
This version of the app can add clients, remove clients, exprort and import selected clients list.

## Setup

Clone the repository with `git clone https://github.com/Minionan/GOapp_invoice.git`

### Initialising user database

1. Ensure that there is a db folder in main project folder, if not create one yourself
2. Run `init_db.go` script by typing in terminal `go run init_db.go`
3. Verify if a new SQLite database file was created in db folder
4. The default data.db file has no records, you will have to create them manually

## Run app

1. Run `go mod tidy`
2. Run `go run main.go`
