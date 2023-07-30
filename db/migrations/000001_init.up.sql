-- Add new schema named "public"
CREATE SCHEMA IF NOT EXISTS "public";

--
-- PGroonga
-- https://pgroonga.github.io/ja/tutorial/
--

CREATE EXTENSION IF NOT EXISTS pgroonga;
SET enable_seqscan = off;

--
-- データ型を登録
--

DROP TYPE if exists type_chef_link CASCADE;

CREATE TYPE type_chef_link AS (
    label TEXT,
    url TEXT
);

DROP TYPE if exists type_recipe_method CASCADE;

CREATE TYPE type_recipe_method AS (
    html TEXT,
    supplement JSONB
);

-- Project Name : チーム03
-- Date/Time    : 2023/07/28 22:39:47
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

-- フォロー中ユーザー
-- * BackupToTempTable
DROP TABLE if exists following_user CASCADE;

-- * RestoreFromTempTable
CREATE TABLE following_user (
  id UUID DEFAULT GEN_RANDOM_UUID() NOT NULL
  , followee_id UUID NOT NULL
  , follower_id UUID NOT NULL
  , created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , CONSTRAINT following_user_PKC PRIMARY KEY (id)
) ;

-- 買い物明細
-- * BackupToTempTable
DROP TABLE if exists shopping_item CASCADE;

-- * RestoreFromTempTable
CREATE TABLE shopping_item (
  id UUID DEFAULT GEN_RANDOM_UUID() NOT NULL
  , shopping_list_id UUID NOT NULL
  , ingredient_id UUID
  , idx INTEGER NOT NULL
  , created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , CONSTRAINT shopping_item_PKC PRIMARY KEY (id)
) ;

-- 材料
-- * BackupToTempTable
DROP TABLE if exists ingredient CASCADE;

-- * RestoreFromTempTable
CREATE TABLE ingredient (
  id UUID DEFAULT GEN_RANDOM_UUID() NOT NULL
  , recipe_id UUID NOT NULL
  , idx INTEGER NOT NULL
  , name TEXT NOT NULL
  , supplement TEXT
  , created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
  , CONSTRAINT ingredient_PKC PRIMARY KEY (id)
) ;

-- ファボ／リム履歴
-- * BackupToTempTable
DROP TABLE if exists fav_history CASCADE;

-- * RestoreFromTempTable
CREATE TABLE fav_history (
  recipe_id UUID NOT NULL
  , usr_id UUID NOT NULL
  , is_fav BOOLEAN NOT NULL
  , created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
) ;

-- フォロー／リム履歴
-- * BackupToTempTable
DROP TABLE if exists follow_chef_history CASCADE;

-- * RestoreFromTempTable
CREATE TABLE follow_chef_history (
  chef_id UUID NOT NULL
  , usr_id UUID NOT NULL
  , is_follow BOOLEAN NOT NULL
  , created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
) ;

-- 買い物リスト
-- * BackupToTempTable
DROP TABLE if exists shopping_list CASCADE;

-- * RestoreFromTempTable
CREATE TABLE shopping_list (
  id UUID DEFAULT GEN_RANDOM_UUID() NOT NULL
  , usr_id UUID NOT NULL
  , recipe_id UUID
  , description TEXT
  , is_fair_copy BOOLEAN NOT NULL
  , created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , CONSTRAINT shopping_list_PKC PRIMARY KEY (id)
) ;

-- ファボ中
-- * BackupToTempTable
DROP TABLE if exists favoring CASCADE;

-- * RestoreFromTempTable
CREATE TABLE favoring (
  id UUID DEFAULT GEN_RANDOM_UUID() NOT NULL
  , recipe_id UUID NOT NULL
  , usr_id UUID NOT NULL
  , created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , CONSTRAINT favoring_PKC PRIMARY KEY (id)
) ;

-- フォロー中シェフ
-- * BackupToTempTable
DROP TABLE if exists following_chef CASCADE;

