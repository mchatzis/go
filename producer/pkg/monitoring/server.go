package monitoring

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/mchatzis/go/producer/pkg/migrate"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartServer(port string) {
	http.Handle("/metrics", promhttp.Handler())

	migrator, err := migrate.NewMigrator()
	if err != nil {
		logger.Fatalf("failed to create migrator: %v", err)
	}
	defer migrator.Close()
	registerMigratorEndpoints(migrator)

	logger.Infof("Starting server on port %s", port)
	http.ListenAndServe(":"+port, nil)
}

func registerMigratorEndpoints(migrator *migrate.Migrator) {
	http.HandleFunc("/migrate/up", func(w http.ResponseWriter, r *http.Request) {
		if err := migrator.Up(); err != nil {
			http.Error(w, fmt.Sprintf("Migration up failed: %v", err), http.StatusInternalServerError)
			logger.Errorf("Migration up failed: %v", err)
			return
		}
		fmt.Fprintf(w, "Migration up successful")
	})

	http.HandleFunc("/migrate/down", func(w http.ResponseWriter, r *http.Request) {
		if err := migrator.Down(); err != nil {
			http.Error(w, fmt.Sprintf("Migration down failed: %v", err), http.StatusInternalServerError)
			logger.Errorf("Migration down failed: %v", err)
			return
		}
		fmt.Fprintf(w, "Migration down successful")
	})
}
