/*
audit table creation

Zach Snyder
Oracle XE version 21
*/
-- audit tables
CREATE TABLE bot.users_log_t (
    user_id number(20) not null,
    user_name varchar2(50) not null,
    action varchar2(1 CHAR),
    action_user varchar2(50),
    action_time date
);
CREATE TABLE bot.guilds_log_t (
    guild_id number(20) not null,
    guild_name varchar2(50) not null,
    action varchar2(1 CHAR),
    action_user varchar2(50),
    action_time date
);
CREATE TABLE bot.message_log_t (
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
CREATE OR REPLACE VIEW bot.users_log AS SELECT * FROM bot.users_log_t;
CREATE OR REPLACE VIEW bot.guilds_log AS SELECT * FROM bot.guilds_log_t;
CREATE OR REPLACE VIEW bot.message_log AS SELECT * FROM bot.message_log_t;
--