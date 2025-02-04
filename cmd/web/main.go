package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/marcusgeorgievski/snippetbox/internal/models"
)

type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	// Flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name (dsn)")
	flag.Parse()

	// Database connection
	db, err := OpenDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	// Template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Form decoder
	formDecoder := form.NewDecoder()

	// Session manager
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.NewWithCleanupInterval(db, 12 * time.Hour)


	// Application struct
	app := &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
		sessionManager: sessionManager,
	}

	// Start server

	logger.Info(fmt.Sprintf("starting server at http://localhost%s", *addr), slog.String("addr", *addr))

	err = http.ListenAndServe(*addr, app.routes())

	logger.Error(err.Error())
	os.Exit(1)
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", "web:pass@/snippetbox?parseTime=true")
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
