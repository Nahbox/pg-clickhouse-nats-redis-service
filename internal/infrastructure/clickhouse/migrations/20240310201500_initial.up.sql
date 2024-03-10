CREATE TABLE IF NOT EXISTS logs(
                                   id Int64,
                                   project_id Int64,
                                   name String,
                                   description String,
                                   priority Int64,
                                   removed Bool,
                                   event_time DateTime DEFAULT now()
) ENGINE = MergeTree() ORDER BY id;