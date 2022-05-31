/*
discorder schema setup script

Zach Snyder
Oracle XE version 21
*/
-- Workaround for XE in container.
alter session set "_ORACLE_SCRIPT"=TRUE;
--
CREATE TABLESPACE tbs_bot_01
    DATAFILE 'tbs_bot_01.dat'
    SIZE 10M
    REUSE
    AUTOEXTEND ON NEXT 10M MAXSIZE 1G
    ONLINE
    LOGGING;
--
CREATE TEMPORARY TABLESPACE tbs_bot_temp
    TEMPFILE 'tbs_bot_temp.dbf'
        SIZE 5M
        AUTOEXTEND ON;
--
CREATE USER bot
    IDENTIFIED BY out_standing1
    DEFAULT TABLESPACE tbs_bot_01
    TEMPORARY TABLESPACE tbs_bot_temp
    QUOTA 20M on tbs_bot_01;
--
GRANT create session TO bot
GRANT create table TO bot;
GRANT create view TO bot;
GRANT create any trigger TO bot;
GRANT create any procedure TO bot;
GRANT create sequence TO bot;
GRANT create synonym TO bot;
GRANT create session TO bot;
-- init tables
-- lookup
DROP TABLE bot.event_t (
    event_id number(20) not null,
    event_type_id number(20) not null,
    CONSTRAINT event
);
-- data
DROP TABLE bot.guilds_t;
CREATE TABLE bot.guilds_t
    (
        guild_id number(20) not null,
        guild_name varchar2(50) not null,
        CONSTRAINT guild_pk PRIMARY KEY (guild_id)
);
DROP TABLE bot.users_t;
CREATE TABLE bot.users_t
    (
        user_id number(20) not null,
        user_name varchar2(50) not null,
        CONSTRAINT user_pk PRIMARY KEY (user_id)
);
DROP TABLE bot.message_t;
CREATE TABLE bot.message_t
    (
        message_id number(20) not null,
        message_time date not null,
        message_guild_id number(20) not null,
        message_channel_id number(20) not null,
        message_user_id number(20) not null,
        message_username varchar2(50) not null,
        message_content varchar2(4000) null,
        CONSTRAINT message_pk PRIMARY KEY (message_id)
);
ALTER TABLE bot.message_t
ADD CONSTRAINT fk_message_guild_id
	FOREIGN KEY (message_guild_id)
	REFERENCES bot.guild_t(guild_id);
-- views
CREATE OR REPLACE VIEW bot.guilds AS SELECT * from bot.guilds_t;
CREATE OR REPLACE VIEW bot.users AS SELECT * from bot.users_t;
CREATE OR REPLACE VIEW bot.message AS SELECT * from bot.message_t;
--
GRANT INSERT, UPDATE, DELETE on bot.message TO bot;
GRANT INSERT, UPDATE, DELETE on bot.users TO bot;
GRANT INSERT, UPDATE, DELETE on bot.guilds TO bot;
--
commit;
