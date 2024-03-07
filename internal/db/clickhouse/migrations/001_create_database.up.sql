CREATE TABLE IF NOT EXISTS goods(
    Id Int32,
    ProjectId Int32,
    Name String,
    Description String,
    Priority Int32,
    Removed Bool,
    EventTime DateTime DEFAULT now()
) ENGINE = NATS SETTINGS
    nats_url = 'nats://localhost:4222',
    nats_subjects = 'logs',
    nats_format = 'JSONEachRow',
    nats_max_block_size = 100,
    nats_flush_interval_ms = 5000;