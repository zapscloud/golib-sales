package sales_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-sales/sales_repository/mongodb_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// MediaDao - Card DAO Repository
type MediaDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)
	//List - List all Collections
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Get - Get by code
	Get(mediaId string) (utils.Map, error)
	// Find - Find by filter
	Find(filter string) (utils.Map, error)
	// Create - Create Collection
	Create(indata utils.Map) (utils.Map, error)
	// Update - Update Collection
	Update(mediaId string, indata utils.Map) (utils.Map, error)
	// Delete - Delete Collection
	Delete(mediaId string) (int64, error)
}

// NewMediaDao - Contruct Business Media Dao
func NewMediaDao(client utils.Map, business_id string) MediaDao {
	var daoMedia MediaDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoMedia = &mongodb_repository.MediaMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoMedia != nil {
		// Initialize the Dao
		daoMedia.InitializeDao(client, business_id)
	}

	return daoMedia
}
