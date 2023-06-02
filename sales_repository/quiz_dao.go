package sales_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-sales/sales_repository/mongodb_repository"

	"github.com/zapscloud/golib-utils/utils"
)

// QuizDao - Card DAO Repository
type QuizDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)
	//List - List all Collections
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Get - Get by code
	Get(quizId string) (utils.Map, error)
	// Find - Find by filter
	Find(filter string) (utils.Map, error)
	// Create - Create Collection
	Create(indata utils.Map) (utils.Map, error)
	// Update - Update Collection
	Update(quizId string, indata utils.Map) (utils.Map, error)
	// Delete - Delete Collection
	Delete(quizId string) (int64, error)
}

// NewQuizDao - Contruct Business Quiz Dao
func NewQuizDao(client utils.Map, business_id string) QuizDao {
	var daoQuiz QuizDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoQuiz = &mongodb_repository.QuizMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoQuiz != nil {
		// Initialize the Dao
		daoQuiz.InitializeDao(client, business_id)
	}

	return daoQuiz
}
