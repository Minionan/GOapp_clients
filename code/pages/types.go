// code/pages/types.go
package pages

import "BBCapp/code/database"

type PageData struct {
	Name    string
	Clients []database.Client
	Error   string
}

type ClientEditData struct {
	Client database.Client
	Error  string
}
