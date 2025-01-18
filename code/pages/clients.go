// code/pages/clients.go
package pages

import (
	"BBCapp/code/database"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func ClientsHandler(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clients, err := database.GetClients()
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		data := PageData{
			Name:    "Admin", // Default name, since no user authentication
			Clients: clients,
		}

		tpl.ExecuteTemplate(w, "clients.html", data)
	}
}

func ClientAddHandler(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// Parse form data
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Error parsing form data", http.StatusBadRequest)
				return
			}

			// Create a new Client struct from the form data
			client := database.Client{
				ClientName:   r.FormValue("clientName"),
				ParentName:   r.FormValue("parentName"),
				Address1:     r.FormValue("address1"),
				Address2:     r.FormValue("address2"),
				Phone:        r.FormValue("phone"),
				Email:        r.FormValue("email"),
				Abbreviation: r.FormValue("abbreviation"),
				Active:       r.FormValue("active") == "1",
				Invoice_lock: false, // Default value
			}

			// Insert the new client into the database
			err = database.CreateClient(client)
			if err != nil {
				log.Printf("Error adding client to database: %v", err)
				tpl.ExecuteTemplate(w, "clientAdd.html", map[string]string{"Error": "Error adding client to database"})
				return
			}

			// Redirect to the clients page
			http.Redirect(w, r, "/clients", http.StatusSeeOther)
			return
		}

		// Render the add client form
		tpl.ExecuteTemplate(w, "clientAdd.html", nil)
	}
}

func ClientEditHandler(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		clientID, err := strconv.Atoi(parts[3])
		if err != nil {
			http.Error(w, "Invalid client ID", http.StatusBadRequest)
			return
		}

		if r.Method == "GET" {
			client, err := database.GetClientByID(clientID)
			if err != nil {
				http.Error(w, "Client not found", http.StatusNotFound)
				return
			}
			data := ClientEditData{Client: *client}
			tpl.ExecuteTemplate(w, "clientEdit.html", data)
			return
		}

		if r.Method == "POST" {
			client := database.Client{
				ID:           clientID,
				ClientName:   r.FormValue("clientName"),
				ParentName:   r.FormValue("parentName"),
				Address1:     r.FormValue("address1"),
				Address2:     r.FormValue("address2"),
				Phone:        r.FormValue("phone"),
				Email:        r.FormValue("email"),
				Abbreviation: r.FormValue("abbreviation"),
				Active:       r.FormValue("active") == "1",
			}

			if err := database.UpdateClient(client); err != nil {
				data := ClientEditData{Error: "Error updating client"}
				tpl.ExecuteTemplate(w, "clientEdit.html", data)
				return
			}
			http.Redirect(w, r, "/clients", http.StatusSeeOther)
		}
	}
}

func ClientDeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		clientID, err := strconv.Atoi(parts[3])
		if err != nil {
			http.Error(w, "Invalid client ID", http.StatusBadRequest)
			return
		}

		if err := database.DeleteClient(clientID); err != nil {
			http.Error(w, "Error deleting client", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func ClientExportHandler(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			clients, err := database.GetClients() // Use GetClients() instead of GetClientsByUserID
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}

			data := PageData{
				Name:    "Admin", // Assuming a default name
				Clients: clients,
			}
			tpl.ExecuteTemplate(w, "clientExport.html", data)

		case "POST":
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Form parsing error", http.StatusBadRequest)
				return
			}

			selectedIDs := r.PostForm["selectedClients"]
			clientIDs := make([]int, 0, len(selectedIDs))

			for _, idStr := range selectedIDs {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					http.Error(w, "Invalid client ID", http.StatusBadRequest)
					return
				}
				clientIDs = append(clientIDs, id)
			}

			data, err := database.ExportSelectedClientsToJSON(clientIDs) // Remove userID argument
			if err != nil {
				http.Error(w, "Export failed", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Disposition", "attachment; filename=clients.json")
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func ClientImportHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "File upload failed", http.StatusBadRequest)
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "File read failed", http.StatusInternalServerError)
			return
		}

		if err := database.ImportClientsFromJSON(data); err != nil { // Remove userID argument
			http.Error(w, "Import failed", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/clients", http.StatusSeeOther)
	}
}
