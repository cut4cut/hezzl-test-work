CREATE TABLE IF NOT EXISTS new_user_queue (
  id          int,
  name        String
)
ENGINE = Kafka SETTINGS kafka_broker_list = 'kafka:29092',
                            kafka_topic_list = 'user',
                            kafka_group_name = 'group1',
                            kafka_format = 'JSONEachRow';

CREATE TABLE new_user (
  id          int,
  name        String
)
ENGINE = ReplacingMergeTree()
ORDER BY (id, name);

CREATE MATERIALIZED VIEW new_user_consumer
TO new_user
AS SELECT *
FROM new_user_queue;