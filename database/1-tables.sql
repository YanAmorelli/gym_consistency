CREATE TABLE user_info (
    user_id   bigint NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    fullname     character varying(100),
    username     character varying(30),
    passwd       text,
    email        character varying(100)
);

CREATE TABLE user_attendance (
	 presence_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	 dt_presence DATE ,
	 went_gym    BOOL,
	 user_id INT references user_info(user_id)
);

CREATE UNIQUE INDEX dt_presence_user_id ON user_presence 
(dt_presence, user_id);