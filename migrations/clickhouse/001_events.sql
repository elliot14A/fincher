create table if not exists fincher.events (
    id UUID default generateUUIDv4(),
    type LowCardinality(String),
    source LowCardinality(String),
    subject LowCardinality(String) default 'GLOBAL',
    time DateTime64(3, 'UTC') codec(Delta(8), ZSTD(1)),
    data String,
    severity Enum8('INFO' = 1, 'WARN' = 2, 'CRITICAL' = 3),
    datacontenttype LowCardinality(String) default 'application/json',
    created_at DateTime default now()
) engine = MergeTree()
partition by toYYYYMM(time)
order by (type, subject, time, id)
ttl toDateTime(time) + interval 2 year delete;
