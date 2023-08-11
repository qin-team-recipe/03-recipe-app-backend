DROP TABLE if exists shopping_item;
DROP TABLE if exists shopping_list;


-- 買い物リスト
-- * BackupToTempTable
DROP TABLE if exists shopping_list CASCADE;

-- * RestoreFromTempTable
CREATE TABLE shopping_list (
  id UUID DEFAULT GEN_RANDOM_UUID() NOT NULL
  , usr_id UUID NOT NULL
  , recipe_id UUID
  , r_idx INTEGER NOT NULL
  , description TEXT
  , is_fair_copy BOOLEAN NOT NULL
  , created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , CONSTRAINT shopping_list_PKC PRIMARY KEY (id)
) ;

CREATE UNIQUE INDEX shopping_list_IX1
  ON shopping_list(usr_id,recipe_id);

-- 買い物明細
-- * BackupToTempTable
DROP TABLE if exists shopping_item CASCADE;

-- * RestoreFromTempTable
CREATE TABLE shopping_item (
  id UUID DEFAULT GEN_RANDOM_UUID() NOT NULL
  , shopping_list_id UUID NOT NULL
  , ingredient_id UUID
  , idx INTEGER NOT NULL
  , name TEXT NOT NULL
  , supplement TEXT
  , created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , CONSTRAINT shopping_item_PKC PRIMARY KEY (id)
) ;

-- ALTER TABLE

ALTER TABLE shopping_list
  ADD CONSTRAINT shopping_list_FK1 FOREIGN KEY (recipe_id) REFERENCES recipe(id)
  on delete set null;

ALTER TABLE shopping_item
  ADD CONSTRAINT shopping_item_FK1 FOREIGN KEY (ingredient_id) REFERENCES ingredient(id)
  on delete set null;

ALTER TABLE shopping_item
  ADD CONSTRAINT shopping_item_FK2 FOREIGN KEY (shopping_list_id) REFERENCES shopping_list(id)
  on delete cascade;

ALTER TABLE shopping_list
  ADD CONSTRAINT shopping_list_FK2 FOREIGN KEY (usr_id) REFERENCES usr(id)
  on delete cascade;

COMMENT ON TABLE shopping_item IS '買い物明細';
COMMENT ON COLUMN shopping_item.id IS '';
COMMENT ON COLUMN shopping_item.shopping_list_id IS '';
COMMENT ON COLUMN shopping_item.ingredient_id IS '';
COMMENT ON COLUMN shopping_item.idx IS 'インデックス';
COMMENT ON COLUMN shopping_item.name IS '材料名';
COMMENT ON COLUMN shopping_item.supplement IS '補足';
COMMENT ON COLUMN shopping_item.created_at IS '';
COMMENT ON COLUMN shopping_item.updated_at IS '';

COMMENT ON TABLE shopping_list IS '買い物リスト';
COMMENT ON COLUMN shopping_list.id IS '';
COMMENT ON COLUMN shopping_list.usr_id IS '';
COMMENT ON COLUMN shopping_list.recipe_id IS 'NULL：メモリスト／削除レシピ';
COMMENT ON COLUMN shopping_list.r_idx IS 'リバースインデックス';
COMMENT ON COLUMN shopping_list.description IS '「*人前」「メモリスト」';
COMMENT ON COLUMN shopping_list.is_fair_copy IS '清書or下書き';
COMMENT ON COLUMN shopping_list.created_at IS '';
COMMENT ON COLUMN shopping_list.updated_at IS '';
