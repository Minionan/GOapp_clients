// code/database/clients.go
package database

import (
	"encoding/json"
	"fmt"
	"strings"
)

// GetClients retrieves all clients from the database
func GetClients() ([]Client, error) {
	rows, err := db.Query(`
        SELECT id, clientName, parentName, phone, email, abbreviation, active, invoice_lock
        FROM clients`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []Client
	for rows.Next() {
		var c Client
		err := rows.Scan(&c.ID, &c.ClientName, &c.ParentName,
			&c.Phone, &c.Email, &c.Abbreviation, &c.Active, &c.Invoice_lock)
		if err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	return clients, nil
}

// GetClientByID retrieves a specific client by their ID
func GetClientByID(clientID int) (*Client, error) {
	var client Client
	err := db.QueryRow(`
        SELECT id, clientName, parentName, address1, address2, phone, email, abbreviation, active, invoice_lock
        FROM clients WHERE id = ?`, clientID).
		Scan(&client.ID, &client.ClientName, &client.ParentName,
			&client.Address1, &client.Address2, &client.Phone, &client.Email, &client.Abbreviation, &client.Active, &client.Invoice_lock)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// CreateClient creates a new client in the database
func CreateClient(client Client) error {
	_, err := db.Exec(`
        INSERT INTO clients (clientName, parentName, address1, address2, phone, email, abbreviation, active, invoice_lock)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		client.ClientName, client.ParentName, client.Address1,
		client.Address2, client.Phone, client.Email, client.Abbreviation, client.Active, client.Invoice_lock)
	return err
}

// UpdateClient updates an existing client in the database
func UpdateClient(client Client) error {
	_, err := db.Exec(`
        UPDATE clients
        SET clientName=?, parentName=?, address1=?, address2=?, phone=?, email=?, abbreviation=?, active=?, invoice_lock=?
        WHERE id=?`,
		client.ClientName, client.ParentName, client.Address1, client.Address2,
		client.Phone, client.Email, client.Abbreviation, client.Active, client.Invoice_lock, client.ID)
	return err
}

// DeleteClient deletes a client from the database
func DeleteClient(clientID int) error {
	_, err := db.Exec("DELETE FROM clients WHERE id = ?", clientID)
	return err
}

// ExportSelectedClientsToJSON exports selected clients to a JSON file
func ExportSelectedClientsToJSON(clientIDs []int) ([]byte, error) {
	if len(clientIDs) == 0 {
		return nil, fmt.Errorf("no clients selected")
	}

	placeholders := make([]string, len(clientIDs))
	args := make([]interface{}, len(clientIDs))

	for i := range clientIDs {
		placeholders[i] = "?"
		args[i] = clientIDs[i]
	}

	query := fmt.Sprintf(`
        SELECT id, clientName, parentName, address1, address2, 
               phone, email, abbreviation, active, invoice_lock
        FROM clients 
        WHERE id IN (%s)`, strings.Join(placeholders, ","))

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []Client
	for rows.Next() {
		var c Client
		err := rows.Scan(&c.ID, &c.ClientName, &c.ParentName,
			&c.Address1, &c.Address2, &c.Phone, &c.Email,
			&c.Abbreviation, &c.Active, &c.Invoice_lock)
		if err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}

	return json.Marshal(clients)
}

// ImportClientsFromJSON imports clients from a JSON file
func ImportClientsFromJSON(data []byte) error {
	var clients []Client
	if err := json.Unmarshal(data, &clients); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, client := range clients {
		client.ID = 0 // Reset ID for new insertion
		_, err := tx.Exec(`
            INSERT INTO clients (clientName, parentName, address1, address2, 
                               phone, email, abbreviation, active, invoice_lock)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			client.ClientName, client.ParentName, client.Address1,
			client.Address2, client.Phone, client.Email, client.Abbreviation,
			client.Active, client.Invoice_lock)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
