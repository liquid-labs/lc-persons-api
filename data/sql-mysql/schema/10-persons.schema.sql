CREATE TABLE `persons` (
  `id` int(10),
  `display_name` varchar(255),
-- see ../docs/Relational-Schemas.md#reformatting-data-via-a-trigger
  `phone` varchar(12),
  `email` varchar(255) NOT NULL,
  `phone_backup` varchar(12),
  CONSTRAINT `persons_key` PRIMARY KEY ( `id` ),
  CONSTRAINT `persons_ref_users` FOREIGN KEY ( `id` ) REFERENCES `users` ( `id` )
);
DELIMITER //
CREATE TRIGGER `persons_phone_format`
  BEFORE INSERT ON persons FOR EACH ROW
    BEGIN
      SET new.phone=(SELECT NUMERIC_ONLY(new.phone));
      SET new.phone_backup=(SELECT NUMERIC_ONLY(new.phone_backup));
    END;//
DELIMITER ;
