CREATE TABLE IF NOT EXISTS items (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255),
	current_quantity INTEGER CHECK(current_quantity >= 0),
	min_quantity INTEGER CHECK(min_quantity >= 0),
	send_email BOOLEAN)
