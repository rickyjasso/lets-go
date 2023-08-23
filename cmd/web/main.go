package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	// CLI
	// The value of the addr flag is stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are
	// encountered during parsing the application will be terminated.
	flag.Parse()

	// 3 parameters: destination to write the logs to, prefix message, flag to indicate what additional information to include in the log message
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)                  // informational messages
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) // error messages

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// now, the server uses the custom errorLog logger to write log messages instead of the standard logger.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	// The value returned from the flag.String() function is a pointer to the flag
	// value, not the value itself. So we need to dereference the pointer (i.e.
	// prefix it with the * symbol) before using it. Note that we're using the
	// log.Printf() function to interpolate the address with the log message.
	infoLog.Printf("Starting server on %s", *addr) // informational message

	// same thing. because srv is of type *http.Server, ListenAndServe() is a method on that type and can be used.
	err := srv.ListenAndServe()
	// err := http.ListenAndServe(*addr, mux)

	errorLog.Fatal(err) // error message
}
