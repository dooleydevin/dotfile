package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/knoebber/dotfile/db"
)

func setupRoutes(r *mux.Router, secure bool) {
	staticRoutes(r, secure)
	assetRoutes(r)
	dynamicRoutes(r)
}

// Pages that do not have a dynamic route element.
// These conflict with the username wild card.
// Seed the DB on start to prevent these from being registered.
// Important that these routes are setup first so that registerRoutes works as expected.
func staticRoutes(r *mux.Router, secure bool) {
	r.HandleFunc("/", getIndexHandler())
	r.HandleFunc("/about", getAboutHandler())
	r.HandleFunc("/explore", getExploreHandler())
	r.HandleFunc("/signup", signupHandler(secure))
	r.HandleFunc("/login", loginHandler(secure))
	r.HandleFunc("/logout", logoutHandler())
	r.HandleFunc("/email", emailHandler())
	r.HandleFunc("/password", passwordHandler())
	r.HandleFunc("/new_file", createFileHandler())
	registerRoutes(r)
}

func assetRoutes(r *mux.Router) {
	assets := http.FileSystem(http.Dir("assets/"))
	r.Path("/style.css").Handler(http.FileServer(assets))
	r.Path("/favicon.ico").Handler(http.FileServer(assets))
}

func dynamicRoutes(r *mux.Router) {
	r.HandleFunc("/new_file/{alias}", confirmFileHandler())

	r.HandleFunc("/{username}", userHandler())
	r.HandleFunc("/{username}/{alias}", fileHandler())
	r.HandleFunc("/{username}/{alias}/commits", commitsHandler())
	r.HandleFunc("/{username}/{alias}/{hash}", commitHandler())
}

func registerRoutes(r *mux.Router) {
	reserved := []interface{}{}
	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		reserved = append(reserved, pathTemplate[1:])
		return nil
	})

	if err != nil {
		log.Fatalf("walking routes: %s", err)
	}

	if err := db.SeedReservedUsernames(reserved); err != nil {
		log.Fatal(err)
	}
}
