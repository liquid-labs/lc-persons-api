INSERT INTO entities (pub_id) VALUES ('4BE66BE5-2A62-11E9-B987-42010A8003FF');
-- TODO: 'SET' is not ANSI SQL; for this and other reasons, we want to do a
-- replacement scheme. Possibly something like:
-- 1) Name template files with a commen prefix ('.sql.template').
-- 2) Use bash subsitutios, so "VALUES ($JANE_DOE_ID)"
-- 3) Have a 'template.vars' file.
-- 4) source template.vars; for $TEMPLATE in ...; do ...; eval "$(cat "$TEMPLATE")" > $SQL_FILE; done
SET @jane_doe_id=LAST_INSERT_ID();
INSERT INTO users (id, auth_id, active) VALUES (@jane_doe_id,'abcdefg123',0);
INSERT INTO persons (id, display_name, phone, email, phone_backup) VALUES (@jane_doe_id,'Jane Doe','5555551111','janedoe@test.com',NULL);
