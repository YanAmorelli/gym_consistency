CREATE TABLE user_info (
    user_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
    fullname character varying(100),
    username character varying(30),
    passwd text,
    email character varying(100)
);