-- * RestoreFromTempTable
CREATE TABLE following_chef (
  id UUID DEFAULT GEN_RANDOM_UUID() NOT NULL
  , chef_id UUID NOT NULL
  , usr_id UUID NOT NULL
  , created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , CONSTRAINT following_chef_PKC PRIMARY KEY (id)
) ;

-- ユーザー（一般シェフ）
-- * BackupToTempTable
DROP TABLE if exists usr CASCADE;

-- * RestoreFromTempTable
CREATE TABLE usr (
  id UUID DEFAULT GEN_RANDOM_UUID() NOT NULL
  , email TEXT NOT NULL
  , name TEXT NOT NULL
  , image_url TEXT
  , profile TEXT
  , link type_chef_link[] DEFAULT ARRAY[]::type_chef_link[] NOT NULL
  , auth_server TEXT NOT NULL
  , auth_userinfo JSONB NOT NULL
  , created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , num_recipe INTEGER DEFAULT 0 NOT NULL
  , CONSTRAINT usr_PKC PRIMARY KEY (id)
) ;

CREATE UNIQUE INDEX usr_IX1
  ON usr(email);

-- シェフのレシピ＆マイレシピ
-- * BackupToTempTable
DROP TABLE if exists recipe CASCADE;

-- * RestoreFromTempTable
CREATE TABLE recipe (
  id UUID DEFAULT GEN_RANDOM_UUID() NOT NULL
  , chef_id UUID
  , usr_id UUID
  , name TEXT NOT NULL
  , servings INTEGER NOT NULL
  , method type_recipe_method[] DEFAULT ARRAY[]::type_recipe_method[] NOT NULL
  , image_url TEXT
  , introduction TEXT
  , link TEXT[] DEFAULT ARRAY[]::TEXT[] NOT NULL
  , access_level INTEGER NOT NULL
  , created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , num_fav INTEGER DEFAULT 0 NOT NULL
  , CONSTRAINT recipe_PKC PRIMARY KEY (id)
) ;

-- シェフ
-- * BackupToTempTable
DROP TABLE if exists chef CASCADE;

-- * RestoreFromTempTable
CREATE TABLE chef (
  id UUID DEFAULT GEN_RANDOM_UUID() NOT NULL
  , email TEXT
  , name TEXT NOT NULL
  , image_url TEXT
  , profile TEXT
  , link type_chef_link[] DEFAULT ARRAY[]::type_chef_link[] NOT NULL
  , auth_server TEXT
  , auth_userinfo JSONB
  , created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , num_recipe INTEGER DEFAULT 0 NOT NULL
  , num_follower INTEGER DEFAULT 0 NOT NULL
  , CONSTRAINT chef_PKC PRIMARY KEY (id)
) ;

CREATE UNIQUE INDEX chef_IX1
  ON chef(email);

ALTER TABLE shopping_list
  ADD CONSTRAINT shopping_list_FK1 FOREIGN KEY (recipe_id) REFERENCES recipe(id)
  on delete set null;

ALTER TABLE following_user
  ADD CONSTRAINT following_user_FK1 FOREIGN KEY (follower_id) REFERENCES usr(id)
  on delete cascade;

ALTER TABLE following_user
  ADD CONSTRAINT following_user_FK2 FOREIGN KEY (followee_id) REFERENCES usr(id)
  on delete cascade;

ALTER TABLE shopping_item
  ADD CONSTRAINT shopping_item_FK1 FOREIGN KEY (ingredient_id) REFERENCES ingredient(id)
  on delete set null;

ALTER TABLE shopping_item
  ADD CONSTRAINT shopping_item_FK2 FOREIGN KEY (shopping_list_id) REFERENCES shopping_list(id)
  on delete cascade;

ALTER TABLE ingredient
  ADD CONSTRAINT ingredient_FK1 FOREIGN KEY (recipe_id) REFERENCES recipe(id)
  on delete cascade;

ALTER TABLE shopping_list
  ADD CONSTRAINT shopping_list_FK2 FOREIGN KEY (usr_id) REFERENCES usr(id)
  on delete cascade;

ALTER TABLE fav_history
  ADD CONSTRAINT fav_history_FK1 FOREIGN KEY (usr_id) REFERENCES usr(id)
  on delete cascade;

