ALTER TABLE IF EXISTS  user_info
ADD COLUMN IF NOT EXISTS profile_pic bytea[],
ADD COLUMN IF NOT EXISTS changed_passwd boolean default false;