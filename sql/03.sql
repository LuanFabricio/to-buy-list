CREATE TABLE IF NOT EXISTS buy_list_access(
	buy_list_id SERIAL,
	user_id SERIAL,
	CONSTRAINT fk_user_id
		FOREIGN KEY (user_id)
		REFERENCES users(id),
	CONSTRAINT fk_buy_list_id
		FOREIGN KEY (buy_list_id)
		REFERENCES buy_list(id))
