package httpapi

import (
	"bytes"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/example/distributed-property-graph/internal/algorithm"
	"github.com/example/distributed-property-graph/internal/edge"
	"github.com/example/distributed-property-graph/internal/graph"
	"github.com/example/distributed-property-graph/internal/importer"
	"github.com/example/distributed-property-graph/internal/platform"
	"github.com/example/distributed-property-graph/internal/query"
	"github.com/example/distributed-property-graph/internal/schema"
	"github.com/example/distributed-property-graph/internal/shard"
	"github.com/example/distributed-property-graph/internal/snapshot"
	"github.com/example/distributed-property-graph/internal/vertex"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Server struct {
	graphs       *graph.Service
	schemas      *schema.Service
	vertices     *vertex.Service
	edges        *edge.Service
	queries      *query.Service
	algorithms   *algorithm.Service
	snapshots    *snapshot.Service
	shards       *shard.Service
	token        string
	started      time.Time
	requests     atomic.Uint64
	queryMu      sync.RWMutex
	queryResults map[string]query.Result
	importer     *importer.Importer
}

func NewServer(g *graph.Service, s *schema.Service, v *vertex.Service, e *edge.Service, q *query.Service, a *algorithm.Service, ss *snapshot.Service, sh *shard.Service, imp *importer.Importer, token string) *Server {
	return &Server{graphs: g, schemas: s, vertices: v, edges: e, queries: q, algorithms: a, snapshots: ss, shards: sh, importer: imp, queryResults: map[string]query.Result{}, token: token, started: time.Now()}
}
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.health)
	mux.HandleFunc("/readyz", s.ready)
	mux.HandleFunc("/metrics", s.metrics)
	mux.HandleFunc("/v1/graphs", s.graphsHandler)
	mux.HandleFunc("/v1/graphs/", s.graphAction)
	mux.HandleFunc("/v1/vertices", s.verticesHandler)
	mux.HandleFunc("/v1/edges", s.edgesHandler)
	mux.HandleFunc("/v1/queries", s.queryHandler)
	mux.HandleFunc("/v1/queries/", s.queryResultHandler)
	mux.HandleFunc("/v1/algorithms", s.algorithmHandler)
	mux.HandleFunc("/v1/algorithms/", s.algorithmAction)
	mux.HandleFunc("/v1/snapshots", s.snapshotHandler)
	mux.HandleFunc("/v1/shards", s.shardsHandler)
	return s.middleware(mux)
}
func (s *Server) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path != "/healthz" && r.URL.Path != "/readyz" && r.URL.Path != "/metrics" {
			got := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if subtle.ConstantTimeCompare([]byte(got), []byte(s.token)) != 1 {
				writeError(w, platform.ErrInvalid, http.StatusUnauthorized)
				return
			}
		}
		defer func() {
			if recover() != nil {
				writeError(w, errors.New("internal failure"), 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
func decode(w http.ResponseWriter, r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 8<<20)
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	return d.Decode(v)
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeError(w http.ResponseWriter, err error, status int) {
	writeJSON(w, status, map[string]any{"error": err.Error()})
}
func statusFor(err error) int {
	if errors.Is(err, platform.ErrNotFound) {
		return 404
	}
	if errors.Is(err, platform.ErrConflict) {
		return 409
	}
	return 400
}
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]any{"status": "ok", "uptime": time.Since(s.started).String()})
}
func (s *Server) ready(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ready"})
}
func (s *Server) metrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintf(w, "graph_http_requests_total %d\n", s.requests.Load())
}
func (s *Server) graphsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		items, _ := s.graphs.List(r.Context())
		writeJSON(w, 200, map[string]any{"items": items})
		return
	}
	var in struct {
		Name string `json:"name"`
	}
	if err := decode(w, r, &in); err != nil {
		writeError(w, err, 400)
		return
	}
	g, err := s.graphs.Create(r.Context(), in.Name)
	if err != nil {
		writeError(w, err, 400)
		return
	}
	writeJSON(w, 201, g)
}
func (s *Server) graphAction(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		writeError(w, platform.ErrInvalid, 400)
		return
	}
	id := parts[2]
	if len(parts) == 4 && parts[3] == "publish" {
		g, err := s.graphs.Publish(r.Context(), id)
		if err != nil {
			writeError(w, err, statusFor(err))
			return
		}
		writeJSON(w, 200, g)
		return
	}
	if len(parts) == 4 && parts[3] == "schema" {
		var in schema.Schema
		if err := decode(w, r, &in); err != nil {
			writeError(w, err, 400)
			return
		}
		in.GraphID = id
		out, err := s.schemas.Publish(r.Context(), in)
		if err != nil {
			writeError(w, err, 400)
			return
		}
		writeJSON(w, 200, out)
		return
	}
	if len(parts) == 4 && parts[3] == "mutations" && r.Method == http.MethodPost {
		s.mutations(w, r, id)
		return
	}
	if len(parts) == 4 && parts[3] == "imports" && r.Method == http.MethodPost {
		s.importRecords(w, r, id)
		return
	}
	g, err := s.graphs.Get(r.Context(), id)
	if err != nil {
		writeError(w, err, 404)
		return
	}
	writeJSON(w, 200, g)
}
func (s *Server) verticesHandler(w http.ResponseWriter, r *http.Request) {
	var in struct {
		GraphID    string         `json:"graph_id"`
		ID         string         `json:"id"`
		Type       string         `json:"type"`
		Properties map[string]any `json:"properties"`
		Version    int64          `json:"version"`
	}
	if err := decode(w, r, &in); err != nil {
		writeError(w, err, 400)
		return
	}
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		v, err := s.vertices.Upsert(r.Context(), in.GraphID, in.ID, in.Type, in.Properties, in.Version)
		if err != nil {
			writeError(w, err, statusFor(err))
			return
		}
		writeJSON(w, 201, v)
		return
	}
	v, err := s.vertices.Get(r.Context(), in.GraphID, in.ID)
	if err != nil {
		writeError(w, err, 404)
		return
	}
	writeJSON(w, 200, v)
}
func (s *Server) edgesHandler(w http.ResponseWriter, r *http.Request) {
	var in struct {
		GraphID    string         `json:"graph_id"`
		ID         string         `json:"id"`
		Type       string         `json:"type"`
		FromID     string         `json:"from_id"`
		ToID       string         `json:"to_id"`
		Properties map[string]any `json:"properties"`
		Version    int64          `json:"version"`
	}
	if err := decode(w, r, &in); err != nil {
		writeError(w, err, 400)
		return
	}
	e, err := s.edges.Upsert(r.Context(), in.GraphID, in.ID, in.Type, in.FromID, in.ToID, in.Properties, in.Version)
	if err != nil {
		writeError(w, err, statusFor(err))
		return
	}
	writeJSON(w, 201, e)
}
func (s *Server) queryHandler(w http.ResponseWriter, r *http.Request) {
	var in query.Request
	if err := decode(w, r, &in); err != nil {
		writeError(w, err, 400)
		return
	}
	out, err := s.queries.Traverse(r.Context(), in)
	if err != nil {
		writeError(w, err, 400)
		return
	}
	queryID := platform.NewID("query")
	s.queryMu.Lock()
	s.queryResults[queryID] = out
	s.queryMu.Unlock()
	writeJSON(w, 200, map[string]any{"id": queryID, "result": out})
}

