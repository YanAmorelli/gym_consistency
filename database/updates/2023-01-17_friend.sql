CREATE TABLE request_types(
      type_id           INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
      status_desc       VARCHAR(20)
);

INSERT INTO request_types (status_desc) VALUES
    ('Accepted'),
    ('Denied'),
    ('Pending');