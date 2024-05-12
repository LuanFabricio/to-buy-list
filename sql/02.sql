CREATE TABLE IF NOT EXISTS buy_list(
	id SERIAL PRIMARY KEY,
	name VARCHAR(64) NOT NULL,
	owner_user_id SERIAL,
	CONSTRAINT fk_owner_user_id
		FOREIGN KEY (owner_user_id)
		REFERENCES users(id))
