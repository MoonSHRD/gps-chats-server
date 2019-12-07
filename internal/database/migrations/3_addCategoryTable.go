package migrations

const (
	_3Up = `
		CREATE TABLE chatCategories 
		(
			id SERIAL
				CONSTRAINT chatCategories_pk
					PRIMARY KEY,
			categoryName TEXT NOT NULL
		);

		CREATE TABLE roomsChatCategoriesLink 
		(
			id SERIAL
				CONSTRAINT roomsChatCategoriesLink_pk
					PRIMARY KEY,
			categoryId INT NOT NULL,
			roomId INT NOT NULL
		);

		ALTER TABLE rooms 
			DROP COLUMN category;
	`

	_3Down = `
		DROP TABLE chatCategories;
		DROP TABLE roomsChatCategoriesLink;
		ALTER TABLE rooms-name
  			ADD category TEXT;
	`
)
