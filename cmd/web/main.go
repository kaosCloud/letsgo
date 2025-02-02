package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"log/slog"
	"net/http"
	"os"
	"snippetbox.kaos.cloud/internal/models"
)

type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func main() {
	// Define a new command-line flag with name "addr" and default value of :4000.
	addr := flag.String("addr", ":4000", "HTTP network address")
	// And one for MYSQL DSN string.
	dsn := flag.String("dsn", "web:mysql@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	// Initialize a new structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	// Initialize a new instance of our application struct, containing the
	// dependencies.
	app := &application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db},
	}

	// The value returned from flag.String() function is pointer to the flag value,
	// not the value itself. That means that addr variable is actually a pointer, and we need
	// to dereference it before using it.
	logger.Info("starting server", "addr", *addr)
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)

}
