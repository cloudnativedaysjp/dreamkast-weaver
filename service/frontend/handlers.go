package frontend

import (
	"fmt"
	"net/http"
)

func (fe *Server) voteHandler(w http.ResponseWriter, r *http.Request) {
	reversed, _ := fe.cfpSvc.Vote(r.Context(), "!dlroW ,olleH")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, reversed)
}
