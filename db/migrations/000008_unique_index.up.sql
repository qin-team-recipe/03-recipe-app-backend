DROP INDEX if exists following_user_IX1 CASCADE;
DROP INDEX if exists favoring_IX1 CASCADE;
DROP INDEX if exists following_chef_IX1 CASCADE;

CREATE UNIQUE INDEX following_user_IX1
  ON following_user(followee_id,follower_id);

CREATE UNIQUE INDEX favoring_IX1
  ON favoring(recipe_id,usr_id);

CREATE UNIQUE INDEX following_chef_IX1
  ON following_chef(chef_id,usr_id);