func (s *Server) queryResultHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/queries/")
	s.queryMu.RLock()
	result, ok := s.queryResults[id]
	s.queryMu.RUnlock()
	if !ok {
		writeError(w, platform.ErrNotFound, http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"id": id, "result": result})
}

type mutationRequest struct {
	Vertices []struct {
		ID         string         `json:"id"`
		Type       string         `json:"type"`
		Properties map[string]any `json:"properties"`
		Version    int64          `json:"version"`
	} `json:"vertices"`
	Edges []struct {
		ID         string         `json:"id"`
		Type       string         `json:"type"`
		FromID     string         `json:"from_id"`
		ToID       string         `json:"to_id"`
		Properties map[string]any `json:"properties"`
		Version    int64          `json:"version"`
	} `json:"edges"`
}

func (s *Server) mutations(w http.ResponseWriter, r *http.Request, g string) {
	var in mutationRequest
	if err := decode(w, r, &in); err != nil {
		writeError(w, err, 400)
		return
	}
	for _, v := range in.Vertices {
		if v.ID == "" || v.Type == "" {
			writeError(w, platform.ErrInvalid, 400)
			return
		}
	}
	for _, e := range in.Edges {
		if e.ID == "" || e.Type == "" || e.FromID == "" || e.ToID == "" {
			writeError(w, platform.ErrInvalid, 400)
			return
		}
	}
	for _, v := range in.Vertices {
		if _, err := s.vertices.Upsert(r.Context(), g, v.ID, v.Type, v.Properties, v.Version); err != nil {
			writeError(w, err, statusFor(err))
			return
		}
	}
	for _, e := range in.Edges {
		if _, err := s.vertices.Get(r.Context(), g, e.FromID); err != nil {
			writeError(w, err, statusFor(err))
			return
		}
		if _, err := s.vertices.Get(r.Context(), g, e.ToID); err != nil {
			writeError(w, err, statusFor(err))
			return
		}
		if _, err := s.edges.Upsert(r.Context(), g, e.ID, e.Type, e.FromID, e.ToID, e.Properties, e.Version); err != nil {
			writeError(w, err, statusFor(err))
			return
		}
	}
	writeJSON(w, 200, map[string]any{"vertices": len(in.Vertices), "edges": len(in.Edges), "status": "committed"})
}
func (s *Server) importRecords(w http.ResponseWriter, r *http.Request, g string) {
	var in struct {
		Records   []importer.Record `json:"records"`
		StartLine int               `json:"start_line"`
	}
	if err := decode(w, r, &in); err != nil {
		writeError(w, err, 400)
		return
	}
	var source bytes.Buffer
	enc := json.NewEncoder(&source)
	for _, record := range in.Records {
		if err := enc.Encode(record); err != nil {
			writeError(w, err, 400)
			return
		}
	}
	summary, err := s.importer.Import(r.Context(), g, &source, in.StartLine)
	if err != nil {
		writeError(w, err, statusFor(err))
		return
	}
	writeJSON(w, 202, summary)
}
func (s *Server) algorithmHandler(w http.ResponseWriter, r *http.Request) {
	var in struct {
		GraphID    string         `json:"graph_id"`
		Name       string         `json:"name"`
		Parameters map[string]any `json:"parameters"`
	}
	if err := decode(w, r, &in); err != nil {
		writeError(w, err, 400)
		return
	}
	task, err := s.algorithms.Submit(r.Context(), in.GraphID, in.Name, in.Parameters)
	if err != nil {
		writeError(w, err, 400)
		return
	}
	writeJSON(w, 202, task)
}
func (s *Server) algorithmAction(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		writeError(w, platform.ErrInvalid, 400)
		return
	}
	task, err := s.algorithms.Get(r.Context(), parts[2])
	if err != nil {
		writeError(w, err, 404)
		return
	}
	writeJSON(w, 200, task)
}
func (s *Server) snapshotHandler(w http.ResponseWriter, r *http.Request) {
	var in struct {
		GraphID    string `json:"graph_id"`
		TTLSeconds int    `json:"ttl_seconds"`
	}
	if err := decode(w, r, &in); err != nil {
		writeError(w, err, 400)
		return
	}
	ttl := time.Duration(in.TTLSeconds) * time.Second
	if ttl <= 0 {
		ttl = time.Hour
	}
	snap, err := s.snapshots.Create(r.Context(), in.GraphID, ttl)
	if err != nil {
		writeError(w, err, 400)
		return
	}
	writeJSON(w, 201, snap)
}
func (s *Server) shardsHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"items": s.shards.List(r.Context())})
}
