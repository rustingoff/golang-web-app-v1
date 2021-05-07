package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rustingoff/models"
	"github.com/rustingoff/sessions"
	"github.com/rustingoff/utils"
)

func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	updates, err := models.GetAllUpdates()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	utils.ExecuteTemplate(w, "index.html", struct {
		Title       string
		Updates     []*models.Update
		DisplayForm bool
	}{
		Title:       "All updates",
		Updates:     updates,
		DisplayForm: true,
	})
}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "token")
	untypedUserId := session.Values["user_id"]
	userId, ok := untypedUserId.(int64)
	if !ok {
		utils.InternalServerError(w)
		return
	}
	r.ParseForm()
	body := r.PostForm.Get("update")
	err := models.PostUpdate(userId, body)
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "login.html", nil)
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	user, err := models.AuthenticatedUser(username, password)
	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			utils.ExecuteTemplate(w, "login.html", "unknown user")
		case models.ErrInvalidLogin:
			utils.ExecuteTemplate(w, "login.html", "invalid login")
		default:
			utils.InternalServerError(w)
		}
		return
	}
	userId, err := user.GetId()
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	session, _ := sessions.Store.Get(r, "token")
	session.Values["user_id"] = userId
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func logoutGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "token")
	delete(session.Values, "user_id")
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func registerGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "register.html", nil)
}

func registerPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	err := models.RegisterUser(username, password)
	if err == models.ErrUsernameTaken {
		utils.ExecuteTemplate(w, "register.html", "username taken")
	} else if err != nil {
		utils.InternalServerError(w)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}

func userGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "token")
	untypedUserId := session.Values["user_id"]
	currentUserId, ok := untypedUserId.(int64)
	if !ok {
		utils.InternalServerError(w)
		return
	}
	vars := mux.Vars(r)
	username := vars["username"]
	user, err := models.GetUserByUsername(username)
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	userId, err := user.GetId()
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	updates, err := models.GetUpdates(userId)
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	utils.ExecuteTemplate(w, "index.html", struct {
		Title       string
		Updates     []*models.Update
		DisplayForm bool
	}{
		Title:       username,
		Updates:     updates,
		DisplayForm: currentUserId == userId,
	})
}
