CREATE TABLE user_info (
    user_id             BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    fullname            VARCHAR(100),
    username            VARCHAR(30),
    passwd              TEXT,
    email               VARCHAR(100)
);

CREATE TABLE user_attendance (
	 attendance_id    BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	 dt_attendance    DATE ,
	 went_gym         BOOL,
	 user_id          INT references user_info(user_id)
);

CREATE TABLE friend_request (
      request_id        INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
      user_sent         INT,
      user_received     INT,
      request_status    BOOL,
      dt_sented         timestamptz DEFAULT NOW(),
      dt_replied        timestamptz
);

CREATE TABLE user_friendship (
      user              INT,
      friend            INT,
      PRIMARY KEY (user,friend)
);

CREATE UNIQUE INDEX dt_attendance_user_id ON user_attendance 
(dt_attendance, user_id);