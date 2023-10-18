package sales_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-sales/sales_repository/mongodb_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// CompareProdsDao - Card DAO Repository
type CompareProdsDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)
	//List - List all Collections
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Find - Find by code
	Get(materialTypeid string) (utils.Map, error)
	// Find - Find by filter
	Find(filter string) (utils.Map, error)
	// Create - Create Collection
	Create(indata utils.Map) (utils.Map, error)
	// Update - Update Collection
	Update(categoriId string, indata utils.Map) (utils.Map, error)
	// Delete - Delete Collection
	Delete(materialTypeId string) (int64, error)
}

// NewCompareProdsDao - Contruct Business CompareProds Dao
func NewCompareProdsDao(client utils.Map, business_id string) CompareProdsDao {
	var daoCompareProds CompareProdsDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoCompareProds = &mongodb_repository.CompareProdsMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoCompareProds != nil {
		// Initialize the Dao
		daoCompareProds.InitializeDao(client, business_id)
	}

	return daoCompareProds
}
