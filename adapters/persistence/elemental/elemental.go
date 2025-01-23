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
	query := `SELECT 
    p.id,
    p.user_id,
    p.template,
    p.description,
    p.title,
    p.elemental_type_id,
    p.is_premium,
    COALESCE(p.video, '') AS video,
    COALESCE(p.url, '') AS prompt_url,
    COALESCE(STRING_AGG(pi2.url, ', '), '') AS images,
    p.price AS price,
    COALESCE(
        (SELECT STRING_AGG(
            CONCAT(tutorial_step.title, ': ', 
                REGEXP_REPLACE(tutorial_step.description, '<[^>]+>', '', 'g')
            ), ', '
        )
        FROM tutorial_step
        WHERE tutorial_step.prompt_id = p.id AND p.is_tutorial_hidden = false), 
        ''
    ) AS tutorial
FROM 
    prompt p
LEFT JOIN 
    prompt_image pi2 ON p.id = pi2.prompt_id 
WHERE 
    p.id = $1
GROUP BY 
    p.id, p.user_id, p.template, p.description, p.title,  p.elemental_type_id, p.is_premium, p.video, p.url, p.price
   `
	row := ep.db.QueryRow(query, id)

	var elemental dto.Elemental
	err := row.Scan(&elemental.Id, &elemental.UserId, &elemental.Template, &elemental.Description, &elemental.Title,
		&elemental.ElementalTypeId, &elemental.IsPremium, &elemental.Url, &elemental.Video, &elemental.Images, &elemental.Price, &elemental.Tutorial)
	if err != nil {
		return dto.Elemental{}, fmt.Errorf("failed to find elemental: %w", err)
	}
	return elemental, nil
}
