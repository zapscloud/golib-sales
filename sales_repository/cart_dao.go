package sales_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-sales/sales_repository/mongodb_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// CartDao - Card DAO Repository
type CartDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string, customerId string)
	//List - List all Collections
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Get - Get by code
	Get(cartId string) (utils.Map, error)
	// Find - Find by filter
	Find(filter string) (utils.Map, error)
	// Create - Create Collection
	Create(indata utils.Map) (utils.Map, error)
	// Update - Update Collection
	Update(cartId string, indata utils.Map) (utils.Map, error)
	// Delete - Delete Collection
	Delete(cartId string) (int64, error)
}

// NewCartDao - Contruct Business Cart Dao
func NewCartDao(client utils.Map, businessId string, customerId string) CartDao {
	var daoCart CartDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoCart = &mongodb_repository.CartMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoCart != nil {
		// Initialize the Dao
		daoCart.InitializeDao(client, businessId, customerId)
	}

	return daoCart
}
