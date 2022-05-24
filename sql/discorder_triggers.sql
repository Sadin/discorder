DROP TRIGGER bot.messages_aiud;
CREATE OR REPLACE TRIGGER bot.messages_aiud AFTER INSERT OR UPDATE OR DELETE ON bot.message_t FOR EACH ROW
    DECLARE
        v_date date:= SYSDATE;
        v_action varchar2(1 CHAR):= 'I';
        v_user varchar2(50):= NULL;
    BEGIN
        IF DELETING THEN
            INSERT INTO bot.message_log VALUES (
                :old.message_id,
                :old.message_time,
                :old.message_guild_id,
                :old.message_channel_id,
                :old.message_user_id,
                :old.message_username,
                :old.message_content,
                'D',
                v_user,
                v_date
            );
        ELSE
            IF INSERTING THEN
                INSERT INTO bot.message_log VALUES (
                    :new.message_id,
                    :new.message_time,
                    :new.message_guild_id,
                    :new.message_channel_id,
                    :new.message_user_id,
                    :new.message_username,
                    :new.message_content,
                    v_action,
                    v_user,
                    v_date
                );
            ELSE
                v_action := 'U';
                INSERT INTO bot.message_log VALUES (
                    :old.message_id,
                    :old.message_time,
                    :old.message_guild_id,
                    :old.message_channel_id,
                    :old.message_user_id,
                    :old.message_username,
                    :old.message_content,
                    v_action,
                    v_user,
                    v_date
                );
        END IF;
    END IF;
END;