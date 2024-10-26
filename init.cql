CREATE KEYSPACE IF NOT EXISTS dataservices WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

USE dataservices;

CREATE TABLE IF NOT EXISTS messages (
  channel_id bigint,
  message_id bigint,
  author_id bigint,
  content text,
  PRIMARY KEY (channel_id, message_id)
) WITH CLUSTERING ORDER BY (message_id DESC);

INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (1, 1000, 5000, 'Hello, this is the first message.') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (2, 2000, 6000, 'Second channel first message.') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (3, 3000, 7000, 'Channel 3000') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (4, 4000, 8000, 'Channel 4000') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (5, 5000, 9000, 'Channel 5000') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (6, 6000, 10000, 'Channel 6000') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (7, 7000, 11000, 'Channel 7000') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (8, 8000, 12000, 'Channel 8000') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (9, 9000, 13000, 'Channel 9000') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (10, 10000, 14000, 'Channel 10000') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (11, 11000, 15000, 'Channel 11000') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (12, 12000, 16000, 'Channel 12000') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (13, 13000, 17000, 'Channel 13000') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (14, 14000, 18000, 'Channel 14000') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (15, 15000, 19000, 'Channel 15000') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (16, 16000, 20000, 'Channel 16000') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (17, 17000, 21000, 'Channel 17000') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (18, 18000, 22000, 'Channel 18000') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (19, 19000, 23000, 'Channel 19000') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (20, 20000, 24000, 'Channel 20000') IF NOT EXISTS;