DROP TYPE if exists type_chef_link CASCADE;
DROP TYPE if exists type_recipe_method CASCADE;
DROP TABLE if exists following_user CASCADE;
DROP TABLE if exists shopping_item CASCADE;
DROP TABLE if exists ingredient CASCADE;
DROP TABLE if exists fav_history CASCADE;
DROP TABLE if exists follow_chef_history CASCADE;
DROP TABLE if exists shopping_list CASCADE;
DROP TABLE if exists favoring CASCADE;
DROP TABLE if exists following_chef CASCADE;
DROP TABLE if exists usr CASCADE;
DROP TABLE if exists recipe CASCADE;
DROP TABLE if exists chef CASCADE;
DROP VIEW if exists v_usr CASCADE;
DROP VIEW if exists v_recipe CASCADE;
DROP VIEW if exists v_chef CASCADE;
DROP FUNCTION if exists refresh_updated_at CASCADE;
DROP FUNCTION if exists refresh_num_recipe CASCADE;
DROP FUNCTION if exists refresh_follow CASCADE;
DROP FUNCTION if exists refresh_fav CASCADE;
DROP FUNCTION if exists insert_chef CASCADE;
DROP FUNCTION if exists update_chef CASCADE;
DROP FUNCTION if exists insert_usr CASCADE;
DROP FUNCTION if exists update_usr CASCADE;
DROP FUNCTION if exists insert_recipe CASCADE;
DROP FUNCTION if exists update_recipe CASCADE;
