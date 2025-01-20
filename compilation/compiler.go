package compilation

import (
	"database/sql"

	"github.com/gpessoni/compiled/application/dto"
	_ "github.com/lib/pq"
)

func GetAllCompiledTextList(db *sql.DB, listId int64, authUserId, token, format, groupBy string) (dto.CompiledList, error) {
	response, err := prepareListResponse(db, listId, authUserId, token, format, groupBy)
	if err != nil {
		return dto.CompiledList{}, err
	}
	return response.(dto.CompiledList), nil
}

func GetAllCompiledTextElemental(db *sql.DB, elementalId string, authUserId, token, format, groupBy string) (dto.CompiledList, error) {
	response, err := prepareResponseElemental(db, elementalId, authUserId, token, format, groupBy)
	if err != nil {
		return dto.CompiledList{}, err
	}
	return response.(dto.CompiledList), nil
}

func GetAllCompiledJsonElemental(db *sql.DB, elementalId string, authUserId, token, format, groupBy string) (map[string]interface{}, error) {
	response, err := prepareResponseElemental(db, elementalId, authUserId, token, format, groupBy)
	if err != nil {
		return nil, err
	}
	return response.(map[string]interface{}), nil
}

func GetAllCompiledJsonList(db *sql.DB, listId int64, authUserId, token, format, groupBy string) (map[string]interface{}, error) {
	response, err := prepareListResponse(db, listId, authUserId, token, format, groupBy)
	if err != nil {
		return nil, err
	}
	return response.(map[string]interface{}), nil
}
