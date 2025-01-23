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
	query := `SELECT 
		l.id,
		l.title,
		l.description,
		l.is_premium,
		l.is_private,
		l.elemental_type_id,
		l.is_hidden,
		l.table_id,
		l.table_index,
		l.table_orientation,
		COALESCE(l.video, '') AS video,
		COALESCE(STRING_AGG(li.url, ', '), '') AS images,
		l.url AS url,
		l.price,
		COALESCE(
        (SELECT STRING_AGG(
            CONCAT(tutorial_step.title, ': ', 
                REGEXP_REPLACE(tutorial_step.description, '<[^>]+>', '', 'g')
            ), ', '
        )
        FROM tutorial_step
        WHERE tutorial_step.list_id = l.id AND l.is_tutorial_hidden = false), 
        ''
    ) AS tutorial
	FROM 
		list l
	LEFT JOIN 
		list_image li ON l.id = li.list_id
	WHERE 
		l.id = $1
	GROUP BY 
		l.id, 
		l.title, 
		l.description, 
		l.is_premium, 
		l.is_private, 
		l.elemental_type_id, 
		l.is_hidden, 
		l.table_id, 
		l.table_index, 
		l.table_orientation, 
		l.url,
		l.price;
;
`
	row := lp.db.QueryRow(query, id)

	var list dto.List
	err := row.Scan(&list.Id, &list.Title, &list.Description, &list.IsPremium,
		&list.IsPrivate, &list.ElementalTypeId, &list.IsHidden,
		&list.TableId, &list.TableIndex, &list.TableOrientation, &list.Url, &list.Video, &list.Images, &list.Price, &list.Tutorial)
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

	crossedOrientation := "column"
	if groupBy == "column" {
		crossedOrientation = "row"
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
	select 0, null, null, null, null, null, null, null, l.id, l.title, l.description, l.is_premium, l.user_id , l.table_index, 0, l.url, COALESCE(l.video, '') AS video,  (
        SELECT STRING_AGG(li.url, ', ')
        FROM list_image li
        WHERE li.list_id = l.id
    ) AS images, l.price,
		COALESCE(
        (SELECT STRING_AGG(
            CONCAT(tutorial_step.title, ': ', 
                REGEXP_REPLACE(tutorial_step.description, '<[^>]+>', '', 'g')
            ), ', '
        )
        FROM tutorial_step
        WHERE tutorial_step.list_id = l.id AND l.is_tutorial_hidden = false), 
        ''
    ) AS tutorial
		from list l
	where ` + idField + ` = $1 ` + orientation + `

		union all
	
	select lp.list_id, p.id, p.title, p.description, p.template, p.is_premium, p.elemental_type_id, p.user_id, li.id, li.title, li.description, li.is_premium, li.user_id , coalesce(li.table_index, crossed.table_index), fc.level, p.url, COALESCE(p.video, '') AS video,
    (
        SELECT STRING_AGG(pi.url, ', ')
        FROM prompt_image pi
        WHERE pi.prompt_id = p.id
    ) AS images, p.price, COALESCE(
        (SELECT STRING_AGG(
            CONCAT(tutorial_step.title, ': ', 
                REGEXP_REPLACE(tutorial_step.description, '<[^>]+>', '', 'g')
            ), ', '
        )
        FROM tutorial_step
        WHERE tutorial_step.prompt_id = p.id AND p.is_tutorial_hidden = false), 
        ''
    ) AS tutorial
	FROM list_prompt lp
		inner join list l on l.id = lp.list_id
		left join prompt p on lp.prompt_id = p.id
		left join list crossed on crossed.id = p.table_` + crossedOrientation + `
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
		tableIndex      sql.NullInt64
		level           sql.NullInt64
		url             sql.NullString
		video           sql.NullString
		images          sql.NullString
		price           sql.NullInt64
		tutorial        sql.NullString
	)

	rows, err := lp.db.Query(dml, idValue)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	defer rows.Close()

	var listChildren []dto.ListChild

	for rows.Next() {
		if err := rows.Scan(&dblistId, &id, &title, &description, &template, &isPremium, &elementalTypeId, &pUserId, &lId, &lTitle, &lDescription, &lIsPremium, &userId, &tableIndex, &level, &url, &video, &images, &price, &tutorial); err != nil {
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
			TableIndex:      tableIndex.Int64,
			ElementalTypeId: elementalTypeId.Int64,
			Url:             url.String,
			Video:           video.String,
			Images:          images.String,
			Price:           price.Int64,
			Tutorial:        tutorial.String,
		}
		listChildren = append(listChildren, listChild)
	}

	return listChildren, nil
}