ALTER TABLE fav_history
  ADD CONSTRAINT fav_history_FK2 FOREIGN KEY (recipe_id) REFERENCES recipe(id)
  on delete cascade;

ALTER TABLE follow_chef_history
  ADD CONSTRAINT follow_chef_history_FK1 FOREIGN KEY (usr_id) REFERENCES usr(id)
  on delete cascade;

ALTER TABLE follow_chef_history
  ADD CONSTRAINT follow_chef_history_FK2 FOREIGN KEY (chef_id) REFERENCES chef(id)
  on delete cascade;

ALTER TABLE recipe
  ADD CONSTRAINT recipe_FK1 FOREIGN KEY (usr_id) REFERENCES usr(id)
  on delete cascade;

ALTER TABLE recipe
  ADD CONSTRAINT recipe_FK2 FOREIGN KEY (chef_id) REFERENCES chef(id)
  on delete cascade;

ALTER TABLE favoring
  ADD CONSTRAINT favoring_FK1 FOREIGN KEY (recipe_id) REFERENCES recipe(id)
  on delete cascade;

ALTER TABLE favoring
  ADD CONSTRAINT favoring_FK2 FOREIGN KEY (usr_id) REFERENCES usr(id)
  on delete cascade;

ALTER TABLE following_chef
  ADD CONSTRAINT following_chef_FK1 FOREIGN KEY (usr_id) REFERENCES usr(id)
  on delete cascade;

ALTER TABLE following_chef
  ADD CONSTRAINT following_chef_FK2 FOREIGN KEY (chef_id) REFERENCES chef(id)
  on delete cascade;

COMMENT ON TABLE following_user IS 'フォロー中ユーザー:Table to show favorite chefs of a user';
COMMENT ON COLUMN following_user.id IS '';
COMMENT ON COLUMN following_user.followee_id IS '';
COMMENT ON COLUMN following_user.follower_id IS '';
COMMENT ON COLUMN following_user.created_at IS '';

COMMENT ON TABLE shopping_item IS '買い物明細';
COMMENT ON COLUMN shopping_item.id IS '';
COMMENT ON COLUMN shopping_item.shopping_list_id IS '';
COMMENT ON COLUMN shopping_item.ingredient_id IS '';
COMMENT ON COLUMN shopping_item.idx IS 'インデックス';
COMMENT ON COLUMN shopping_item.created_at IS '';
COMMENT ON COLUMN shopping_item.updated_at IS '';

COMMENT ON TABLE ingredient IS '材料';
COMMENT ON COLUMN ingredient.id IS '';
COMMENT ON COLUMN ingredient.recipe_id IS '';
COMMENT ON COLUMN ingredient.idx IS 'インデックス';
COMMENT ON COLUMN ingredient.name IS '材料名';
COMMENT ON COLUMN ingredient.supplement IS '補足';
COMMENT ON COLUMN ingredient.created_at IS '';
COMMENT ON COLUMN ingredient.updated_at IS '';

COMMENT ON TABLE fav_history IS 'ファボ／リム履歴';
COMMENT ON COLUMN fav_history.recipe_id IS '';
COMMENT ON COLUMN fav_history.usr_id IS '';
COMMENT ON COLUMN fav_history.is_fav IS 'ファボ／リム';
COMMENT ON COLUMN fav_history.created_at IS '';

COMMENT ON TABLE follow_chef_history IS 'フォロー／リム履歴';
COMMENT ON COLUMN follow_chef_history.chef_id IS '';
COMMENT ON COLUMN follow_chef_history.usr_id IS '';
COMMENT ON COLUMN follow_chef_history.is_follow IS 'フォロー／リム';
COMMENT ON COLUMN follow_chef_history.created_at IS '';

COMMENT ON TABLE shopping_list IS '買い物リスト';
COMMENT ON COLUMN shopping_list.id IS '';
COMMENT ON COLUMN shopping_list.usr_id IS '';
COMMENT ON COLUMN shopping_list.recipe_id IS 'NULL：メモリスト／削除レシピ';
COMMENT ON COLUMN shopping_list.description IS '「*人前」「メモリスト」';
COMMENT ON COLUMN shopping_list.is_fair_copy IS '清書or下書き';
COMMENT ON COLUMN shopping_list.created_at IS '';
COMMENT ON COLUMN shopping_list.updated_at IS '';

