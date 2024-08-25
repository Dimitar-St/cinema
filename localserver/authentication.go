package localserver

import (
	"cinema/db"
	"cinema/db/models"
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func init() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
	}

	gob.Register(models.User{})
}

func Signup(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)

	user.Password = string(hashedPassword)

	println(user.Username)
	println(user.Password)

	db.DB.Where(models.User{Username: user.Username}).FirstOrCreate(user)
}

func Profile(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := getUser(session)

	if auth := user.Authenticated; !auth {
		session.AddFlash("You don't have access!")
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return
	}

	pageToLoad, _ := os.ReadFile("./build/user.html")

	render(w, r, parseTemplate(string(pageToLoad)))
}

func getUser(s *sessions.Session) models.User {
	val := s.Values["user"]
	var user = models.User{}
	user, ok := val.(models.User)
	if !ok {
		return models.User{Authenticated: false}
	}
	user.Authenticated = true
	return user
}

func Signin(w http.ResponseWriter, r *http.Request) {
	creds := &models.User{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result := &models.User{}
	db.DB.Model(models.User{}).First(result, "Username = ?", creds.Username)
	if result == nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(creds.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	session, err := store.Get(r, "cookie-name")
	if err != nil {
		println(result.Username)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["user"] = result

	err = session.Save(r, w)
	if err != nil {
		println("save")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	println("asdasdsadasdas")
	http.Redirect(w, r, "/profile", http.StatusFound)

}
