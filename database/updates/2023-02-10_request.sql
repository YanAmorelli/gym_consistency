ALTER TABLE request_types 
    RENAME COLUMN status_desc TO status_acro;

ALTER TABLE request_types
    ADD COLUMN status_desc VARCHAR(15);

ALTER TABLE user_info
    ALTER COLUMN username SET NOT NULL;

