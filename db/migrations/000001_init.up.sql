-- Project Name : チーム03
-- Date/Time    : 2023/06/20 18:57:43
-- Author       : kaned
-- RDBMS Type   : PostgreSQL
-- Application  : A5:SQL Mk-2

/*
  << 注意！！ >>
  BackupToTempTable, RestoreFromTempTable疑似命令が付加されています。
  これにより、drop table, create table 後もデータが残ります。
  この機能は一時的に $$TableName のような一時テーブルを作成します。
  この機能は A5:SQL Mk-2でのみ有効であることに注意してください。
*/
CREATE SCHEMA IF NOT EXISTS "public";

-- atlas:delimiter -- end
CREATE OR REPLACE FUNCTION nanoid(
    size int DEFAULT 21,
    alphabet text DEFAULT '_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ'
)
    RETURNS text
    LANGUAGE plpgsql
    volatile
AS
$$
DECLARE
    idBuilder      text := '';
    counter        int  := 0;
    bytes          bytea;
    alphabetIndex  int;
    alphabetArray  text[];
    alphabetLength int;
    mask           int;
    step           int;
BEGIN
    alphabetArray := regexp_split_to_array(alphabet, '');
    alphabetLength := array_length(alphabetArray, 1);
    mask := (2 << cast(floor(log(alphabetLength - 1) / log(2)) as int)) - 1;
    step := cast(ceil(1.6 * mask * size / alphabetLength) AS int);

    while true
        loop
            bytes := gen_random_bytes(step);
            while counter < step
                loop
                    alphabetIndex := (get_byte(bytes, counter) & mask) + 1;
                    if alphabetIndex <= alphabetLength then
                        idBuilder := idBuilder || alphabetArray[alphabetIndex];
                        if length(idBuilder) = size then
                            return idBuilder;
                        end if;
                    end if;
                    counter := counter + 1;
                end loop;

            counter := 0;
        end loop;
END
$$;
-- end

CREATE TYPE type_shopping_item AS (
    item VARCHAR,
    is_on BOOLEAN
);

-- ファボ／リム履歴
-- * BackupToTempTable
DROP TABLE if exists fav_history CASCADE;

-- * RestoreFromTempTable
CREATE TABLE fav_history (
  recipe_id CHAR(21) NOT NULL
  , usr_id CHAR(21) NOT NULL
  , is_fav BOOLEAN NOT NULL
  , created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
) ;

-- フォロー／リム履歴
-- * BackupToTempTable
DROP TABLE if exists follow_history CASCADE;

-- * RestoreFromTempTable
CREATE TABLE follow_history (
  chef_id CHAR(21) NOT NULL
  , usr_id CHAR(21) NOT NULL
  , is_follow BOOLEAN NOT NULL
  , created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
) ;

-- 買い物リスト
-- * BackupToTempTable
DROP TABLE if exists shopping_list CASCADE;

-- * RestoreFromTempTable
CREATE TABLE shopping_list (
  id CHAR(21) DEFAULT nanoid() NOT NULL
  , usr_id CHAR(21) NOT NULL
  , recipe_id CHAR(21)
  , is_fair_copy BOOLEAN NOT NULL
  , shopping_item type_shopping_item[] DEFAULT ARRAY[]::type_shopping_item[] NOT NULL
  , created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
  , updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
  , CONSTRAINT shopping_list_PKC PRIMARY KEY (id)
) ;

-- ファボ中
-- * BackupToTempTable
DROP TABLE if exists favoring CASCADE;

-- * RestoreFromTempTable
CREATE TABLE favoring (
  id CHAR(21) DEFAULT nanoid() NOT NULL
  , recipe_id CHAR(21) NOT NULL
  , usr_id CHAR(21) NOT NULL
  , created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , CONSTRAINT favoring_PKC PRIMARY KEY (id)
) ;

-- フォロー中
-- * BackupToTempTable
DROP TABLE if exists following CASCADE;

-- * RestoreFromTempTable
CREATE TABLE following (
  id CHAR(21) DEFAULT nanoid() NOT NULL
  , chef_id CHAR(21) NOT NULL
  , usr_id CHAR(21) NOT NULL
  , created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
  , CONSTRAINT following_PKC PRIMARY KEY (id)
) ;

