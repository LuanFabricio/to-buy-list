CREATE OR REPLACE FUNCTION create_user_access_after_insert()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO buy_list_access(buy_list_id, user_id)
    VALUES (NEW.id, NEW.owner_user_id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql
