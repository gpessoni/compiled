package elementals

import (
	"database/sql"

	"github.com/gpessoni/compiled/application/dto"
	compile "github.com/gpessoni/compiled/compile"
)

func GetAllCompiledText(db *sql.DB, elementalId string, authUserId, token, format, groupBy, fields string) (dto.CompiledList, error) {
	response, err := compile.PrepareResponseElemental(db, elementalId, authUserId, token, format, groupBy, fields)
	if err != nil {
		return dto.CompiledList{}, err
	}
	return response.(dto.CompiledList), nil
}

func GetAllCompiledJson(db *sql.DB, elementalId string, authUserId, token, format, groupBy, fields string) (map[string]interface{}, error) {
	response, err := compile.PrepareResponseElemental(db, elementalId, authUserId, token, format, groupBy, fields)
	if err != nil {
		return nil, err
	}
	return response.(map[string]interface{}), nil
}
