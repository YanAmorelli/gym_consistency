CREATE TABLE request_types(
      type_id           INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
      status_desc       VARCHAR(20)
);

INSERT INTO request_types (status_desc) VALUES
    ('P'),
    ('A'),
    ('D');

ALTER TABLE friend_request
    DROP COLUMN request_status;

ALTER TABLE friend_request
    ADD COLUMN request_status INT REFERENCES(type_id);

ALTER TABLE friend_request 
    ADD CONSTRAINT fk_status_type_id 
        FOREIGN KEY(request_status) REFERENCES request_types(type_id);



