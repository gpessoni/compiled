package persistence

import (
	"database/sql"
	"fmt"

	"github.com/gpessoni/compiled/application/dto"
)

type ListPersistence struct {
	db *sql.DB
}

func NewListPersistence(db *sql.DB) *ListPersistence {
	return &ListPersistence{db: db}
}

func (lp *ListPersistence) FindByID(id int64) (dto.List, error) {
	query := `SELECT title, description, price, video, is_premium, is_private, stripe_is_product,
	elemental_type_id, is_hidden, price_original, price_type_id, table_id, table_index, table_orientation
	FROM list WHERE id = $1`
	row := lp.db.QueryRow(query, id)

	var list dto.List
	err := row.Scan(&list.Title, &list.Description, &list.Price, &list.Video, &list.IsPremium,
		&list.IsPrivate, &list.StripeIsProduct, &list.ElementalTypeId, &list.IsHidden,
		&list.PriceOriginal, &list.PriceTypeId, &list.TableId, &list.TableIndex, &list.TableOrientation)
	if err != nil {
		return dto.List{}, fmt.Errorf("failed to find list: %w", err)
	}
	return list, nil
}

func (lp ListPersistence) GetAllCompiledTextList(identifier interface{}, groupBy string) ([]dto.ListChild, error) {
	filter := ""
	justListOwner := false
	if justListOwner {
		filter = "AND l.user_id in (li.user_id, p.user_id)"
	}

	var idField string
	var idValue interface{}
	var orientation string

	switch v := identifier.(type) {
	case string:
		idField = "table_id"
		idValue = v
		orientation = "and table_orientation = 'row'"
		if groupBy == "column" {
			orientation = "and table_orientation = 'column'"
		}
	case int64:
		idField = "id"
		idValue = v
	}

	dml := `
		with recursive folder_content as (
		select *, 0::bigint as list_id, array[id] as path, 1 as level
		from list
		where ` + idField + ` = $1 ` + orientation + `
		--
		union all
		--
		select l.*, lp.list_id, fc.path || l.id, fc.level + 1
		from list l 
			inner join list_prompt lp on lp.list_item_id = l.id
			INNER JOIN folder_content fc ON lp.list_id = fc.id
		WHERE l.id <> ALL(fc.path)
	)
	select 0, null, null, null, null, null, null, null, l.id, l.title, l.description, l.is_premium, l.user_id , 0
		from list l
	where ` + idField + ` = $1 ` + orientation + `

		union all
	
	select lp.list_id, p.id, p.title, p.description, p.template, p.is_premium, p.elemental_type_id, p.user_id, li.id, li.title, li.description, li.is_premium, li.user_id , fc.level
	FROM list_prompt lp
		inner join list l on l.id = lp.list_id
		left join prompt p on lp.prompt_id = p.id
		left join list li on lp.list_item_id = li.id
		inner join folder_content fc on fc.id = lp.list_id
	WHERE 1=1 ` + filter

	var (
		dblistId        sql.NullInt64
		id              sql.NullString
		title           sql.NullString
		description     sql.NullString
		template        sql.NullString
		isPremium       sql.NullBool
		elementalTypeId sql.NullInt64
		pUserId         sql.NullString
		lId             sql.NullInt64
		lTitle          sql.NullString
		lDescription    sql.NullString
		lIsPremium      sql.NullBool
		userId          sql.NullString
		level           sql.NullInt64
	)

	rows, err := lp.db.Query(dml, idValue)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var listChildren []dto.ListChild

	for rows.Next() {
		if err := rows.Scan(&dblistId, &id, &title, &description, &template, &isPremium, &elementalTypeId, &pUserId, &lId, &lTitle, &lDescription, &lIsPremium, &userId, &level); err != nil {
			return nil, err
		}

		isList := false
		if lId.Int64 != 0 {
			isList = true
		}

		listChild := dto.ListChild{
			Id:              id.String,
			LId:             lId.Int64,
			IsList:          isList,
			ListId:          dblistId.Int64,
			Title:           title.String + lTitle.String,
			Description:     description.String + lDescription.String,
			Template:        template.String,
			IsPremium:       isPremium.Bool || lIsPremium.Bool,
			UserId:          userId.String + pUserId.String,
			Level:           level.Int64,
			ElementalTypeId: elementalTypeId.Int64,
		}
		listChildren = append(listChildren, listChild)
	}

	return listChildren, nil
}