COMMENT ON TABLE favoring IS 'ファボ中';
COMMENT ON COLUMN favoring.id IS '';
COMMENT ON COLUMN favoring.recipe_id IS '';
COMMENT ON COLUMN favoring.usr_id IS '';
COMMENT ON COLUMN favoring.created_at IS '';

COMMENT ON TABLE following_chef IS 'フォロー中シェフ:Table to show favorite chefs of a user';
COMMENT ON COLUMN following_chef.id IS '';
COMMENT ON COLUMN following_chef.chef_id IS '';
COMMENT ON COLUMN following_chef.usr_id IS '';
COMMENT ON COLUMN following_chef.created_at IS '';

COMMENT ON TABLE usr IS 'ユーザー（一般シェフ）';
COMMENT ON COLUMN usr.id IS '';
COMMENT ON COLUMN usr.email IS 'ログインemail';
COMMENT ON COLUMN usr.name IS 'ニックネーム';
COMMENT ON COLUMN usr.image_url IS 'プロフィール画像（任意）';
COMMENT ON COLUMN usr.profile IS '自己紹介（任意）';
COMMENT ON COLUMN usr.link IS 'リンク（任意）';
COMMENT ON COLUMN usr.auth_server IS 'google／apple';
COMMENT ON COLUMN usr.auth_userinfo IS '認証USERINFO';
COMMENT ON COLUMN usr.created_at IS '';
COMMENT ON COLUMN usr.updated_at IS '';
COMMENT ON COLUMN usr.num_recipe IS 'マイレシピ数';

COMMENT ON TABLE recipe IS 'シェフのレシピ＆マイレシピ';
COMMENT ON COLUMN recipe.id IS '';
COMMENT ON COLUMN recipe.chef_id IS '';
COMMENT ON COLUMN recipe.usr_id IS '';
COMMENT ON COLUMN recipe.name IS 'レシピ名';
COMMENT ON COLUMN recipe.servings IS '＊人前';
COMMENT ON COLUMN recipe.method IS '作り方';
COMMENT ON COLUMN recipe.image_url IS '画像';
COMMENT ON COLUMN recipe.introduction IS 'レシピの紹介文';
COMMENT ON COLUMN recipe.link IS 'リンク';
COMMENT ON COLUMN recipe.access_level IS '公開等:公開、限定公開、非公開、下書き';
COMMENT ON COLUMN recipe.created_at IS '';
COMMENT ON COLUMN recipe.updated_at IS '';
COMMENT ON COLUMN recipe.num_fav IS 'ファボられ数';

COMMENT ON TABLE chef IS 'シェフ';
COMMENT ON COLUMN chef.id IS '';
COMMENT ON COLUMN chef.email IS 'ログインemail';
COMMENT ON COLUMN chef.name IS '登録名';
COMMENT ON COLUMN chef.image_url IS 'プロフィール画像';
COMMENT ON COLUMN chef.profile IS 'シェフ紹介';
COMMENT ON COLUMN chef.link IS 'リンク';
COMMENT ON COLUMN chef.auth_server IS 'google／apple';
COMMENT ON COLUMN chef.auth_userinfo IS '認証USERINFO';
COMMENT ON COLUMN chef.created_at IS '';
COMMENT ON COLUMN chef.updated_at IS '';
COMMENT ON COLUMN chef.num_recipe IS 'レシピ数';
COMMENT ON COLUMN chef.num_follower IS 'フォロワー数';

--
-- VIEW作成
-- 参考：https://qiita.com/kanedaq/items/8b47443df2c42a15eaa9
--

DROP VIEW if exists v_usr CASCADE;

CREATE VIEW v_usr AS
SELECT
    id,
    email,
    name,
    image_url,
    profile,
    TO_JSONB(link) AS link,
    auth_server,
    auth_userinfo,
    created_at,
    updated_at,
    num_recipe
