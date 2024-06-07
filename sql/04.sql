CREATE TABLE IF NOT EXISTS items (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255),
	current_quantity INTEGER CHECK(current_quantity >= 0),
	min_quantity INTEGER CHECK(min_quantity >= 0),
	send_email BOOLEAN,
	buy_list_id SERIAL,
	CONSTRAINT fk_buy_list_id
		FOREIGN KEY (buy_list_id)
		REFERENCES buy_list(id)
		ON DELETE CASCADE)
