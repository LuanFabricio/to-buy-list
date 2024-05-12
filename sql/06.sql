CREATE OR REPLACE TRIGGER after_insert_buy_list_trigger
AFTER INSERT ON buy_list
FOR EACH ROW
EXECUTE FUNCTION create_user_access_after_insert()