FROM
    usr;

DROP VIEW if exists v_recipe CASCADE;

CREATE VIEW v_recipe AS
SELECT
    id,
    chef_id,
    usr_id,
    name,
    servings,
    TO_JSONB(method) AS method,
    image_url,
    Introduction,
    link,
    access_level,
    created_at,
    updated_at
FROM
    recipe;

DROP VIEW if exists v_chef CASCADE;

CREATE VIEW v_chef AS
SELECT
    id,
    email,
    name,
    image_url,
    profile,
    TO_JSONB(link) AS link,
    auth_server,
    auth_userinfo,
    created_at,
    updated_at,
    num_recipe,
    num_follower
FROM
    chef;

--
-- TRIGGER関連作成
-- 参考：https://blog.kumano-te.com/activities/auto-update-last-updated-at-postgresql
--

DROP FUNCTION if exists refresh_updated_at CASCADE;

CREATE OR REPLACE FUNCTION refresh_updated_at() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := NOW();
    return NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER t_shopping_item_updated_at
    BEFORE UPDATE ON shopping_item
    FOR EACH ROW
    EXECUTE FUNCTION refresh_updated_at();


CREATE TRIGGER t_ingredient_updated_at
    BEFORE UPDATE ON ingredient
    FOR EACH ROW
    EXECUTE FUNCTION refresh_updated_at();


CREATE TRIGGER t_shopping_list_updated_at
    BEFORE UPDATE ON shopping_list
    FOR EACH ROW
    EXECUTE FUNCTION refresh_updated_at();


CREATE TRIGGER t_usr_updated_at
    BEFORE UPDATE ON usr
    FOR EACH ROW
    EXECUTE FUNCTION refresh_updated_at();


CREATE TRIGGER t_recipe_updated_at
    BEFORE UPDATE ON recipe
    FOR EACH ROW
    EXECUTE FUNCTION refresh_updated_at();


CREATE TRIGGER t_chef_updated_at
    BEFORE UPDATE ON chef
    FOR EACH ROW
    EXECUTE FUNCTION refresh_updated_at();

DROP FUNCTION if exists refresh_num_recipe CASCADE;

CREATE OR REPLACE FUNCTION refresh_num_recipe() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
        IF OLD.chef_id IS NOT NULL THEN
            UPDATE chef SET
                num_recipe = (
                    SELECT
                        COUNT(*)
                    FROM
                        recipe
                    WHERE
                        chef_id = OLD.chef_id
                )
            WHERE
                id = OLD.chef_id;
        END IF;

        IF OLD.usr_id IS NOT NULL THEN
            UPDATE usr SET
                num_recipe = (
                    SELECT
                        COUNT(*)
                    FROM
                        recipe
                    WHERE
                        usr_id = OLD.usr_id
                )
            WHERE
                id = OLD.usr_id;
        END IF;
    ELSEIF TG_OP = 'INSERT' THEN
        IF NEW.chef_id IS NOT NULL THEN
            UPDATE chef SET
                num_recipe = (
                    SELECT
                        COUNT(*)
                    FROM
                        recipe
                    WHERE
                        chef_id = NEW.chef_id
                )
            WHERE
                id = NEW.chef_id;
        END IF;

        IF NEW.usr_id IS NOT NULL THEN
            UPDATE usr SET
                num_recipe = (
                    SELECT
                        COUNT(*)
                    FROM
                        recipe
                    WHERE
                        usr_id = NEW.usr_id
                )
            WHERE
                id = NEW.usr_id;
        END IF;
    ELSEIF TG_OP = 'UPDATE' THEN
        IF NEW.chef_id = OLD.chef_id OR (NEW.chef_id IS NULL AND OLD.chef_id IS NULL) THEN
        ELSE
            UPDATE chef SET
                num_recipe = (
                    SELECT
                        COUNT(*)
                    FROM
                        recipe
                    WHERE
                        chef_id = OLD.chef_id
                )
            WHERE
                id = OLD.chef_id;

            UPDATE chef SET
                num_recipe = (
                    SELECT
                        COUNT(*)
                    FROM
                        recipe
                    WHERE
                        chef_id = NEW.chef_id
                )
            WHERE
                id = NEW.chef_id;
        END IF;

        IF NEW.usr_id = OLD.usr_id OR (NEW.usr_id IS NULL AND OLD.usr_id IS NULL) THEN
        ELSE
            UPDATE usr SET
                num_recipe = (
                    SELECT
                        COUNT(*)
                    FROM
                        recipe
                    WHERE
                        usr_id = OLD.usr_id
                )
            WHERE
                id = OLD.usr_id;

            UPDATE usr SET
                num_recipe = (
                    SELECT
                        COUNT(*)
                    FROM
                        recipe
                    WHERE
                        usr_id = NEW.usr_id
                )
            WHERE
                id = NEW.usr_id;
        END IF;
    END IF;

    return NEW;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE TRIGGER t_recipe_num_recipe
    AFTER UPDATE OF chef_id, usr_id OR INSERT OR DELETE ON recipe
    FOR EACH ROW
    EXECUTE FUNCTION refresh_num_recipe();

