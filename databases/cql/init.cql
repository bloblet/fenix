CREATE KEYSPACE IF NOT EXISTS Users
	WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };

DELETE TYPE IF EXISTS  User;
CREATE TYPE User (
	UserID uuid,
	Token varchar,
	Email varchar,
	Salt blob,
	Password blob,
	Username varchar,
	Discriminator varchar,
	Servers list<uuid>,
	Friends list<uuid>,
	Activity Activity,
	Settings UserSettings
);

DELETE TYPE IF EXISTS Activity;
CREATE TYPE Activity (
	StopAt timestamp,
	Name varchar,
	Emoji varchar
);

DELETE TYPE IF EXISTS UserSettings;
CREATE TYPE UserSettings (
);

CREATE KEYSPACE IF NOT EXISTS Channels
	WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };

DELETE TYPE IF EXISTS Channel;
CREATE TYPE Channel (
	ChannelID uuid,
	Topic varchar,
	ChannelPermissions Permissions
);

DELETE TYPE IF EXISTS Permissions;
CREATE TYPE Permissions (
	RolePermissions list<RolePermission>,
	UserPermissions list<UserPermission>
);

DELETE TYPE IF EXISTS RolePermission;
CREATE TYPE RolePermission (
	RoleID uuid,
	Read boolean,
	Talk boolean,
	React boolean,
	DeleteMessages boolean,
	DeleteReactions boolean,
	ChangeTopic boolean,
	Mute boolean
);

DELETE TYPE IF EXISTS UserPermission;
CREATE TYPE UserPermission (
	UserID uuid,
	Read boolean,
	Talk boolean,
	React boolean,
	DeleteMessages boolean,
	DeleteReactions boolean,
	ChangeTopic boolean,
	Mute boolean
);

CREATE KEYSPACE IF NOT EXISTS Roles
WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };

DELETE TYPE IF EXISTS Role;
CREATE TYPE Role (
	RoleID uuid,
	Color tuple<tinyint, tinyint, tinyint>,
	Name varchar,
	Members list<uuid>
	Pingable boolean,
	hierarchy tinyint,
	display boolean
);

CREATE KEYSPACE IF NOT EXISTS Servers
WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };
DELETE TYPE IF EXISTS Server;
CREATE TYPE Server (
	ServerID uuid,
	Name varchar,
	Members list<varchar>
);

CREATE KEYSPACE IF NOT EXISTS Messages
WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };
DELETE TYPE IF EXISTS Message;
CREATE TYPE Message (
	ServerID uuid,
	ChannelID uuid,
	UserID uuid,
	MessageID uuid,
	Sent Timestamp,
	Contents varchar
);

CREATE KEYSPACE IF NOT EXISTS Comments
WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };
DELETE TYPE IF EXISTS Comment;
CREATE TYPE Comment (
	CommentID uuid,
	MessageID uuid,
	AuthorID uuid,
	Contents varchar	
);

CREATE KEYSPACE IF NOT EXISTS Reactions
WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };
DELETE TYPE IF EXISTS Reactions;
CREATE TYPE Reactions (
	ReactionID uuid,
	MessageID uuid,
	Emoji varchar,
	Reactees list<uuid>
);
