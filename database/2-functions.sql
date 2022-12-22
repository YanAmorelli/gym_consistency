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

CREATE OR REPLACE FUNCTION login_username(
	prm_username character varying,
	prm_passwd text)
RETURNS TABLE(user_id bigint, 
              fullname character varying, 
              email character varying, 
              username character varying, 
              passwd text) 
LANGUAGE 'sql'
AS $$
SELECT ui.user_id,ui.fullname,ui.email,ui.username,ui.passwd
FROM user_info ui
WHERE ui.username = prm_username AND 
    (ui.passwd = crypt(prm_passwd, ui.passwd ));
$$;