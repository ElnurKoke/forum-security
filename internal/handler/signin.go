package handler

import (
	"forum/internal/models"
	"log"
	"net/http"
)

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signin" {
		h.ErrorPage(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		username := r.FormValue("username")
		password := r.FormValue("password")

		token, expired, err := h.Service.CheckUser(models.User{
			Username: username,
			Password: password,
		})
		if err != nil {
			info := models.InfoSign{
				Error:    err.Error(),
				Username: username,
				Password: password,
			}
			if err := h.Temp.ExecuteTemplate(w, "signin.html", info); err != nil {
				h.ErrorPage(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			models.InfoLog.Printf("URL: %s\n        Method:   %s\n        Message:  %s\n        Status:   %s\n", r.URL.Path, r.Method, err, "fail")
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   token,
			Path:    "/",
			Expires: expired,
		})
		models.InfoLog.Printf("URL: %s\n        Method:   %s\n        Message:  %s\n        Status:   %s\n", r.URL.Path, r.Method, username, "successful")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	case http.MethodGet:

		if err := h.Temp.ExecuteTemplate(w, "signin.html", nil); err != nil {
			h.ErrorPage(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	default:
		h.ErrorPage(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}
