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

// DealerMongoDBDao - Dealer DAO Repository
type DealerMongoDBDao struct {
	client     utils.Map
	businessId string
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func (p *DealerMongoDBDao) InitializeDao(client utils.Map, businessId string) {
	log.Println("Initialize Dealer Mongodb DAO")
	p.client = client
	p.businessId = businessId
}

// List - List all Collections
func (t *DealerMongoDBDao) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {
	var results []utils.Map

	log.Println("Begin - Find All Collection Dao", sales_common.DbDealers)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(t.client, sales_common.DbDealers)
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

	log.Println("End - Find All Collection Dao", results)

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
func (p *DealerMongoDBDao) Get(dealerId string) (utils.Map, error) {
	// Get a single document
	var result utils.Map

	log.Println("DealerMongoDBDao::Get:: Begin ", dealerId)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, sales_common.DbDealers)
	log.Println("Get:: Got Collection ")

	filter := bson.D{{Key: sales_common.FLD_DEALER_ID, Value: dealerId}, {}}

	filter = append(filter,
		bson.E{Key: sales_common.FLD_BUSINESS_ID, Value: p.businessId},
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

	log.Printf("Business DealerMongoDBDao::Get:: End Found a single document: %+v\n", result)
	return result, nil
}

// Find - Find by Filter
func (p *DealerMongoDBDao) Find(filter string) (utils.Map, error) {
	// Find a single document
	var result utils.Map

	log.Println("DealerDBDao::Find:: Begin ", filter)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, sales_common.DbDealers)
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

	log.Println("DealerDBDao::Find:: End Found a single document: \n", err)
	return result, nil
}

// Create - Create Collection
func (t *DealerMongoDBDao) Create(indata utils.Map) (utils.Map, error) {

	log.Println("Dealer Save - Begin", indata)
	//Business_dealer
	collection, ctx, err := mongo_utils.GetMongoDbCollection(t.client, sales_common.DbDealers)
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
	log.Println("Save - End", indata[sales_common.FLD_DEALER_ID])

	return t.Get(indata[sales_common.FLD_DEALER_ID].(string))
}

// Update - Update Collection
func (t *DealerMongoDBDao) Update(dealerId string, indata utils.Map) (utils.Map, error) {

	log.Println("Update - Begin")

	//dealer
	collection, ctx, err := mongo_utils.GetMongoDbCollection(t.client, sales_common.DbDealers)
	if err != nil {
		return utils.Map{}, err
	}
	// Modify Fields for Update
	indata = db_common.AmendFldsforUpdate(indata)
	log.Printf("Update - Values %v", indata)

	filterDealer := bson.D{{Key: sales_common.FLD_DEALER_ID, Value: dealerId}}
	updateResult1, err := collection.UpdateOne(ctx, filterDealer, bson.D{{Key: "$set", Value: indata}})
	if err != nil {
		return utils.Map{}, err
	}
	log.Println("Update a single document: ", updateResult1.ModifiedCount)

	log.Println("Update - End")
	return t.Get(dealerId)
}

// Delete - Delete Collection
func (t *DealerMongoDBDao) Delete(dealerId string) (int64, error) {

	log.Println("DealerMongoDBDao::Delete - Begin ", dealerId)

	//BusinessDealer
	collection, ctx, err := mongo_utils.GetMongoDbCollection(t.client, sales_common.DbDealers)
	if err != nil {
		return 0, err
	}
	optsDealer := options.Delete().SetCollation(&options.Collation{
		Locale:    db_common.LOCALE,
		Strength:  1,
		CaseLevel: false,
	})

	filterDealer := bson.D{{Key: sales_common.FLD_DEALER_ID, Value: dealerId}}
	resDealer, err := collection.DeleteOne(ctx, filterDealer, optsDealer)
	if err != nil {
		log.Println("Error in delete ", err)
		return 0, err
	}
	log.Printf("DealerMongoDBDao::Delete - End deleted %v documents\n", resDealer.DeletedCount)
	return resDealer.DeletedCount, nil
}
