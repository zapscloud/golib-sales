package customer_mongodb_repository

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

// CustomerCartMongoDBDao - Cart DAO Repository
type CustomerCartMongoDBDao struct {
	client     utils.Map
	businessId string
	customerId string
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func (p *CustomerCartMongoDBDao) InitializeDao(client utils.Map, businessId string, customerId string) {
	log.Println("Initialize Cart Mongodb DAO")
	p.client = client
	p.businessId = businessId
	p.customerId = customerId
}

// List - List all Collections
func (t *CustomerCartMongoDBDao) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {
	var results []utils.Map

	log.Println("Begin - Find All Collection Dao", sales_common.DbCustomerCarts)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(t.client, sales_common.DbCustomerCarts)
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

	// Append customerId as filter if it available
	if len(t.customerId) > 0 {
		filterdoc = append(filterdoc, bson.E{Key: sales_common.FLD_CUSTOMER_ID, Value: t.customerId})
	}

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
	for _, value := range results {
		// log.Println("Item ", idx)
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

	// Append customerId as filter if it available
	if len(t.customerId) > 0 {
		basefilterdoc = append(basefilterdoc, bson.E{Key: sales_common.FLD_CUSTOMER_ID, Value: t.customerId})
	}
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
func (t *CustomerCartMongoDBDao) Get(cartId string) (utils.Map, error) {
	// Get a single document
	var result utils.Map

	log.Println("CustomerCartMongoDBDao::Get:: Begin ", cartId)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(t.client, sales_common.DbCustomerCarts)
	log.Println("Get:: Got Collection ")

	filter := bson.D{{Key: sales_common.FLD_CART_ID, Value: cartId}, {}}

	filter = append(filter,
		bson.E{Key: sales_common.FLD_BUSINESS_ID, Value: t.businessId},
		bson.E{Key: db_common.FLD_IS_DELETED, Value: false})

	// Append customerId as filter if it available
	if len(t.customerId) > 0 {
		filter = append(filter, bson.E{Key: sales_common.FLD_CUSTOMER_ID, Value: t.customerId})
	}

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

	log.Printf("Business CustomerCartMongoDBDao::Get:: End Found a single document\n")
	return result, nil
}

// Find - Find by Filter
func (p *CustomerCartMongoDBDao) Find(filter string) (utils.Map, error) {
	// Find a single document
	var result utils.Map

	log.Println("CartDBDao::Find:: Begin ", filter)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, sales_common.DbCustomerCarts)
	log.Println("Find:: Got Collection ", err)

	bfilter := bson.D{}
	err = bson.UnmarshalExtJSON([]byte(filter), true, &bfilter)
	if err != nil {
		fmt.Println("Error on filter Unmarshal", err)
	}
	bfilter = append(bfilter,
		bson.E{Key: sales_common.FLD_BUSINESS_ID, Value: p.businessId},
		bson.E{Key: db_common.FLD_IS_DELETED, Value: false})

	// Append customerId as filter if it available
	if len(p.customerId) > 0 {
		bfilter = append(bfilter, bson.E{Key: sales_common.FLD_CUSTOMER_ID, Value: p.customerId})
	}

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

	log.Println("CartDBDao::Find:: End Found a single document: \n", err)
	return result, nil
}

// Create - Create Collection
func (t *CustomerCartMongoDBDao) Create(indata utils.Map) (utils.Map, error) {

	log.Println("Cart Save - Begin", indata)
	//Sales Cart
	collection, ctx, err := mongo_utils.GetMongoDbCollection(t.client, sales_common.DbCustomerCarts)
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
	log.Println("Save - End", indata[sales_common.FLD_CART_ID])

	return t.Get(indata[sales_common.FLD_CART_ID].(string))
}

// Update - Update Collection
func (t *CustomerCartMongoDBDao) Update(cartId string, indata utils.Map) (utils.Map, error) {

	log.Println("Update - Begin")

	//Sales Cart
	collection, ctx, err := mongo_utils.GetMongoDbCollection(t.client, sales_common.DbCustomerCarts)
	if err != nil {
		return utils.Map{}, err
	}
	// Modify Fields for Update
	indata = db_common.AmendFldsforUpdate(indata)
	log.Printf("Update - Values %v", indata)

	filterCart := bson.D{{Key: sales_common.FLD_CART_ID, Value: cartId}}
	updateResult1, err := collection.UpdateOne(ctx, filterCart, bson.D{{Key: "$set", Value: indata}})
	if err != nil {
		return utils.Map{}, err
	}
	log.Println("Update a single document: ", updateResult1.ModifiedCount)

	log.Println("Update - End")
	return t.Get(cartId)
}

// Delete - Delete Collection
func (t *CustomerCartMongoDBDao) Delete(cartId string) (int64, error) {

	log.Println("CustomerCartMongoDBDao::Delete - Begin ", cartId)

	// Sales Cart
	collection, ctx, err := mongo_utils.GetMongoDbCollection(t.client, sales_common.DbCustomerCarts)
	if err != nil {
		return 0, err
	}
	optsCart := options.Delete().SetCollation(&options.Collation{
		Locale:    db_common.LOCALE,
		Strength:  1,
		CaseLevel: false,
	})

	filterCart := bson.D{{Key: sales_common.FLD_CART_ID, Value: cartId}}
	resCart, err := collection.DeleteOne(ctx, filterCart, optsCart)
	if err != nil {
		log.Println("Error in delete ", err)
		return 0, err
	}
	log.Printf("CustomerCartMongoDBDao::Delete - End deleted %v documents\n", resCart.DeletedCount)
	return resCart.DeletedCount, nil
}
