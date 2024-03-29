CREATE OR REPLACE FUNCTION public.crpt_passwd()
    RETURNS trigger
    LANGUAGE 'plpgsql'
AS $$
BEGIN NEW.passwd:= crypt(new.passwd, gen_salt('md5'));
RETURN NEW;
END;
$$;

CREATE TRIGGER tgr_crpt_passwd
BEFORE INSERT OR UPDATE 
ON public.user_info
FOR EACH ROW
EXECUTE FUNCTION public.crpt_passwd();

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