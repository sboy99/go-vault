package http

import (
	nethttp "net/http"
	"strconv"
)

func (r *router) handleGetJob(w nethttp.ResponseWriter, req *nethttp.Request) {
	id := req.PathValue("id")
	j, err := r.jobs.Get(req.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeOK(w, j)
}

func (r *router) handleListJobs(w nethttp.ResponseWriter, req *nethttp.Request) {
	limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(req.URL.Query().Get("offset"))
	jobs, err := r.jobs.List(req.Context(), limit, offset)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeOK(w, jobs)
}
