CREATE TABLE IF NOT EXISTS buy_list(
	id SERIAL PRIMARY KEY,
	user_id SERIAL,
	CONSTRAINT fk_user_id
		FOREIGN KEY (user_id)
		REFERENCES users(id))
