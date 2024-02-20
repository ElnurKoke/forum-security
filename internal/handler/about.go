package handler

import (
	"net/http"
)

func (h *Handler) info(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/about" {
		h.ErrorPage(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err := h.Temp.ExecuteTemplate(w, "about.html", nil); err != nil {
		h.ErrorPage(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
