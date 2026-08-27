create materialized view if not exists fincher.qc
engine = MergeTree()
partition by toYYYYMM(inspected_at)
order by (vendor_id, component, inspected_at)
as select
    id as event_id,
    subject as title_slug,
    JSONExtractString(data, 'package_id') as package_id,
    JSONExtractString(data, 'vendor_id') as vendor_id,
    transform(
        JSONExtractString(data, 'component'),
        ['AUDIO', 'VIDEO', 'SUBTITLE'],
        ['AUDIO', 'VIDEO', 'SUBTITLE']::Array(Enum8('UNKNOWN' = 0, 'AUDIO' = 1, 'VIDEO' = 2, 'SUBTITLE' = 3)),
        'UNKNOWN'
    ) as component,
    JSONExtractString(data, 'language') as language,
    transform(
        JSONExtractString(data, 'status'),
        ['PASSED', 'FAILED', 'WARNING'],
        ['PASSED', 'FAILED', 'WARNING']::Array(Enum8('UNKNOWN' = 0, 'PASSED' = 1, 'FAILED' = 2, 'WARNING' = 3)),
        'UNKNOWN'
    ) as status,
    if(JSONHas(data, 'sync_drift_ms'), JSONExtractFloat(data, 'sync_drift_ms'), NaN) as sync_drift_ms,
    if(JSONHas(data, 'video_corruption_score'), JSONExtractFloat(data, 'video_corruption_score'), NaN) as video_corruption_score,
    transform(
        JSONExtractString(data, 'defect_category'),
        ['NONE', 'AUDIO_SYNC_DRIFT', 'CORRUPT_FRAME', 'SUBTITLE_OVERLAP', 'SLA_BREACH', 'OTHER'],
        ['NONE', 'AUDIO_SYNC_DRIFT', 'CORRUPT_FRAME', 'SUBTITLE_OVERLAP', 'SLA_BREACH', 'OTHER']::Array(Enum8('NONE' = 0, 'AUDIO_SYNC_DRIFT' = 1, 'CORRUPT_FRAME' = 2, 'SUBTITLE_OVERLAP' = 3, 'SLA_BREACH' = 4, 'OTHER' = 5, 'UNKNOWN' = 99)),
        'UNKNOWN'
    ) as defect_category,
    source as inspector_agent,
    time as inspected_at
from fincher.events
where type = 'fincher.qc.completed';
