package sales_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-sales/sales_repository/mongodb_repository"

	"github.com/zapscloud/golib-utils/utils"
)

// CustomerDao - Card DAO Repository
type CustomerDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)
	//List - List all Collections
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Get - Get by code
	Get(customerId string) (utils.Map, error)
	// Find - Find by filter
	Find(filter string) (utils.Map, error)
	// Create - Create Collection
	Create(indata utils.Map) (utils.Map, error)
	// Update - Update Collection
	Update(customerId string, indata utils.Map) (utils.Map, error)
	// Delete - Delete Collection
	Delete(customerId string) (int64, error)

	// Authenticate
	Authenticate(auth_key string, auth_login string, auth_pwd string) (utils.Map, error)
}

// NewCustomerDao - Contruct Business Customer Dao
func NewCustomerDao(client utils.Map, business_id string) CustomerDao {
	var daoCustomer CustomerDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoCustomer = &mongodb_repository.CustomerMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoCustomer != nil {
		// Initialize the Dao
		daoCustomer.InitializeDao(client, business_id)
	}

	return daoCustomer
}
