package elementals

import (
	"database/sql"
	"fmt"

	"github.com/gpessoni/compiled/application/dto"
	compile "github.com/gpessoni/compiled/compile"
)

func GetAllCompiledText(db *sql.DB, elementalId string, authUserId, token, format, groupBy, fields string) (dto.CompiledList, error) {
	if fields == "" {
		return dto.CompiledList{}, fmt.Errorf("Error, no fields specified")
	}
	response, err := compile.PrepareResponseElemental(db, elementalId, authUserId, token, format, groupBy, fields)
	if err != nil {
		return dto.CompiledList{}, err
	}
	return response.(dto.CompiledList), nil
}

func GetAllCompiledJson(db *sql.DB, elementalId string, authUserId, token, format, groupBy, fields string) (map[string]interface{}, error) {
	if fields == "" {
		return nil, fmt.Errorf("Error, no fields specified")
	}
	response, err := compile.PrepareResponseElemental(db, elementalId, authUserId, token, format, groupBy, fields)
	if err != nil {
		return nil, err
	}
	return response.(map[string]interface{}), nil
}