-- sns
-- * BackupToTempTable
DROP TABLE if exists sns CASCADE;

-- * RestoreFromTempTable
CREATE TABLE sns (
  id CHAR(21) DEFAULT nanoid() NOT NULL
  , chef_id CHAR(21) NOT NULL
  , name VARCHAR NOT NULL
  , account_name VARCHAR NOT NULL
  , num_followers integer DEFAULT 0 NOT NULL
  , link VARCHAR
  , CONSTRAINT sns_PKC PRIMARY KEY (id)
) ;

-- ユーザー
-- * BackupToTempTable
DROP TABLE if exists usr CASCADE;

-- * RestoreFromTempTable
CREATE TABLE usr (
  id CHAR(21) DEFAULT nanoid() NOT NULL
  , email VARCHAR NOT NULL
  , name VARCHAR NOT NULL
  , image_url VARCHAR
  , created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , CONSTRAINT usr_PKC PRIMARY KEY (id)
) ;

-- シェフのレシピ＆マイレシピ
-- * BackupToTempTable
DROP TABLE if exists recipe CASCADE;

-- * RestoreFromTempTable
CREATE TABLE recipe (
  id CHAR(21) DEFAULT nanoid() NOT NULL
  , chef_id CHAR(21)
  , usr_id CHAR(21)
  , title VARCHAR NOT NULL
  , comment VARCHAR NOT NULL
  , servings integer NOT NULL
  , ingredient VARCHAR[] DEFAULT ARRAY[]::VARCHAR[] NOT NULL
  , method VARCHAR[] DEFAULT ARRAY[]::VARCHAR[] NOT NULL
  , image_url VARCHAR
  , link VARCHAR[] DEFAULT ARRAY[]::VARCHAR[] NOT NULL
  , access_level integer NOT NULL
  , created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , CONSTRAINT recipe_PKC PRIMARY KEY (id)
) ;

-- シェフ
-- * BackupToTempTable
DROP TABLE if exists chef CASCADE;

-- * RestoreFromTempTable
CREATE TABLE chef (
  id CHAR(21) DEFAULT nanoid() NOT NULL
  , email VARCHAR
  , name VARCHAR NOT NULL
  , image_url VARCHAR
  , comment VARCHAR
  , created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , CONSTRAINT chef_PKC PRIMARY KEY (id)
) ;

ALTER TABLE shopping_list
  ADD CONSTRAINT shopping_list_FK1 FOREIGN KEY (recipe_id) REFERENCES recipe(id);

ALTER TABLE fav_history
  ADD CONSTRAINT fav_history_FK1 FOREIGN KEY (usr_id) REFERENCES usr(id);

ALTER TABLE fav_history
  ADD CONSTRAINT fav_history_FK2 FOREIGN KEY (recipe_id) REFERENCES recipe(id);

ALTER TABLE follow_history
  ADD CONSTRAINT follow_history_FK1 FOREIGN KEY (usr_id) REFERENCES usr(id);

ALTER TABLE follow_history
  ADD CONSTRAINT follow_history_FK2 FOREIGN KEY (chef_id) REFERENCES chef(id);

ALTER TABLE recipe
  ADD CONSTRAINT recipe_FK1 FOREIGN KEY (usr_id) REFERENCES usr(id);

ALTER TABLE shopping_list
  ADD CONSTRAINT shopping_list_FK2 FOREIGN KEY (usr_id) REFERENCES usr(id);

ALTER TABLE sns
  ADD CONSTRAINT sns_FK1 FOREIGN KEY (chef_id) REFERENCES chef(id);

ALTER TABLE recipe
  ADD CONSTRAINT recipe_FK2 FOREIGN KEY (chef_id) REFERENCES chef(id);

ALTER TABLE favoring
  ADD CONSTRAINT favoring_FK1 FOREIGN KEY (recipe_id) REFERENCES recipe(id);

ALTER TABLE favoring
  ADD CONSTRAINT favoring_FK2 FOREIGN KEY (usr_id) REFERENCES usr(id);

ALTER TABLE following
  ADD CONSTRAINT following_FK1 FOREIGN KEY (usr_id) REFERENCES usr(id);

ALTER TABLE following
  ADD CONSTRAINT following_FK2 FOREIGN KEY (chef_id) REFERENCES chef(id);

