package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"os"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/plan"
	"sigs.k8s.io/external-dns/provider"
)

type API struct {
	chi.Router
	provider  provider.Provider
	hostsFile string
}

func New(p provider.Provider, hostsFile string) *API {
	a := &API{
		chi.NewRouter(),
		p,
		hostsFile,
	}

	a.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	a.Get("/", a.negotiate)
	a.Get("/records", a.records)
	a.Post("/records", a.applyChanges)
	a.Post("/adjustendpoints", a.adjustEndpoints)

	a.Get("/hosts", a.hosts)
	return a
}

func (a API) hosts(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(a.hostsFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, err = io.Copy(w, f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a API) negotiate(w http.ResponseWriter, r *http.Request) {
	b, err := a.provider.GetDomainFilter().MarshalJSON()
	if err != nil {
		fmt.Printf("failed to marshal domain filter, request method: %s, request path: %s\n", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(contentTypeHeader, string(mediaTypeVersion1))
	if _, writeError := w.Write(b); writeError != nil {
		fmt.Printf("error writing response: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a API) records(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	records, err := a.provider.Records(ctx)
	if err != nil {
		fmt.Printf("error getting records: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentTypeHeader, string(mediaTypeVersion1))
	w.Header().Set(varyHeader, contentTypeHeader)
	err = json.NewEncoder(w).Encode(records)
	if err != nil {
		fmt.Printf("error encoding records: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a API) applyChanges(w http.ResponseWriter, r *http.Request) {
	var changes plan.Changes
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).Decode(&changes); err != nil {
		w.Header().Set(contentTypeHeader, contentTypePlaintext)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("error decoding changes: %s\n", err.Error())

		return
	}

	if err := a.provider.ApplyChanges(ctx, &changes); err != nil {
		w.Header().Set(contentTypeHeader, contentTypePlaintext)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error applying changes: %s\n", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a API) adjustEndpoints(w http.ResponseWriter, r *http.Request) {
	var pve []*endpoint.Endpoint
	if err := json.NewDecoder(r.Body).Decode(&pve); err != nil {
		w.Header().Set(contentTypeHeader, contentTypePlaintext)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("failed to decode request body: %v\n", err)
		return
	}
	pve, err := a.provider.AdjustEndpoints(pve)
	if err != nil {
		fmt.Printf("error adjusting endpoints %v\n", err)
		return
	}
	out, _ := json.Marshal(&pve)
	w.Header().Set(contentTypeHeader, string(mediaTypeVersion1))
	w.Header().Set(varyHeader, contentTypeHeader)
	if _, writeError := fmt.Fprint(w, string(out)); writeError != nil {
		fmt.Printf("error writing response %v\n", writeError)
		return
	}
}