DROP FUNCTION if exists refresh_follow CASCADE;

CREATE OR REPLACE FUNCTION refresh_follow() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
        INSERT INTO follow_chef_history
        (
            chef_id,
            usr_id,
            is_follow
        )
        VALUES
        (
            OLD.chef_id,
            OLD.usr_id,
            FALSE
        );

        UPDATE chef SET
            num_follower = (
                SELECT
                    COUNT(*)
                FROM
                    following_chef
                WHERE
                    chef_id = OLD.chef_id
            )
        WHERE
            id = OLD.chef_id;
    ELSEIF TG_OP = 'INSERT' THEN
        INSERT INTO follow_chef_history
        (
            chef_id,
            usr_id,
            is_follow
        )
        VALUES
        (
            NEW.chef_id,
            NEW.usr_id,
            TRUE
        );

        UPDATE chef SET
            num_follower = (
                SELECT
                    COUNT(*)
                FROM
                    following_chef
                WHERE
                    chef_id = NEW.chef_id
            )
        WHERE
            id = NEW.chef_id;
    END IF;

    return NEW;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE TRIGGER t_recipe_follow
    AFTER INSERT OR DELETE ON following_chef
    FOR EACH ROW
    EXECUTE FUNCTION refresh_follow();

DROP FUNCTION if exists refresh_fav CASCADE;

CREATE OR REPLACE FUNCTION refresh_fav() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
        INSERT INTO fav_history
        (
            recipe_id,
            usr_id,
            is_fav
        )
        VALUES
        (
            OLD.recipe_id,
            OLD.usr_id,
            FALSE
        );

        UPDATE recipe SET
            num_fav = (
                SELECT
                    COUNT(*)
                FROM
                    favoring
                WHERE
                    recipe_id = OLD.recipe_id
            )
        WHERE
            id = OLD.recipe_id;
    ELSEIF TG_OP = 'INSERT' THEN
        INSERT INTO fav_history
        (
            recipe_id,
            usr_id,
            is_fav
        )
        VALUES
        (
            NEW.recipe_id,
            NEW.usr_id,
            TRUE
        );

        UPDATE recipe SET
            num_fav = (
                SELECT
                    COUNT(*)
                FROM
                    favoring
                WHERE
                    recipe_id = NEW.recipe_id
            )
        WHERE
            id = NEW.recipe_id;
    END IF;

    return NEW;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE TRIGGER t_recipe_fav
    AFTER INSERT OR DELETE ON favoring
    FOR EACH ROW
    EXECUTE FUNCTION refresh_fav();


--
-- ストアドプロシージャ／ファンクション作成
-- 参考：https://qiita.com/kanedaq/items/6c86d8dff79a5fa1dda3
--

DROP FUNCTION if exists insert_chef CASCADE;

CREATE OR REPLACE FUNCTION insert_chef(
    data JSONB
)
    RETURNS v_chef
    LANGUAGE plpgsql
    volatile
