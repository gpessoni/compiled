package list

import (
	"database/sql"

	"github.com/gpessoni/compiled/application/dto"
	compile "github.com/gpessoni/compiled/compile"
)

func GetAllCompiledText(db *sql.DB, listId int64, authUserId, token, format, groupBy, fields string) (dto.CompiledList, error) {
	response, err := compile.PrepareListResponse(db, listId, authUserId, token, format, groupBy, fields)
	if err != nil {
		return dto.CompiledList{}, err
	}
	return response.(dto.CompiledList), nil
}

func GetAllCompiledJson(db *sql.DB, listId int64, authUserId, token, format, groupBy, fields string) (map[string]interface{}, error) {
	response, err := compile.PrepareListResponse(db, listId, authUserId, token, format, groupBy, fields)
	if err != nil {
		return nil, err
	}
	return response.(map[string]interface{}), nil
}
