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

// CouponMongoDBDao - Coupon DAO Repository
type CouponMongoDBDao struct {
	client     utils.Map
	businessId string
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func (p *CouponMongoDBDao) InitializeDao(client utils.Map, businessId string) {
	log.Println("Initialize Coupon Mongodb DAO")
	p.client = client
	p.businessId = businessId
}

// List - List all Collections
func (t *CouponMongoDBDao) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {
	var results []utils.Map

	log.Println("Begin - Find All Collection Dao", sales_common.DbCoupons)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(t.client, sales_common.DbCoupons)
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
func (p *CouponMongoDBDao) Get(CouponId string) (utils.Map, error) {
	// Get a single document
	var result utils.Map

	log.Println("CouponMongoDBDao::Get:: Begin ", CouponId)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, sales_common.DbCoupons)
	log.Println("Get:: Got Collection ")

	filter := bson.D{{Key: sales_common.FLD_COUPON_ID, Value: CouponId}, {}}

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

	log.Printf("Business CouponMongoDBDao::Get:: End Found a single document\n")
	return result, nil
}

// Find - Find by Filter
func (p *CouponMongoDBDao) Find(filter string) (utils.Map, error) {
	// Find a single document
	var result utils.Map

	log.Println("CouponDBDao::Find:: Begin ", filter)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, sales_common.DbCoupons)
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

	log.Println("CouponDBDao::Find:: End Found a single document: \n", err)
	return result, nil
}

// Create - Create Collection
func (t *CouponMongoDBDao) Create(indata utils.Map) (utils.Map, error) {

	log.Println("Coupon Save - Begin", indata)
	//Business_Coupon
	collection, ctx, err := mongo_utils.GetMongoDbCollection(t.client, sales_common.DbCoupons)
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
	log.Println("Save - End", indata[sales_common.FLD_COUPON_ID])

	return t.Get(indata[sales_common.FLD_COUPON_ID].(string))
}

// Update - Update Collection
func (t *CouponMongoDBDao) Update(CouponId string, indata utils.Map) (utils.Map, error) {

	log.Println("Update - Begin")

	//Coupon
	collection, ctx, err := mongo_utils.GetMongoDbCollection(t.client, sales_common.DbCoupons)
	if err != nil {
		return utils.Map{}, err
	}
	// Modify Fields for Update
	indata = db_common.AmendFldsforUpdate(indata)
	log.Printf("Update - Values %v", indata)

	filterCoupon := bson.D{{Key: sales_common.FLD_COUPON_ID, Value: CouponId}}
	updateResult1, err := collection.UpdateOne(ctx, filterCoupon, bson.D{{Key: "$set", Value: indata}})
	if err != nil {
		return utils.Map{}, err
	}
	log.Println("Update a single document: ", updateResult1.ModifiedCount)

	log.Println("Update - End")
	return t.Get(CouponId)
}

// Delete - Delete Collection
func (t *CouponMongoDBDao) Delete(CouponId string) (int64, error) {

	log.Println("CouponMongoDBDao::Delete - Begin ", CouponId)

	//BusinessCoupon
	collection, ctx, err := mongo_utils.GetMongoDbCollection(t.client, sales_common.DbCoupons)
	if err != nil {
		return 0, err
	}
	optsCoupon := options.Delete().SetCollation(&options.Collation{
		Locale:    db_common.LOCALE,
		Strength:  1,
		CaseLevel: false,
	})

	filterCoupon := bson.D{{Key: sales_common.FLD_COUPON_ID, Value: CouponId}}
	resCoupon, err := collection.DeleteOne(ctx, filterCoupon, optsCoupon)
	if err != nil {
		log.Println("Error in delete ", err)
		return 0, err
	}
	log.Printf("CouponMongoDBDao::Delete - End deleted %v documents\n", resCoupon.DeletedCount)
	return resCoupon.DeletedCount, nil
}
