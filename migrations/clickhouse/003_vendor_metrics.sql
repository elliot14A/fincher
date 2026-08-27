create materialized view if not exists fincher.vendor_metrics
engine = SummingMergeTree()
partition by toYYYYMM(recorded_date)
order by (vendor_id, component, recorded_date)
as select
    toDate(time) as recorded_date,
    JSONExtractString(data, 'vendor_id') as vendor_id,
    transform(
        JSONExtractString(data, 'component'),
        ['AUDIO', 'VIDEO', 'SUBTITLE'],
        ['AUDIO', 'VIDEO', 'SUBTITLE']::Array(Enum8('UNKNOWN' = 0, 'AUDIO' = 1, 'VIDEO' = 2, 'SUBTITLE' = 3)),
        'UNKNOWN'
    ) as component,
    count() as total_inspections,
    countIf(transform(JSONExtractString(data, 'status'), ['PASSED', 'FAILED', 'WARNING'], ['PASSED', 'FAILED', 'WARNING']::Array(Enum8('UNKNOWN' = 0, 'PASSED' = 1, 'FAILED' = 2, 'WARNING' = 3)), 'UNKNOWN') = 'FAILED') as failed_inspections,
    countIf(transform(JSONExtractString(data, 'status'), ['PASSED', 'FAILED', 'WARNING'], ['PASSED', 'FAILED', 'WARNING']::Array(Enum8('UNKNOWN' = 0, 'PASSED' = 1, 'FAILED' = 2, 'WARNING' = 3)), 'UNKNOWN') = 'WARNING') as warning_inspections,
    countIf(transform(JSONExtractString(data, 'status'), ['PASSED', 'FAILED', 'WARNING'], ['PASSED', 'FAILED', 'WARNING']::Array(Enum8('UNKNOWN' = 0, 'PASSED' = 1, 'FAILED' = 2, 'WARNING' = 3)), 'UNKNOWN') in ('PASSED', 'FAILED', 'WARNING')) as measured_status_count,
    sumIf(JSONExtractFloat(data, 'sync_drift_ms'), JSONHas(data, 'sync_drift_ms')) as total_sync_drift_ms,
    countIf(JSONHas(data, 'sync_drift_ms')) as measured_drift_count
from fincher.events
where type = 'fincher.qc.completed'
group by recorded_date, vendor_id, component;
