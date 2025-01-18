// code/database/models.go
package database

// Client represents a client in the system
type Client struct {
	ID           int
	UserID       int
	ClientName   string
	ParentName   string
	Address1     string
	Address2     string
	Phone        string
	Email        string
	Abbreviation string
	Active       bool
	Invoice_lock bool
}