AS
$$
DECLARE
    inserting_link type_chef_link[];
    rec v_chef%ROWTYPE;
BEGIN
    inserting_link = ARRAY[]::type_chef_link[];
    FOR i IN 0..JSONB_ARRAY_LENGTH(data->'link') - 1 LOOP
        inserting_link = ARRAY_APPEND(inserting_link,
            ROW(data->'link'->i->>'label',
                data->'link'->i->>'url')::type_chef_link);
    END LOOP;

    INSERT INTO chef
    (
        email,
        name,
        image_url,
        profile,
        link,
        auth_server,
        auth_userinfo
    )
    VALUES
    (
        data->>'email',
        data->>'name',
        data->>'imageUrl',
        data->>'profile',
        inserting_link,
        data->>'authServer',
        data->'authUserinfo'
    )
    RETURNING
        id,
        email,
        name,
        image_url,
        profile,
        TO_JSONB(link) AS link,
        auth_server,
        auth_userinfo,
        created_at,
        updated_at,
        num_recipe,
        num_follower
    INTO
        rec;

    RETURN rec;
END
$$;

DROP FUNCTION if exists update_chef CASCADE;

CREATE OR REPLACE FUNCTION update_chef(
    id UUID,
    data JSONB
)
    RETURNS v_chef
    LANGUAGE plpgsql
    volatile
AS
$$
DECLARE
    updating_link type_chef_link[];
    rec v_chef%ROWTYPE;
BEGIN
    updating_link = ARRAY[]::type_chef_link[];
    FOR i IN 0..JSONB_ARRAY_LENGTH(data->'link') - 1 LOOP
        updating_link = ARRAY_APPEND(updating_link,
            ROW(data->'link'->i->>'label',
                data->'link'->i->>'url')::type_chef_link);
    END LOOP;

    UPDATE chef SET
--         email         = data->>'email',
        name          = data->>'name',
        image_url     = data->>'imageUrl',
        profile       = data->>'profile',
        link          = updating_link
--         auth_server   = data->>'authServer',
--         auth_userinfo = data->'authUserinfo'
    WHERE
        chef.id = update_chef.id
    RETURNING
        chef.id,
        email,
        name,
        image_url,
        profile,
        TO_JSONB(link) AS link,
        auth_server,
        auth_userinfo,
        created_at,
        updated_at,
        num_recipe,
        num_follower
    INTO
        rec;

    RETURN rec;
END
$$;

DROP FUNCTION if exists insert_usr CASCADE;

CREATE OR REPLACE FUNCTION insert_usr(
    data JSONB
)
    RETURNS v_usr
    LANGUAGE plpgsql
    volatile
AS
$$
DECLARE
    inserting_link type_chef_link[];
    rec v_usr%ROWTYPE;
BEGIN
    inserting_link = ARRAY[]::type_chef_link[];
    FOR i IN 0..JSONB_ARRAY_LENGTH(data->'link') - 1 LOOP
        inserting_link = ARRAY_APPEND(inserting_link,
            ROW(data->'link'->i->>'label',
                data->'link'->i->>'url')::type_chef_link);
    END LOOP;

    INSERT INTO usr
    (
        email,
        name,
        image_url,
        profile,
        link,
        auth_server,
        auth_userinfo
    )
    VALUES
    (
        data->>'email',
        data->>'name',
        data->>'imageUrl',
        data->>'profile',
        inserting_link,
        data->>'authServer',
        data->'authUserinfo'
    )
    RETURNING
        id,
        email,
        name,
        image_url,
        profile,
        TO_JSONB(link) AS link,
        auth_server,
        auth_userinfo,
        created_at,
        updated_at,
        num_recipe
    INTO
        rec;

    RETURN rec;
END
$$;

DROP FUNCTION if exists update_usr CASCADE;

CREATE OR REPLACE FUNCTION update_usr(
    id UUID,
    data JSONB
)
    RETURNS v_usr
    LANGUAGE plpgsql
    volatile
AS
$$
DECLARE
    updating_link type_chef_link[];
    rec v_usr%ROWTYPE;
