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
--
CREATE TABLE bot.guilds
    (
        guild_id number(20) not null,
        guild_name varchar2(50) not null,
        CONSTRAINT guild_pk PRIMARY KEY (guild_id)
    );
CREATE TABLE bot.users
    (
        user_id number(20) not null,
        user_name varchar2(50) not null,
        CONSTRAINT user_pk PRIMARY KEY (user_id)
    );

