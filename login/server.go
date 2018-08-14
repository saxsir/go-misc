package login

import (
	"database/sql"
	"log"
	"net/http"

	"io"

	"github.com/gorilla/context"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"

	_ "github.com/mattn/go-sqlite3"
	"github.com/saxsir/talks/2018/treasure-go/login/handler"
)

// Server is whole server implementation for this wiki app.
// This holds database connection and router settings.
type Server struct {
	db     *sql.DB
	router *mux.Router
}

func NewServer() *Server {
	return &Server{}
}

// Init initialize server state. Connecting to database, compiling templates,
// and settings router.
func (s *Server) Init() {
	db, err := sql.Open("sqlite3", "dev.db")
	if err != nil {
		log.Fatalf("db initialization failed: %s", err)
	}
	s.db = db
	s.router = s.Route()
}

// Close makes the database connection to close.
func (s *Server) Close() error {
	return s.db.Close()
}

// csrfProtectKey should have 32 byte length.
var csrfProtectKey = []byte("32-byte-long-auth-key")

// Run starts running http server.
func (s *Server) Run(addr string) {
	log.Printf("start listening on %s", addr)

	// NOTE: when you serve on TLS, make csrf.Secure(true)
	CSRF := csrf.Protect(csrfProtectKey, csrf.Secure(false))
	http.ListenAndServe(addr, context.ClearHandler(CSRF(s.router)))
}

func (s *Server) Route() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "pong")
	}).Methods("GET")

	u := &handler.User{DB: s.db}
	r.Handle("/signup", handler.Handler(u.GetSignupHandler)).Methods("GET")
	r.Handle("/signup", handler.Handler(u.PostSignupHandler)).Methods("POST")
	r.Handle("/login", handler.Handler(u.GetLoginHandler)).Methods("GET")
	r.Handle("/login", handler.Handler(u.PostLoginHandler)).Methods("POST")
	r.Handle("/logout", handler.Handler(u.GetLogoutHandler)).Methods("GET")
	r.Handle("/logout", handler.Handler(u.PostLogoutHandler)).Methods("POST")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})

	return r
}
