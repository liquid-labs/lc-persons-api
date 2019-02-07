INSERT INTO entities (pub_id) VALUES ('4BE66BE5-2A62-11E9-B987-42010A8003FF');
SET @jane_doe_id=LAST_INSERT_ID();
INSERT INTO users (id, active) VALUES (@jane_doe_id,0);
INSERT INTO persons (id, name, phone, email, phone_backup) VALUES (@jane_doe_id,'Jane Doe','5555551111','janedoe@test.com',NULL);
