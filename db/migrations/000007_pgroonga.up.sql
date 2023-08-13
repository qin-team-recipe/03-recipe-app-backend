DROP INDEX if exists pgroonga_chef_name_index;
DROP INDEX if exists pgroonga_chef_profile_index;
DROP INDEX if exists pgroonga_recipe_name_index;
DROP INDEX if exists pgroonga_recipe_introduction_index;

SET enable_seqscan = on;

CREATE INDEX pgroonga_chef_name_index ON chef USING pgroonga (name);
CREATE INDEX pgroonga_chef_profile_index ON chef USING pgroonga (profile);

CREATE INDEX pgroonga_recipe_name_index ON recipe USING pgroonga (name);
CREATE INDEX pgroonga_recipe_introduction_index ON recipe USING pgroonga (introduction);
