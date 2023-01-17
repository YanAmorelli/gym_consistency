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

CREATE OR REPLACE FUNCTION auth_login(p_user text,p_passwd text) 
RETURNS boolean
LANGUAGE 'plpgsql'
AS $$
DECLARE 
    auth	boolean;
BEGIN
SELECT 
	(passwd = crypt(p_passwd,passwd)) AS authed_login
INTO 
	auth
FROM 
	user_info
WHERE
	username = p_user;
IF auth IS NOT true THEN  
	RETURN false;
		ELSE RETURN true;
END IF;
END;
$$;