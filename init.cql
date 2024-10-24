CREATE KEYSPACE IF NOT EXISTS dataservices WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

USE dataservices;

CREATE TABLE IF NOT EXISTS messages (
  channel_id bigint,
  message_id bigint,
  author_id bigint,
  content text,
  PRIMARY KEY (channel_id, message_id)
) WITH CLUSTERING ORDER BY (message_id DESC);

INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (1, 1001, 5001, 'Hello, this is the first message.') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (1, 1002, 5002, 'Here is another message.') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (1, 1003, 5003, 'Bulk insert message example.') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (1, 1004, 5004, 'Just adding some dummy data.') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (2, 2001, 6001, 'Second channel first message.') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (2, 2002, 6002, 'More messages in another channel.') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (2, 2003, 6003, 'Sample content goes here.') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (2, 2004, 6004, 'Just more dummy data.') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (3, 3001, 7001, 'Channel 3, starting new batch.') IF NOT EXISTS;
INSERT INTO messages (channel_id, message_id, author_id, content) VALUES (3, 3002, 7002, 'Another message in channel 3.') IF NOT EXISTS;
