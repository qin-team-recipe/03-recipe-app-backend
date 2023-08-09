DROP TABLE if exists shopping_item;
DROP TABLE if exists shopping_list;


-- 買い物リスト

-- * RestoreFromTempTable
CREATE TABLE shopping_list (
  id UUID DEFAULT GEN_RANDOM_UUID() NOT NULL
  , usr_id UUID NOT NULL
  , recipe_id UUID
  , idx INTEGER NOT NULL
  , description TEXT
  , is_fair_copy BOOLEAN NOT NULL
  , created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
  , CONSTRAINT shopping_list_PKC PRIMARY KEY (id)
) ;

COMMENT ON COLUMN shopping_list.idx IS 'インデックス';

-- 買い物明細

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
