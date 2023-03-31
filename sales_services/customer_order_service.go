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

type Customer_orderService interface {
	// List - List All records
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Get - Find By Code
	Get(customer_orderId string) (utils.Map, error)
	// Find - Find the item
	Find(filter string) (utils.Map, error)
	// Create - Create Service
	Create(indata utils.Map) (utils.Map, error)
	// Update - Update Service
	Update(bcustomer_orderId string, indata utils.Map) (utils.Map, error)
	// Delete - Delete Service
	Delete(customer_orderId string, delete_permanent bool) error

	EndService()
}

type customer_orderBaseService struct {
	db_utils.DatabaseService
	daoCustomer_order sales_repository.Customer_orderDao
	daoBusiness       platform_repository.BusinessDao
	child             Customer_orderService
	businessId        string
}

// NewCustomer_orderService - Construct Customer_order
func NewCustomer_orderService(props utils.Map) (Customer_orderService, error) {
	funcode := sales_common.GetServiceModuleCode() + "M" + "01"

	p := customer_orderBaseService{}
	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Customer_orderService ")
	// Verify whether the business id data passed
	businessId, err := utils.IsMemberExist(props, sales_common.FLD_BUSINESS_ID)
	if err != nil {
		return nil, err
	}

	// Assign the BusinessId
	p.businessId = businessId
	p.initializeService()

	_, err = p.daoBusiness.GetDetails(businessId)
	if err != nil {
		err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "Invalid business_id", ErrorDetail: "Given app_business_id is not exist"}
		return nil, err
	}

	p.child = &p

	return &p, err
}

// EndLoyaltyCardService - Close all the services
func (p *customer_orderBaseService) EndService() {
	log.Printf("EndService ")
	p.CloseDatabaseService()
}

func (p *customer_orderBaseService) initializeService() {
	log.Printf("customer_orderBaseService:: GetBusinessDao ")
	p.daoCustomer_order = sales_repository.NewCustomer_orderDao(p.GetClient(), p.businessId)
	p.daoBusiness = platform_repository.NewBusinessDao(p.GetClient())
}

// List - List All records
func (p *customer_orderBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("customer_orderBaseService::FindAll - Begin")

	listdata, err := p.daoCustomer_order.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("customer_orderBaseService::FindAll - End ")
	return listdata, nil
}

// Get - Find By Code
func (p *customer_orderBaseService) Get(customerorderId string) (utils.Map, error) {
	log.Printf("customer_orderBaseService::Get::  Begin %v", customerorderId)

	data, err := p.daoCustomer_order.Get(customerorderId)

	log.Println("customer_orderBaseService::Get:: End ", data, err)
	return data, err
}

func (p *customer_orderBaseService) Find(filter string) (utils.Map, error) {
	fmt.Println("customer_orderBaseService::FindByCode::  Begin ", filter)

	data, err := p.daoCustomer_order.Find(filter)
	log.Println("customer_orderBaseService::FindByCode:: End ", data, err)
	return data, err
}

// Create - Create Service
func (p *customer_orderBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("customer_orderBaseService::Create - Begin")
	var customerorderId string

	dataval, dataok := indata[sales_common.FLD_CUSTOMER_ORDER_ID]
	if dataok {
		customerorderId = strings.ToLower(dataval.(string))
	} else {
		customerorderId = utils.GenerateUniqueId("c_order")
		log.Println("Unique customer_order ID", customerorderId)
	}

	// Assign BusinessId
	indata[sales_common.FLD_BUSINESS_ID] = p.businessId
	indata[sales_common.FLD_CUSTOMER_ORDER_ID] = customerorderId

	data, err := p.daoCustomer_order.Create(indata)
	if err != nil {
		return utils.Map{}, err
	}

	log.Println("customer_orderBaseService::Create - End ")
	return data, nil
}

// Update - Update Service
func (p *customer_orderBaseService) Update(customerorderId string, indata utils.Map) (utils.Map, error) {

	log.Println("customer_orderService::Update - Begin")

	data, err := p.daoCustomer_order.Update(customerorderId, indata)

	log.Println("customer_orderService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *customer_orderBaseService) Delete(customerorderId string, delete_permanent bool) error {

	log.Println("customer_orderService::Delete - Begin", customerorderId)

	if delete_permanent {
		result, err := p.daoCustomer_order.Delete(customerorderId)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(customerorderId, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("customer_orderService::Delete - End")
	return nil
}
