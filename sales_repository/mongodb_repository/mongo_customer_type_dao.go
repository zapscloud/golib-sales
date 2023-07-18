package mongodb_repository

import (
	"fmt"
	"log"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/mongo_utils"
	"github.com/zapscloud/golib-sales/sales_common"
	"github.com/zapscloud/golib-utils/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CustomerTypeMongoDBDao - CustomerType DAO Repository
type CustomerTypeMongoDBDao struct {
	client     utils.Map
	businessId string
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func (p *CustomerTypeMongoDBDao) InitializeDao(client utils.Map, businessId string) {
	log.Println("Initialize CustomerType Mongodb DAO")
	p.client = client
	p.businessId = businessId
}

// List - List all Collections
func (t *CustomerTypeMongoDBDao) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {
	var results []utils.Map

	log.Println("Begin - Find All Collection Dao", sales_common.DbCustomerTypes)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(t.client, sales_common.DbCustomerTypes)
	if err != nil {
		return nil, err
	}

	log.Println("Get Collection - Find All Collection Dao", filter, len(filter), sort, len(sort))

	opts := options.Find()

	filterdoc := bson.D{}
	if len(filter) > 0 {
		// filters, _ := strconv.Unquote(string(filter))
		err = bson.UnmarshalExtJSON([]byte(filter), true, &filterdoc)
		if err != nil {
			log.Println("Unmarshal Ext JSON error", err)
			log.Println(filterdoc)
		}
	}

	if len(sort) > 0 {
		var sortdoc interface{}
		err = bson.UnmarshalExtJSON([]byte(sort), true, &sortdoc)
		if err != nil {
			log.Println("Sort Unmarshal Error ", sort)
		} else {
			opts.SetSort(sortdoc)
		}
	}

	if skip > 0 {
		log.Println(filterdoc)
		opts.SetSkip(skip)
	}

	if limit > 0 {
		log.Println(filterdoc)
		opts.SetLimit(limit)
	}
	filterdoc = append(filterdoc,
		bson.E{Key: sales_common.FLD_BUSINESS_ID, Value: t.businessId},
		bson.E{Key: db_common.FLD_IS_DELETED, Value: false})

	log.Println("Parameter values ", filterdoc, opts)
	cursor, err := collection.Find(ctx, filterdoc, opts)
	if err != nil {
		return nil, err
	}

	// get a list of all returned documents and print them out
	// see the mongo.Cursor documentation for more examples of using cursors
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	listdata := []utils.Map{}
	for idx, value := range results {
		log.Println("Item ", idx)
		// Remove fields from result
		value = db_common.AmendFldsForGet(value)
		listdata = append(listdata, value)
	}

	log.Println("Parameter values ", filterdoc)
	filtercount, err := collection.CountDocuments(ctx, filterdoc)
	if err != nil {
		return nil, err
	}

	basefilterdoc := bson.D{
		{Key: sales_common.FLD_BUSINESS_ID, Value: t.businessId},
		{Key: db_common.FLD_IS_DELETED, Value: false}}
	totalcount, err := collection.CountDocuments(ctx, basefilterdoc)
	if err != nil {
		return nil, err
	}

	response := utils.Map{
		db_common.LIST_SUMMARY: utils.Map{
			db_common.LIST_TOTALSIZE:    totalcount,
			db_common.LIST_FILTEREDSIZE: filtercount,
			db_common.LIST_RESULTSIZE:   len(listdata),
		},
		db_common.LIST_RESULT: listdata,
	}

	return response, nil
}

// Get - Get by code
func (t *CustomerTypeMongoDBDao) Get(CustomerTypeId string) (utils.Map, error) {
	// Get a single document
	var result utils.Map

	log.Println("CustomerTypeMongoDBDao::Get:: Begin ", CustomerTypeId)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(t.client, sales_common.DbCustomerTypes)
	log.Println("Get:: Got Collection ")

	filter := bson.D{{Key: sales_common.FLD_CUSTOMER_TYPE_ID, Value: CustomerTypeId}, {}}

	filter = append(filter,
		bson.E{Key: sales_common.FLD_BUSINESS_ID, Value: t.businessId},
		bson.E{Key: db_common.FLD_IS_DELETED, Value: false})

	log.Println("Get:: Got filter ", filter)
	singleResult := collection.FindOne(ctx, filter)
	if singleResult.Err() != nil {
		log.Println("Get:: Record not found ", singleResult.Err())
		return result, singleResult.Err()
	}
	singleResult.Decode(&result)
	if err != nil {
		log.Println("Error in decode", err)
		return result, err
	}

	// Remove fields from result
	result = db_common.AmendFldsForGet(result)

	log.Printf("Business CustomerTypeMongoDBDao::Get:: End Found a single document\n")
	return result, nil
}

// Find - Find by Filter
func (p *CustomerTypeMongoDBDao) Find(filter string) (utils.Map, error) {
	// Find a single document
	var result utils.Map

	log.Println("CustomerTypeDBDao::Find:: Begin ", filter)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, sales_common.DbCustomerTypes)
	log.Println("Find:: Got Collection ", err)

	bfilter := bson.D{}
	err = bson.UnmarshalExtJSON([]byte(filter), true, &bfilter)
	if err != nil {
		fmt.Println("Error on filter Unmarshal", err)
	}
	bfilter = append(bfilter,
		bson.E{Key: sales_common.FLD_BUSINESS_ID, Value: p.businessId},
		bson.E{Key: db_common.FLD_IS_DELETED, Value: false})

	log.Println("Find:: Got filter ", bfilter)
	singleResult := collection.FindOne(ctx, bfilter)
	if singleResult.Err() != nil {
		log.Println("Find:: Record not found ", singleResult.Err())
		return result, singleResult.Err()
	}
	singleResult.Decode(&result)
	if err != nil {
		log.Println("Error in decode", err)
		return result, err
	}

	// Remove fields from result
	result = db_common.AmendFldsForGet(result)

	log.Println("CustomerTypeDBDao::Find:: End Found a single document: \n", err)
	return result, nil
}

// Create - Create Collection
func (t *CustomerTypeMongoDBDao) Create(indata utils.Map) (utils.Map, error) {

	log.Println("CustomerType Save - Begin", indata)
	//Sales CustomerType
	collection, ctx, err := mongo_utils.GetMongoDbCollection(t.client, sales_common.DbCustomerTypes)
	if err != nil {
		log.Println("Error in insert ", err)
		return utils.Map{}, err
	}
	// Add Fields for Create
	indata = db_common.AmendFldsforCreate(indata)

	insertResult1, err := collection.InsertOne(ctx, indata)
	if err != nil {
		log.Println("Error in insert ", err)
		return utils.Map{}, err

	}
	log.Println("Inserted a single document: ", insertResult1.InsertedID)
	log.Println("Save - End", indata[sales_common.FLD_CUSTOMER_TYPE_ID])

	return t.Get(indata[sales_common.FLD_CUSTOMER_TYPE_ID].(string))
}

// Update - Update Collection
func (t *CustomerTypeMongoDBDao) Update(CustomerTypeId string, indata utils.Map) (utils.Map, error) {

	log.Println("Update - Begin")

	//Sales CustomerType
	collection, ctx, err := mongo_utils.GetMongoDbCollection(t.client, sales_common.DbCustomerTypes)
	if err != nil {
		return utils.Map{}, err
	}
	// Modify Fields for Update
	indata = db_common.AmendFldsforUpdate(indata)
	log.Printf("Update - Values %v", indata)

	filterCustomerType := bson.D{{Key: sales_common.FLD_CUSTOMER_TYPE_ID, Value: CustomerTypeId}}
	updateResult1, err := collection.UpdateOne(ctx, filterCustomerType, bson.D{{Key: "$set", Value: indata}})
	if err != nil {
		return utils.Map{}, err
	}
	log.Println("Update a single document: ", updateResult1.ModifiedCount)

	log.Println("Update - End")
	return t.Get(CustomerTypeId)
}

// Delete - Delete Collection
func (t *CustomerTypeMongoDBDao) Delete(CustomerTypeId string) (int64, error) {

	log.Println("CustomerTypeMongoDBDao::Delete - Begin ", CustomerTypeId)

	// Sales CustomerType
	collection, ctx, err := mongo_utils.GetMongoDbCollection(t.client, sales_common.DbCustomerTypes)
	if err != nil {
		return 0, err
	}
	optsCustomerType := options.Delete().SetCollation(&options.Collation{
		Locale:    db_common.LOCALE,
		Strength:  1,
		CaseLevel: false,
	})

	filterCustomerType := bson.D{{Key: sales_common.FLD_CUSTOMER_TYPE_ID, Value: CustomerTypeId}}
	resCustomerType, err := collection.DeleteOne(ctx, filterCustomerType, optsCustomerType)
	if err != nil {
		log.Println("Error in delete ", err)
		return 0, err
	}
	log.Printf("CustomerTypeMongoDBDao::Delete - End deleted %v documents\n", resCustomerType.DeletedCount)
	return resCustomerType.DeletedCount, nil
}
