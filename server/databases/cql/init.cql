CREATE KEYSPACE IF NOT EXISTS Users WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };
CREATE KEYSPACE IF NOT EXISTS Channels WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };
CREATE KEYSPACE IF NOT EXISTS Roles WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };
CREATE KEYSPACE IF NOT EXISTS Servers WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };
CREATE KEYSPACE IF NOT EXISTS Messages WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };
CREATE KEYSPACE IF NOT EXISTS Comments WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };
CREATE KEYSPACE IF NOT EXISTS Reactions WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };

prepare 