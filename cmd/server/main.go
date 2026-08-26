package main

import (
	"context"
	"github.com/example/distributed-property-graph/internal/algorithm"
	"github.com/example/distributed-property-graph/internal/edge"
	"github.com/example/distributed-property-graph/internal/graph"
	"github.com/example/distributed-property-graph/internal/httpapi"
	"github.com/example/distributed-property-graph/internal/importer"
	"github.com/example/distributed-property-graph/internal/platform"
	"github.com/example/distributed-property-graph/internal/query"
	"github.com/example/distributed-property-graph/internal/schema"
	"github.com/example/distributed-property-graph/internal/shard"
	"github.com/example/distributed-property-graph/internal/snapshot"
	"github.com/example/distributed-property-graph/internal/vertex"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := platform.LoadConfig()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	clock := platform.SystemClock{}
	graphRepo := graph.NewMemoryRepository()
	graphSvc := graph.NewService(graphRepo, clock)
	schemaSvc := schema.NewService(schema.NewMemoryRepository())
	vertexRepo := vertex.NewMemoryRepository()
	vertexSvc := vertex.NewService(vertexRepo, clock)
	edgeRepo := edge.NewMemoryRepository()
	edgeSvc := edge.NewService(edgeRepo, clock)
	importSvc := importer.New(vertexSvc, edgeSvc, 100000)
	querySvc := query.NewService(vertexRepo, edgeRepo)
	algorithmSvc := algorithm.NewService(vertexRepo, edgeRepo, clock)
	snapshotSvc := snapshot.NewService(clock)
	shardSvc := shard.NewService(cfg.ShardCount, clock)
	algorithmSvc.Run(context.Background(), 2)
	handler := httpapi.NewServer(graphSvc, schemaSvc, vertexSvc, edgeSvc, querySvc, algorithmSvc, snapshotSvc, shardSvc, importSvc, cfg.AuthToken).Handler()
	server := &http.Server{Addr: cfg.HTTPAddr, Handler: handler, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 20 * time.Second, WriteTimeout: 20 * time.Second}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	go func() {
		logger.Info("server started", "address", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			cancel()
		}
	}()
	<-ctx.Done()
	shutdown, stop := context.WithTimeout(context.Background(), 10*time.Second)
	defer stop()
	_ = server.Shutdown(shutdown)
	logger.Info("server stopped")
}
