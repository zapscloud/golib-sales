package sales_services

import (
	"fmt"
	"log"
	"strings"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/db_utils"
	"github.com/zapscloud/golib-platform/platform_repository"
	"github.com/zapscloud/golib-sales/sales_common"
	"github.com/zapscloud/golib-sales/sales_repository"
	"github.com/zapscloud/golib-utils/utils"
)

type CartService interface {
	// List - List All records
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Get - Find By Code
	Get(cartId string) (utils.Map, error)
	// Find - Find the item
	Find(filter string) (utils.Map, error)
	// Create - Create Service
	Create(indata utils.Map) (utils.Map, error)
	// Update - Update Service
	Update(cartId string, indata utils.Map) (utils.Map, error)
	// Delete - Delete Service
	Delete(cartId string, delete_permanent bool) error

	EndService()
}

type cartBaseService struct {
	db_utils.DatabaseService
	daoCart     sales_repository.CartDao
	daoBusiness platform_repository.BusinessDao
	daoCustomer sales_repository.CustomerDao

	child      CartService
	businessId string
	customerId string
}

// NewCartService - Construct Cart
func NewCartService(props utils.Map) (CartService, error) {
	funcode := sales_common.GetServiceModuleCode() + "M" + "01"

	p := cartBaseService{}
	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("CartService ")
	// Verify whether the business id data passed
	businessId, err := utils.IsMemberExist(props, sales_common.FLD_BUSINESS_ID)
	if err != nil {
		return nil, err
	}

	// Verify whether the User id data passed
	customerId, err := utils.IsMemberExist(props, sales_common.FLD_CUSTOMER_ID)
	if err != nil {
		return nil, err
	}

	// Assign the BusinessId
	p.businessId = businessId
	p.customerId = customerId
	p.initializeService()

	// Verify the Business Exists
	_, err = p.daoBusiness.GetDetails(businessId)
	if err != nil {
		err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "Invalid BusinessId", ErrorDetail: "Given BusinessId is not exist"}
		return nil, err
	}

	// Verify the Customer Exist
	_, err = p.daoCustomer.Get(customerId)
	if err != nil {
		err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "Invalid CustomerId", ErrorDetail: "Given CustomerId is not exist"}
		return nil, err
	}

	p.child = &p

	return &p, err
}

// EndLoyaltyCardService - Close all the services
func (p *cartBaseService) EndService() {
	log.Printf("EndService ")
	p.CloseDatabaseService()
}

func (p *cartBaseService) initializeService() {
	log.Printf("CartService:: GetBusinessDao ")
	p.daoCart = sales_repository.NewCartDao(p.GetClient(), p.businessId, p.customerId)
	p.daoBusiness = platform_repository.NewBusinessDao(p.GetClient())
	p.daoCustomer = sales_repository.NewCustomerDao(p.GetClient(), p.businessId)
}

// List - List All records
func (p *cartBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("cartBaseService::FindAll - Begin")

	listdata, err := p.daoCart.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("cartBaseService::FindAll - End ")
	return listdata, nil
}

// Get - Find By Code
func (p *cartBaseService) Get(cartId string) (utils.Map, error) {
	log.Printf("cartBaseService::Get::  Begin %v", cartId)

	data, err := p.daoCart.Get(cartId)

	log.Println("cartBaseService::Get:: End ", data, err)
	return data, err
}

func (p *cartBaseService) Find(filter string) (utils.Map, error) {
	fmt.Println("CartService::FindByCode::  Begin ", filter)

	data, err := p.daoCart.Find(filter)
	log.Println("CartService::FindByCode:: End ", data, err)
	return data, err
}

// Create - Create Service
func (p *cartBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("CartService::Create - Begin")
	var cartId string

	dataval, dataok := indata[sales_common.FLD_CART_ID]
	if dataok {
		cartId = strings.ToLower(dataval.(string))
	} else {
		cartId = utils.GenerateUniqueId("crt")
		log.Println("Unique Cart ID", cartId)
	}

	// Assign BusinessId
	indata[sales_common.FLD_BUSINESS_ID] = p.businessId
	indata[sales_common.FLD_CUSTOMER_ID] = p.customerId
	indata[sales_common.FLD_CART_ID] = cartId

	data, err := p.daoCart.Create(indata)
	if err != nil {
		return utils.Map{}, err
	}

	log.Println("CartService::Create - End ")
	return data, nil
}

// Update - Update Service
func (p *cartBaseService) Update(cartId string, indata utils.Map) (utils.Map, error) {

	log.Println("CartService::Update - Begin")

	// Delete Key values
	delete(indata, sales_common.FLD_BUSINESS_ID)
	delete(indata, sales_common.FLD_CUSTOMER_ID)
	delete(indata, sales_common.FLD_CART_ID)

	data, err := p.daoCart.Update(cartId, indata)

	log.Println("CartService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *cartBaseService) Delete(cartId string, delete_permanent bool) error {

	log.Println("CartService::Delete - Begin", cartId)

	if delete_permanent {
		result, err := p.daoCart.Delete(cartId)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(cartId, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("CartService::Delete - End")
	return nil
}
