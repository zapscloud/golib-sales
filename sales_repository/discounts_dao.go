package sales_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-sales/sales_repository/mongodb_repository"

	"github.com/zapscloud/golib-utils/utils"
)

// DiscountDao - Card DAO Repository
type DiscountDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)
	//List - List all Collections
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Get - Get by code
	Get(discountId string) (utils.Map, error)
	// Find - Find by filter
	Find(filter string) (utils.Map, error)
	// Create - Create Collection
	Create(indata utils.Map) (utils.Map, error)
	// Update - Update Collection
	Update(discountId string, indata utils.Map) (utils.Map, error)
	// Delete - Delete Collection
	Delete(discountId string) (int64, error)
}

// NewDiscountDao - Contruct Business Discount Dao
func NewDiscountDao(client utils.Map, business_id string) DiscountDao {
	var daoDiscount DiscountDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoDiscount = &mongodb_repository.DiscountMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoDiscount != nil {
		// Initialize the Dao
		daoDiscount.InitializeDao(client, business_id)
	}

	return daoDiscount
}
