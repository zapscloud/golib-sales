package sales_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-sales/sales_repository/mongodb_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// PageDao - Card DAO Repository
type PageDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)
	//List - List all Collections
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Get - Get by code
	Get(pageId string) (utils.Map, error)
	// Find - Find by filter
	Find(filter string) (utils.Map, error)
	// Create - Create Collection
	Create(indata utils.Map) (utils.Map, error)
	// Update - Update Collection
	Update(pageId string, indata utils.Map) (utils.Map, error)
	// Delete - Delete Collection
	Delete(pageId string) (int64, error)
}

// NewPageDao - Contruct Business Page Dao
func NewPageDao(client utils.Map, business_id string) PageDao {
	var daoPage PageDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoPage = &mongodb_repository.PageMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoPage != nil {
		// Initialize the Dao
		daoPage.InitializeDao(client, business_id)
	}

	return daoPage
}
