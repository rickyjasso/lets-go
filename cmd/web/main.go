package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql" // (the blank identifier means that the package is imported for its side-effects only)
	"snippetbox.rickyjasso.com/internal/models"
)

// create a new struct type to hold the application-wide dependencies for the web application.
// the scope of this struct is the entire application, so we define it at the top-level of the main.go file.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
}

func main() {
	// CLI
	// The value of the addr flag is stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	dsn := flag.String("dsn", "web:9091@/snippetbox?parseTime=true", "MySQL data source name")

	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are
	// encountered during parsing the application will be terminated.
	flag.Parse()

	// 3 parameters: destination to write the logs to, prefix message, flag to indicate what additional information to include in the log message
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)                  // informational messages
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) // error messages

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	// the app has to be a reference to the application struct type so that we can assign the loggers to the errorLog and infoLog fields.
	// if it was a copy, the changes would not be reflected in the main function.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
	}

	// now, the server uses the custom errorLog logger to write log messages instead of the standard logger.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// The value returned from the flag.String() function is a pointer to the flag
	// value, not the value itself. So we need to dereference the pointer (i.e.
	// prefix it with the * symbol) before using it. Note that we're using the
	// log.Printf() function to interpolate the address with the log message.
	infoLog.Printf("Starting server on %s", *addr) // informational message

	// same thing. because srv is of type *http.Server, ListenAndServe() is a method on that type and can be used.
	err = srv.ListenAndServe()
	// err := http.ListenAndServe(*addr, mux)

	errorLog.Fatal(err) // error message
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		// if there is an error, we return a nil *sql.DB connection pool and the error.
		return nil, err
	}
	return db, nil
}
