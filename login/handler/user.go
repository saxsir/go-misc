package handler

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/alecthomas/template"
	"github.com/gorilla/csrf"
	"github.com/saxsir/talks/2018/treasure-go/login/model"
	"github.com/saxsir/talks/2018/treasure-go/login/sessions"
)

type User struct {
	DB *sql.DB
}

func (u *User) GetSignupHandler(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("template/signup.tmpl")
	if err != nil {
		return err
	}

	return t.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

func (u *User) PostSignupHandler(w http.ResponseWriter, r *http.Request) error {
	var m model.User
	m.Name = r.PostFormValue("name")
	m.Email = r.PostFormValue("email")

	password := r.PostFormValue("password")

	if err := TXHandler(u.DB, func(tx *sql.Tx) error {
		if _, err := m.Insert(tx, password); err != nil {
			return err
		}
		return tx.Commit()
	}); err != nil {
		return err
	}

	http.Redirect(w, r, "/", 301)
	return nil
}

func (u *User) GetLoginHandler(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("template/login.tmpl")
	if err != nil {
		return err
	}

	return t.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

func (u *User) PostLoginHandler(w http.ResponseWriter, r *http.Request) error {
	m, err := model.Auth(u.DB, r.PostFormValue("email"), r.PostFormValue("password"))
	if err != nil {
		log.Printf("/login: login failed: %s", err)
		sess, _ := sessions.Get(r, "user")
		sess.AddFlash("login failed.")
		if err := sessions.Save(r, w, sess); err != nil {
			log.Printf("/login: save session failed: %s", err)
			return err
		}
		http.Redirect(w, r, "/login", http.StatusFound)
		return nil
	}

	log.Printf("authed: %#v", m)

	sess, _ := sessions.Get(r, "user")
	sess.Values["id"] = m.ID
	sess.Values["email"] = m.Email
	sess.Values["name"] = m.Name
	if err := sessions.Save(r, w, sess); err != nil {
		log.Printf("session can't save: %s", err)
		return err
	}

	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}

func (u *User) GetLogoutHandler(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("template/logout.tmpl")
	if err != nil {
		return err
	}

	return t.Execute(w, map[string]interface{}{
		"currentName":    currentName(r),
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

func (u *User) PostLogoutHandler(w http.ResponseWriter, r *http.Request) error {
	sess, _ := sessions.Get(r, "user")
	if err := sessions.Clear(r, w, sess); err != nil {
		return err
	}
	http.Redirect(w, r, "/", 301)

	return nil
}

// CurrentName returns current user name who logged in.
func currentName(r *http.Request) string {
	if r == nil {
		return ""
	}
	sess, _ := sessions.Get(r, "user")
	rawname, ok := sess.Values["name"]
	if !ok {
		return ""
	}
	name, ok := rawname.(string)
	if !ok {
		return ""
	}

	return name
}
