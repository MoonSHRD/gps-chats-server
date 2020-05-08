package repositories

import (
	"github.com/MoonSHRD/sonis/app"
	"github.com/MoonSHRD/sonis/models"
)

type ChatCategoryRepository struct {
	app *app.App
}

func NewChatCategoryRepository(a *app.App) *ChatCategoryRepository {
	return &ChatCategoryRepository{
		app: a,
	}
}

func (ccr *ChatCategoryRepository) AddCategory(category *models.ChatCategory) (*models.ChatCategory, error) {
	stmt, err := ccr.app.DBConn.Preparex("INSERT INTO chatCategories (categoryName) VALUES ($1) RETURNING id;")
	if err != nil {
		return nil, err
	}

	err = stmt.QueryRow(category.CategoryName).Scan(&category.Id)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (ccr *ChatCategoryRepository) GetCategory(id int) (*models.ChatCategory, error) {
	stmt, err := ccr.app.DBConn.Preparex("SELECT * FROM chatCategories WHERE id = ?;")
	if err != nil {
		return nil, err
	}
	var category models.ChatCategory
	err = stmt.Select(&category)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (ccr *ChatCategoryRepository) GetAllCategories() ([]models.ChatCategory, error) {
	stmt, err := ccr.app.DBConn.Preparex("SELECT * FROM chatCategories;")
	if err != nil {
		return nil, err
	}
	var categories []models.ChatCategory
	err = stmt.Select(&categories)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (ccr *ChatCategoryRepository) RemoveCategory(id int) error {
	stmt, err := ccr.app.DBConn.Preparex("DELETE FROM chatCategories WHERE id = ?;")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (ccr *ChatCategoryRepository) UpdateCategoryName(updatedCategory *models.ChatCategory) error {
	stmt, err := ccr.app.DBConn.Preparex("UPDATE chatCategories SET categoryName = $1 WHERE id = $2;")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(updatedCategory.CategoryName, updatedCategory.Id)
	if err != nil {
		return err
	}
	return nil
}
