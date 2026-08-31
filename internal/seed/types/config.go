package types

// SeedConfig defines the configuration options for the data-generation subsystem.
type SeedConfig struct {
	Environment   string `kong:"default='development',env='FINCHER_ENV'"`
	TursoURL      string `kong:"default='fincher.db',env='FINCHER_TURSO_URL'"`
	TursoToken    string `kong:"env='FINCHER_TURSO_TOKEN'"`
	ClickHouseDSN string `kong:"default='127.0.0.1:9000',env='FINCHER_CLICKHOUSE_DSN'"`

	// Seeder-specific flags
	Seed            int64 `kong:"default=20260828,env='FINCHER_SEED',help='deterministic RNG seed'"`
	EventsPerVendor int   `kong:"default=10000,env='FINCHER_SEED_EVENTS_PER_VENDOR'"`
	HistoryDays     int   `kong:"default=120,env='FINCHER_SEED_HISTORY_DAYS'"`
	Titles          int   `kong:"default=7,help='total titles incl hero'"`
	FillerVendors   int   `kong:"default=0,help='additional faker vendors added to curated set'"`
	Reset           bool  `kong:"help='TRUNCATE existing Turso rows + CH fincher.events before seeding'"`
	BatchSize       int   `kong:"default=2000,help='CH insert chunk size'"`
}

// DefaultConfig returns a SeedConfig initialized with default values.
func DefaultConfig() *SeedConfig {
	return &SeedConfig{
		Environment:     "development",
		TursoURL:        "fincher.db",
		ClickHouseDSN:   "127.0.0.1:9000",
		Seed:            20260828,
		EventsPerVendor: 10000,
		HistoryDays:     120,
		Titles:          7,
		FillerVendors:   0,
		Reset:           false,
		BatchSize:       2000,
	}
}
