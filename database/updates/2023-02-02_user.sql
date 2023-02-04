CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--Primeiramente dropamos ambos FK e PK da coluna user_id de suas respectivas tabelas.
ALTER TABLE user_attendance 
    DROP CONSTRAINT user_attendance_user_id_fkey;

ALTER TABLE user_info 
    DROP CONSTRAINT user_info_pkey;

--Excluimos aqui a geração de id's pelo identity.
ALTER TABLE user_info 
    ALTER COLUMN user_id DROP IDENTITY;

--Logo depois iremos mudar o tipo da coluna para uuid e dizer para o mesmo ser gerado automaticamente.
ALTER TABLE user_info 
    ALTER COLUMN user_id TYPE UUID USING (uuid_generate_v4()) ;

--Aqui iremos apagar tudo na tabela user presence.
DELETE FROM user_attendance;

-- E trocar a sua coluna user_id para uuid também, detalhe que precisaremos passar para char para poder fazer o cast.
ALTER TABLE user_attendance
    ALTER COLUMN user_id TYPE CHAR;

ALTER TABLE user_attendance
    ALTER COLUMN user_id TYPE UUID USING user_id::uuid;

--Feito isso poderemos adicionar finalmente as respectivas PK e FK para seu devido lugar.
ALTER TABLE user_info ADD PRIMARY KEY (user_id);

ALTER TABLE user_attendance 
ADD CONSTRAINT user_attendance_user_id_fkey FOREIGN KEY (user_id)
REFERENCES user_info(user_id);

-- Nestes demais "Alters" serão para passar as colunas que possuem id de usuário para char e logo depois para uuid.
ALTER TABLE friend_request
    ALTER COLUMN user_sent TYPE CHAR;

ALTER TABLE friend_request
    ALTER COLUMN user_sent TYPE uuid USING user_sent::uuid;

ALTER TABLE friend_request
    ALTER COLUMN user_received TYPE CHAR;

ALTER TABLE friend_request
    ALTER COLUMN user_received TYPE uuid USING user_received::uuid;

ALTER TABLE user_friendship
    ALTER COLUMN user TYPE CHAR;

ALTER TABLE user_friendship
    ALTER COLUMN user TYPE uuid USING user::uuid;

ALTER TABLE user_friendship
    ALTER COLUMN friend TYPE CHAR;

ALTER TABLE user_friendship
    ALTER COLUMN friend TYPE uuid USING friend::uuid;