BEGIN
    updating_link = ARRAY[]::type_chef_link[];
    FOR i IN 0..JSONB_ARRAY_LENGTH(data->'link') - 1 LOOP
        updating_link = ARRAY_APPEND(updating_link,
            ROW(data->'link'->i->>'label',
                data->'link'->i->>'url')::type_chef_link);
    END LOOP;

    UPDATE usr SET
--         email         = data->>'email',
        name          = data->>'name',
        image_url     = data->>'imageUrl',
        profile       = data->>'profile',
        link          = updating_link
--         auth_server   = data->>'authServer',
--         auth_userinfo = data->'authUserinfo'
    WHERE
        usr.id = update_usr.id
    RETURNING
        usr.id,
        email,
        name,
        image_url,
        profile,
        TO_JSONB(link) AS link,
        auth_server,
        auth_userinfo,
        created_at,
        updated_at,
        num_recipe
    INTO
        rec;

    RETURN rec;
END
$$;

DROP FUNCTION if exists insert_recipe CASCADE;

CREATE OR REPLACE FUNCTION insert_recipe(
    data JSONB
)
    RETURNS v_recipe
    LANGUAGE plpgsql
    volatile
AS
$$
DECLARE
    inserting_method type_recipe_method[];
    rec v_recipe%ROWTYPE;
BEGIN
    inserting_method = ARRAY[]::type_recipe_method[];
    FOR i IN 0..JSONB_ARRAY_LENGTH(data->'method') - 1 LOOP
        inserting_method = ARRAY_APPEND(inserting_method,
            ROW(data->'method'->i->>'html',
                data->'method'->i->>'supplement')::type_recipe_method);
    END LOOP;

    INSERT INTO recipe
    (
        chef_id,
        usr_id,
        name,
        servings,
        method,
        image_url,
        introduction,
        link,
        access_level
    )
    VALUES
    (
        (data->>'chefId')::UUID,
        (data->>'usrId')::UUID,
        data->>'name',
        (data->'servings')::INTEGER,
        inserting_method,
        data->>'imageUrl',
        data->>'introduction',
        (SELECT ARRAY_AGG(value) FROM JSONB_ARRAY_ELEMENTS_TEXT(data->'link')),
        (data->'accessLevel')::INTEGER
    )
    RETURNING
        id,
        chef_id,
        usr_id,
        name,
        servings,
        TO_JSONB(method) AS method,
        image_url,
        introduction,
        link,
        access_level,
        created_at,
        updated_at
    INTO
        rec;

    RETURN rec;
END
$$;

DROP FUNCTION if exists update_recipe CASCADE;

CREATE OR REPLACE FUNCTION update_recipe(
    id UUID,
    data JSONB
)
    RETURNS v_recipe
    LANGUAGE plpgsql
    volatile
AS
$$
DECLARE
    updating_method type_recipe_method[];
    rec v_recipe%ROWTYPE;
BEGIN
    updating_method = ARRAY[]::type_recipe_method[];
    FOR i IN 0..JSONB_ARRAY_LENGTH(data->'method') - 1 LOOP
        updating_method = ARRAY_APPEND(updating_method,
            ROW(data->'method'->i->>'html',
                data->'method'->i->>'supplement')::type_recipe_method);
    END LOOP;

    UPDATE recipe SET
--         chef_id      = (data->>'chefId')::UUID,
--         usr_id       = (data->>'usrId')::UUID,
        name         = data->>'name',
        servings     = (data->'servings')::INTEGER,
        method       = updating_method,
        image_url    = data->>'imageUrl',
        introduction = data->>'introduction',
        link         = (SELECT ARRAY_AGG(value) FROM JSONB_ARRAY_ELEMENTS_TEXT(data->'link')),
        access_level = (data->'accessLevel')::INTEGER
    WHERE
        recipe.id = update_recipe.id
    RETURNING
        recipe.id,
        chef_id,
        usr_id,
        name,
        servings,
        TO_JSONB(method) AS method,
        image_url,
        introduction,
        link,
        access_level,
        created_at,
        updated_at
    INTO
        rec;

    RETURN rec;
END
$$;
