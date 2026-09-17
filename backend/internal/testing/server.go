package testing

import (
	"time"

	"github.com/VikyCham/go-boilerplate/internal/config"
	"github.com/VikyCham/go-boilerplate/internal/database"
	"github.com/VikyCham/go-boilerplate/internal/server"
	"github.com/rs/zerolog"
)

// CreateTestServer creates a server instance for testing
func CreateTestServer(logger *zerolog.Logger, db *TestDB) *server.Server {

	// Set upobservability config with defaults if not present
	if db.Config.Observability == nil {
		db.Config.Observability = &config.ObservabilityConfig{
			ServiceName: "alfred-test",
			Environment: "test",
			Logging: config.LoggingConfig{
				Level:              "info",
				Format:             "json",
				SlowQueryThreshold: 100 * time.Millisecond,
			},
			NewRelic: config.NewRelicConfig{
				LicenseKey:                "",    //empty for tests
				AppLogForwardingEnabled:   false, // disabled for tests
				DistributedTracingEnabled: false, // disabled for tests
				DebugLogging:              false, // disabled for tests
			},
			HealthChecks: config.HealthCheckConfig{
				Enabled: false,
			},
		}
	}

	testServer := &server.Server{
		Logger: logger,
		DB: &database.Database{
			Pool: db.Pool,
		},
		Config: db.Config,
	}

	return testServer
}
