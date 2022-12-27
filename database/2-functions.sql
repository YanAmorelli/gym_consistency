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

CREATE OR REPLACE FUNCTION login_user(
	passwd_inform text,
	passwd_saved text)
    RETURNS boolean
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE PARALLEL UNSAFE
AS $BODY$
DECLARE 
	auth	boolean;
BEGIN 
auth:= passwd_saved = crypt(passwd_inform, passwd_saved );
RETURN auth;
END;
$BODY$;