CREATE OR REPLACE FUNCTION auth_login(pr_username text, pr_pass text)
RETURNS TABLE (user_id uuid, fullname text, username text,email text)
LANGUAGE 'sql'
AS $$
SELECT 
	ui.user_id,
	ui.fullname,
	ui.username,
	ui.email
FROM user_info ui
WHERE
ui.passwd = (crypt(pr_pass,ui.passwd)) and ui.username = pr_username;
$$;

ALTER TABLE user_info
    ALTER COLUMN user_id SET DEFAULT uuid_generate_v4();