/*
audit table creation

Zach Snyder
Oracle XE version 21
*/
-- audit tables
CREATE TABLE users_log_t (
    user_id number(20) not null,
    user_name varchar2(50) not null,
    action varchar2(1 CHAR),
    action_user varchar2(50),
    action_time date
);
CREATE TABLE guilds_log_t (
    guild_id number(20) not null,
    guild_name varchar2(50) not null,
    action varchar2(1 CHAR),
    action_user varchar2(50),
    action_time date
);
CREATE TABLE message_log_t (
    message_id number(20) not null,
    message_time date not null,
    message_guild_id number(20) not null,
    message_channel_id number(20) not null,
    message_user_id number(20) not null,
    message_username varchar2(50) not null,
    message_content varchar2(4000) null,
    action varchar2(1 CHAR),
    action_user varchar2(50),
    action_time date
);
-- views supporting audit tables
CREATE OR REPLACE VIEW users_log AS SELECT * FROM users_log_t;
CREATE OR REPLACE VIEW guilds_log AS SELECT * FROM guilds_log_t;
CREATE OR REPLACE VIEW message_log AS SELECT * FROM message_log_t;
--