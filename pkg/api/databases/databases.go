package databases

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/rodrigoTcarmo/nankion/pkg/notion"
)

const (
	databaseId                    = "databaseId"
	missingDatabaseIdErrorMessage = "Missing database ID!"
)

func GetDatabase(ctx *gin.Context) {
	queryDatabaseId, ok := ctx.GetQuery(databaseId)
	if !ok {
		ctx.AbortWithStatusJSON(422, map[string]string{
			"message": missingDatabaseIdErrorMessage,
		})
		return
	}
	var db notion.DatabaseClient = notion.NewDatabase()

	database, code, err := db.SearchDatabases(queryDatabaseId)
	if err != nil {
		ctx.AbortWithStatusJSON(code, map[string]any{
			"message":        err.Error(),
			"Database found": database,
		})
		return
	}

	err = json.NewEncoder(ctx.Writer).Encode(database)
	if err != nil {
		ctx.AbortWithStatusJSON(500, map[string]any{
			"message": "error trying to encode json",
		})
		return
	}

}
