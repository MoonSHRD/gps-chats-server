package migrations

const (
	_4Up = `
		INSERT INTO chatCategories (categoryName) VALUES ('Тусовки');
		INSERT INTO chatCategories (categoryName) VALUES ('Бизнес ивенты');
		INSERT INTO chatCategories (categoryName) VALUES ('Кружок по интересам');
		INSERT INTO chatCategories (categoryName) VALUES ('Культурные мероприятия');
	`

	_4Down = `
		TRUNCATE TABLE chatCategories;
	`
)
