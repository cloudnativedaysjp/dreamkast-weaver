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

func (fe *Server) showHandler(w http.ResponseWriter, r *http.Request) {
	res, _ := fe.cfpSvc.Show(r.Context())

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, res[len(res)-1].Dt.String())
}