COMMENT ON TABLE fav_history IS 'ファボ／リム履歴';
COMMENT ON COLUMN fav_history.recipe_id IS '';
COMMENT ON COLUMN fav_history.usr_id IS '';
COMMENT ON COLUMN fav_history.is_fav IS 'ファボ／リム';
COMMENT ON COLUMN fav_history.created_at IS '';

COMMENT ON TABLE follow_history IS 'フォロー／リム履歴';
COMMENT ON COLUMN follow_history.chef_id IS '';
COMMENT ON COLUMN follow_history.usr_id IS '';
COMMENT ON COLUMN follow_history.is_follow IS 'フォロー／リム';
COMMENT ON COLUMN follow_history.created_at IS '';

COMMENT ON TABLE shopping_list IS '買い物リスト';
COMMENT ON COLUMN shopping_list.id IS '';
COMMENT ON COLUMN shopping_list.usr_id IS '';
COMMENT ON COLUMN shopping_list.recipe_id IS 'NULL=メモ';
COMMENT ON COLUMN shopping_list.is_fair_copy IS '清書or下書き';
COMMENT ON COLUMN shopping_list.shopping_item IS '明細';
COMMENT ON COLUMN shopping_list.created_at IS '';
COMMENT ON COLUMN shopping_list.updated_at IS '';

COMMENT ON TABLE favoring IS 'ファボ中';
COMMENT ON COLUMN favoring.id IS '';
COMMENT ON COLUMN favoring.recipe_id IS '';
COMMENT ON COLUMN favoring.usr_id IS '';
COMMENT ON COLUMN favoring.created_at IS '';

COMMENT ON TABLE following IS 'フォロー中:Table to show favorite chefs of a user';
COMMENT ON COLUMN following.id IS '';
COMMENT ON COLUMN following.chef_id IS '';
COMMENT ON COLUMN following.usr_id IS '';
COMMENT ON COLUMN following.created_at IS '';

COMMENT ON TABLE sns IS 'sns';
COMMENT ON COLUMN sns.id IS '';
COMMENT ON COLUMN sns.chef_id IS '';
COMMENT ON COLUMN sns.name IS '';
COMMENT ON COLUMN sns.account_name IS '';
COMMENT ON COLUMN sns.num_followers IS '';
COMMENT ON COLUMN sns.link IS '';

COMMENT ON TABLE usr IS 'ユーザー';
COMMENT ON COLUMN usr.id IS '';
COMMENT ON COLUMN usr.email IS 'ログインemail';
COMMENT ON COLUMN usr.name IS '登録名';
COMMENT ON COLUMN usr.image_url IS '登録画像';
COMMENT ON COLUMN usr.created_at IS '';
COMMENT ON COLUMN usr.updated_at IS '';

COMMENT ON TABLE recipe IS 'シェフのレシピ＆マイレシピ';
COMMENT ON COLUMN recipe.id IS '';
COMMENT ON COLUMN recipe.chef_id IS '';
COMMENT ON COLUMN recipe.usr_id IS '';
COMMENT ON COLUMN recipe.title IS 'レシピタイトル';
COMMENT ON COLUMN recipe.comment IS 'レシピコメント';
COMMENT ON COLUMN recipe.servings IS '＊人前';
COMMENT ON COLUMN recipe.ingredient IS '材料';
COMMENT ON COLUMN recipe.method IS '作り方';
COMMENT ON COLUMN recipe.image_url IS '画像';
COMMENT ON COLUMN recipe.link IS 'リンク';
COMMENT ON COLUMN recipe.access_level IS '公開？:公開、限定公開、非公開、下書き';
COMMENT ON COLUMN recipe.created_at IS '';
COMMENT ON COLUMN recipe.updated_at IS '';

COMMENT ON TABLE chef IS 'シェフ';
COMMENT ON COLUMN chef.id IS '';
COMMENT ON COLUMN chef.email IS 'ログインemail';
COMMENT ON COLUMN chef.name IS '登録名';
COMMENT ON COLUMN chef.image_url IS '登録画像';
COMMENT ON COLUMN chef.comment IS 'シェフコメント';
COMMENT ON COLUMN chef.created_at IS '';
COMMENT ON COLUMN chef.updated_at IS '';

