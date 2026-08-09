package router

import (
	"database/sql"
	"net/http"
)

func MainRoute(mux *http.ServeMux, db *sql.DB) {
	AuthRoute(mux, db)
}
