package persistence

import (
	"database/sql"
	"fmt"

	"github.com/gpessoni/compiled/application/dto"
)

type ElementalPersistence struct {
	db *sql.DB
}

func NewElementalPersistence(db *sql.DB) *ElementalPersistence {
	return &ElementalPersistence{db: db}
}

func (ep *ElementalPersistence) FindById(id string) (dto.Elemental, error) {
	query := `SELECT id, user_id, template, description, title, elemental_type_id, is_premium
		FROM prompt WHERE id = $1`
	row := ep.db.QueryRow(query, id)

	var elemental dto.Elemental
	err := row.Scan(&elemental.Id, &elemental.UserId, &elemental.Template, &elemental.Description, &elemental.Title,
		&elemental.ElementalTypeId, &elemental.IsPremium)
	if err != nil {
		fmt.Print(err)
		return dto.Elemental{}, fmt.Errorf("failed to find elemental: %w", err)
	}
	return elemental, nil
}
