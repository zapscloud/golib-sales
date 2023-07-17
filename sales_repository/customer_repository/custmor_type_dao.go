package customer_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-sales/sales_repository/mongodb_repository/customer_mongodb_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// CustomerTypeDao - Card DAO Repository
type CustomerTypeDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string, customerId string)
	//List - List all Collections
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Get - Get by code
	Get(TypeId string) (utils.Map, error)
	// Find - Find by filter
	Find(filter string) (utils.Map, error)
	// Create - Create Collection
	Create(indata utils.Map) (utils.Map, error)
	// Update - Update Collection
	Update(TypeId string, indata utils.Map) (utils.Map, error)
	// Delete - Delete Collection
	Delete(TypeId string) (int64, error)
}

// NewCustomerTypeDao - Contruct Business Type Dao
func NewCustomerTypeDao(client utils.Map, businessId string, customerId string) CustomerTypeDao {
	var daoType CustomerTypeDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoType = &customer_mongodb_repository.CustomerTypeMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoType != nil {
		// Initialize the Dao
		daoType.InitializeDao(client, businessId, customerId)
	}

	return daoType
}
