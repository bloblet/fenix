CREATE KEYSPACE IF NOT EXISTS Users WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };
CREATE KEYSPACE IF NOT EXISTS Channels WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };
CREATE KEYSPACE IF NOT EXISTS Roles WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };
CREATE KEYSPACE IF NOT EXISTS Servers WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };
CREATE KEYSPACE IF NOT EXISTS Messages WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };
CREATE KEYSPACE IF NOT EXISTS Comments WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };
CREATE KEYSPACE IF NOT EXISTS Reactions WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };

DROP TABLE IF EXISTS Users.User;
CREATE TABLE Users.User (
	UserID TEXT PRIMARY KEY,
	Username TEXT,
	Discriminator TEXT,

	Activity TEXT,
	Servers list<TEXT>,
	Friends list<TEXT>,

	UserToken TEXT,
	Email TEXT,
	Salt blob,
	Password blob
);

DROP TABLE IF EXISTS Messages.Messages;
CREATE TABLE Messages.Messages (
    MessageID TEXT,
    UserID TEXT,
    Content TEXT,
    CreatedAt TIMESTAMP,
    PRIMARY KEY (MessageID, CreatedAt)
) WITH CLUSTERING ORDER BY (CreatedAt DESC);
