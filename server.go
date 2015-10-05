package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	store *securecookie.SecureCookie
}

const (
	COOKIE_NAME   string = "PrudonForensics"
	COOKIE_EXPIRE int    = 10 * 60
)

func InitServer(r *mux.Router) {

	s := &Server{securecookie.New([]byte("MiaMySuperSecret"), []byte("MiaMySuperSecret"))}

	r.HandleFunc("/send", s.secure(s.handleSend))
	r.HandleFunc("/login", s.handleLogin)

}

func (s *Server) handleSend(rw http.ResponseWriter, req *http.Request) {

	switch strings.TrimSpace(req.FormValue("type")) {
	case "update":
		s.handleUpdate(rw, req.FormValue("col"), req.FormValue("what"), req.FormValue("value"))
	case "init":
		s.handleInit(rw)
	case "logout":
		s.handleLogout(rw)
	case "change":
		s.handleChange(rw, req.FormValue("nu"), req.FormValue("np"), req.FormValue("ou"), req.FormValue("op"))
	default:
		log.Println("No valid send type given.")
		http.Error(rw, "Internal error", http.StatusInternalServerError)
	}

}

func (s *Server) handleChange(rw http.ResponseWriter, newUser, newPass, oldUser, oldPass string) {

	if !GetUser(oldUser, oldPass) {
		http.Error(rw, "Oude gebruikersnaam/wachtwoord niet goed.", http.StatusInternalServerError)
		return
	}

	SetUser(newUser, newPass)
}

func (s *Server) handleLogout(rw http.ResponseWriter) {
	http.SetCookie(rw, &http.Cookie{
		Name:    COOKIE_NAME,
		Value:   "",
		Path:    "/",
		Expires: time.Now(),
		MaxAge:  -1,
	})
}

func (s *Server) handleInit(rw http.ResponseWriter) {

	d := GetDataJson()

	rw.Header().Set("Content-Length", strconv.Itoa(len(d)))
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(d)

}

func (s *Server) handleUpdate(rw http.ResponseWriter, col, what, value string) {

	i, err := strconv.ParseInt(col, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(rw, "Internal error", http.StatusInternalServerError)
		return
	}

	SetData(i, what == "true", value)
	Render()

	d := GetDataJson()

	rw.Header().Set("Content-Length", strconv.Itoa(len(d)))
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(d)

}

// login
func (s *Server) handleLogin(rw http.ResponseWriter, req *http.Request) {

	u := strings.TrimSpace(req.FormValue("username"))
	if u == "" {
		log.Println("server :: handleLogin :: no valid username")
		http.Error(rw, "Internal error.", http.StatusInternalServerError)
		return
	}

	p := strings.TrimSpace(req.FormValue("password"))
	if p == "" {
		log.Println("server :: handleLogin :: no valid password")
		http.Error(rw, "Internal error.", http.StatusInternalServerError)
		return
	}

	if !GetUser(u, p) {
		http.Error(rw, "Not authorized", http.StatusUnauthorized)
		return
	}

	v := map[string]string{
		"Id": u,
	}

	dv, err := s.store.Encode(COOKIE_NAME, v)
	if err != nil {
		log.Println(err)
		http.Error(rw, "Not authorized", http.StatusUnauthorized)
		return
	}

	s.setCookie(rw, dv)

}

func (s *Server) setCookie(rw http.ResponseWriter, v string) {
	http.SetCookie(rw, &http.Cookie{
		Name:    COOKIE_NAME,
		Value:   v,
		Path:    "/",
		Expires: time.Now().Add(time.Duration(COOKIE_EXPIRE) * time.Second),
		MaxAge:  COOKIE_EXPIRE,
	})
}

func (s *Server) secure(fn http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {

		c, err := req.Cookie(COOKIE_NAME)
		if err != nil {
			log.Println(err)
			http.Error(rw, "Not authorised.", http.StatusUnauthorized)
			return
		}

		var result map[string]string
		err = s.store.Decode(COOKIE_NAME, c.Value, &result)
		if err != nil {
			log.Println(err)
			http.Error(rw, "Not authorised.", http.StatusUnauthorized)
			return
		}

		var ok bool
		if _, ok = result["Id"]; !ok {
			log.Println("server :: secure :: cookie id not found")
			http.Error(rw, "Not authorised.", http.StatusUnauthorized)
			return
		}

		s.setCookie(rw, c.Value)

		fn(rw, req)

	}
}
