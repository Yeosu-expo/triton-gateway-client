package handler

import "net/http"

/* HTML, CSS, JS 요청 */
func (h *Handler) homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "index.html", http.StatusTemporaryRedirect)
}
