package migrations

import "embed"

// FS embeds all ClickHouse migration SQL scripts.
//
//go:embed clickhouse/*.sql
var FS embed.FS
