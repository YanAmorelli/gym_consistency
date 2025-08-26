CREATE TABLE user_info (
       user_id          UUID DEFAULT uuid_generate_v4()PRIMARY KEY,
       fullname         VARCHAR(100),
       username         VARCHAR(30) UNIQUE NOT NULL,
       passwd           TEXT,
       email            VARCHAR(100) UNIQUE,
       profile_pic      bytea[],
       changed_passwd   boolean default false
);

CREATE TABLE user_attendance (
	 attendance_id    BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	 dt_attendance    DATE ,
	 went_gym         BOOL,
	 user_id          UUID references user_info(user_id)
);

CREATE TABLE request_types(
      type_id           INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
      status_acro       CHAR(1),
      status_desc       VARCHAR(15)
);

CREATE TABLE friend_request (
      request_id        INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
      user_sent         UUID,
      user_received     UUID,
      request_status    INT REFERENCES request_types(type_id),
      dt_sented         timestamptz DEFAULT NOW(),
      dt_replied        timestamptz
);


CREATE TABLE user_friendship (
      "user"             UUID,
      friend            UUID,
      PRIMARY KEY ("user",friend)
);

CREATE UNIQUE INDEX dt_attendance_user_id ON user_attendance 
(dt_attendance, user_id);
