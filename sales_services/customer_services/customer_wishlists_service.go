package customer_services

import (
	"fmt"
	"log"
	"strings"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/db_utils"
	"github.com/zapscloud/golib-platform/platform_repository"
	"github.com/zapscloud/golib-sales/sales_common"
	"github.com/zapscloud/golib-sales/sales_repository"
	"github.com/zapscloud/golib-sales/sales_repository/customer_repository"

	"github.com/zapscloud/golib-utils/utils"
)

type CustomerWishlistService interface {
	// List - List All records
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Get - Find By Code
	Get(wishlistId string) (utils.Map, error)
	// Find - Find the item
	Find(filter string) (utils.Map, error)
	// Create - Create Service
	Create(indata utils.Map) (utils.Map, error)
	// Update - Update Service
	Update(wishlistId string, indata utils.Map) (utils.Map, error)
	// Delete - Delete Service
	Delete(wishlistId string, delete_permanent bool) error

	EndService()
}

type customerwishlistBaseService struct {
	db_utils.DatabaseService
	daoCustomerWishlist customer_repository.CustomerWishlistDao
	daoBusiness         platform_repository.BusinessDao
	daoCustomer         sales_repository.CustomerDao

	child      CustomerWishlistService
	businessId string
	customerId string
}

// NewCustomerWishlistService - Construct CustomerWishlist
func NewCustomerWishlistService(props utils.Map) (CustomerWishlistService, error) {
	funcode := sales_common.GetServiceModuleCode() + "M" + "01"

	p := customerwishlistBaseService{}
	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("CustomerWishlistService ")
	// Verify whether the business id data passed
	businessId, err := utils.IsMemberExist(props, sales_common.FLD_BUSINESS_ID)
	if err != nil {
		return nil, err
	}
	// Verify whether the User id data passed, this is optional parameter
	customerId, _ := utils.IsMemberExist(props, sales_common.FLD_CUSTOMER_ID)
	// if err != nil {
	// 	return nil, err
	// }

	// Assign the BusinessId
	p.businessId = businessId
	p.customerId = customerId
	p.initializeService()

	_, err = p.daoBusiness.Get(businessId)
	if err != nil {
		err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "Invalid business_id", ErrorDetail: "Given app_business_id is not exist"}
		return nil, err
	}

	// Verify the Customer Exist
	if len(customerId) > 0 {
		_, err = p.daoCustomer.Get(customerId)
		if err != nil {
			err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "Invalid CustomerId", ErrorDetail: "Given CustomerId is not exist"}
			return nil, err
		}
	}

	p.child = &p

	return &p, err
}

// EndLoyaltyCardService - Close all the services
func (p *customerwishlistBaseService) EndService() {
	log.Printf("EndService ")
	p.CloseDatabaseService()
}

func (p *customerwishlistBaseService) initializeService() {
	log.Printf("CustomerWishlistService:: GetBusinessDao ")
	p.daoCustomerWishlist = customer_repository.NewCustomerWishlistDao(p.GetClient(), p.businessId, p.customerId)
	p.daoBusiness = platform_repository.NewBusinessDao(p.GetClient())
	p.daoCustomer = sales_repository.NewCustomerDao(p.GetClient(), p.businessId)
}

// List - List All records
func (p *customerwishlistBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("customerwishlistBaseService::FindAll - Begin")

	listdata, err := p.daoCustomerWishlist.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("customerwishlistBaseService::FindAll - End ")
	return listdata, nil
}

// Get - Find By Code
func (p *customerwishlistBaseService) Get(wishlistId string) (utils.Map, error) {
	log.Printf("customerwishlistBaseService::Get::  Begin %v", wishlistId)

	data, err := p.daoCustomerWishlist.Get(wishlistId)

	log.Println("customerwishlistBaseService::Get:: End ", err)
	return data, err
}

func (p *customerwishlistBaseService) Find(filter string) (utils.Map, error) {
	fmt.Println("customerwishlistBaseService::FindByCode::  Begin ", filter)

	data, err := p.daoCustomerWishlist.Find(filter)
	log.Println("customerwishlistBaseService::FindByCode:: End ", err)
	return data, err
}

// Create - Create Service
func (p *customerwishlistBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("CustomerWishlistService::Create - Begin")
	var wishlistId string

	dataval, dataok := indata[sales_common.FLD_WISHLIST_ID]
	if dataok {
		wishlistId = strings.ToLower(dataval.(string))
	} else {
		wishlistId = utils.GenerateUniqueId("wish")
		log.Println("Unique CustomerWishlist ID", wishlistId)
	}

	// Assign BusinessId
	indata[sales_common.FLD_BUSINESS_ID] = p.businessId
	indata[sales_common.FLD_CUSTOMER_ID] = p.customerId
	indata[sales_common.FLD_WISHLIST_ID] = wishlistId

	data, err := p.daoCustomerWishlist.Create(indata)
	if err != nil {
		return utils.Map{}, err
	}

	log.Println("CustomerWishlistService::Create - End ")
	return data, nil
}

// Update - Update Service
func (p *customerwishlistBaseService) Update(wishlistId string, indata utils.Map) (utils.Map, error) {

	log.Println("CustomerWishlistService::Update - Begin")

	// Delete Key values
	delete(indata, sales_common.FLD_BUSINESS_ID)
	delete(indata, sales_common.FLD_CUSTOMER_ID)
	delete(indata, sales_common.FLD_WISHLIST_ID)

	data, err := p.daoCustomerWishlist.Update(wishlistId, indata)

	log.Println("CustomerWishlistService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *customerwishlistBaseService) Delete(wishlistId string, delete_permanent bool) error {

	log.Println("CustomerWishlistService::Delete - Begin", wishlistId)

	if delete_permanent {
		result, err := p.daoCustomerWishlist.Delete(wishlistId)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(wishlistId, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("CustomerWishlistService::Delete - End")
	return nil
}
