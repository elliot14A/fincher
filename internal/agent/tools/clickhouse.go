package tools

import (
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/mcp"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

func NewClickhouseQueryTool(client *mcp.Client) (tool.Tool, error) {
	if client == nil {
		return nil, domainerrors.NewWithOp("tools.NewClickhouseQueryTool", domainerrors.CodeInvalidInput, "mcp client cannot be nil", nil)
	}
	return functiontool.New(
		functiontool.Config{
			Name:        "query_clickhouse",
			Description: clickhouseQueryToolDescription,
		},
		func(ctx agent.Context, args map[string]any) ([]map[string]any, error) {
			sql := ensureLimit(ExtractSQLArg(args), clickhouseRowLimit)
			rows, err := mcp.RunQueryRows(ctx, client, sql)
			if err != nil {
				return nil, err
			}
			capped, _ := capRows(rows)
			return capped, nil
		},
	)
}

const clickhouseQueryToolDescription = `Run a read-only SQL SELECT against the Fincher ClickHouse analytical history via the ClickHouse MCP server and get matching rows. The MCP connection is pinned read-only; write/DDL statements will fail. This is the historical record of everything that happened over time.

Schema (database: fincher):
- events(id UUID, type LowCardinality(String), source String, subject LowCardinality(String), time DateTime64(3,'UTC'), data String, severity Enum8('INFO','WARN','CRITICAL')). This is the raw CloudEvents stream. subject = the title slug (or 'GLOBAL'). data is a JSON string; extract fields with JSONExtractString(data,'field') / JSONExtractFloat(data,'field') / JSONHas(data,'field'). Common event types: fincher.qc.completed, fincher.delivery.held, fincher.delivery.released, fincher.master.revised, fincher.package.invalidated.
- qc(event_id UUID, title_slug String, package_id String, vendor_id String, component Enum('AUDIO','VIDEO','SUBTITLE'), language String, status Enum('PASSED','FAILED','WARNING'), sync_drift_ms Float, video_corruption_score Float, defect_category Enum('NONE','AUDIO_SYNC_DRIFT','CORRUPT_FRAME','SUBTITLE_OVERLAP','SLA_BREACH','OTHER'), inspector_agent String, inspected_at DateTime64). Materialized QC inspection results.
- vendor_metrics(recorded_date Date, vendor_id String, component Enum('AUDIO','VIDEO','SUBTITLE'), total_inspections UInt64, failed_inspections UInt64, warning_inspections UInt64, measured_status_count UInt64, total_sync_drift_ms Float, measured_drift_count UInt64). SummingMergeTree rollups; aggregate with sum()/GROUP BY vendor_id, component to get lifetime totals.

Examples:
- Recent events for a title, newest first:
  select toString(time) as t, type, severity, data from fincher.events where subject = '<slug>' order by time desc limit 20
- QC failures for a vendor:
  select title_slug, component, language, defect_category, toString(inspected_at) as t from fincher.qc where vendor_id = '<vendor_id>' and status = 'FAILED' order by inspected_at desc limit 20
- Vendor accuracy rollup (pass rate) for AUDIO:
  select vendor_id, sum(total_inspections) as total, sum(failed_inspections) as failed, 1 - (sum(failed_inspections) / nullif(sum(total_inspections),0)) as pass_rate from fincher.vendor_metrics where component = 'AUDIO' group by vendor_id order by pass_rate desc

Results are capped at 100 rows (and a byte budget); a LIMIT is auto-applied if you omit one. The events table is huge, so ALWAYS filter (by subject/type/time) or aggregate — never select it unbounded. Prefer several small targeted SELECTs over one giant query. Never fabricate rows; use only rows these queries return.`
