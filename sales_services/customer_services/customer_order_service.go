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

type CustomerOrderService interface {
	// List - List All records
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Get - Find By Code
	Get(customerOrderId string) (utils.Map, error)
	// Find - Find the item
	Find(filter string) (utils.Map, error)
	// Create - Create Service
	Create(indata utils.Map) (utils.Map, error)
	// Update - Update Service
	Update(bcustomerOrderId string, indata utils.Map) (utils.Map, error)
	// Delete - Delete Service
	Delete(customerOrderId string, delete_permanent bool) error

	EndService()
}

type customerOrderBaseService struct {
	db_utils.DatabaseService
	daoCustomerOrder customer_repository.CustomerOrderDao
	daoBusiness      platform_repository.BusinessDao
	daoCustomer      sales_repository.CustomerDao

	child      CustomerOrderService
	businessId string
	customerId string
}

// NewCustomerOrderService - Construct CustomerOrder
func NewCustomerOrderService(props utils.Map) (CustomerOrderService, error) {
	funcode := sales_common.GetServiceModuleCode() + "M" + "01"

	p := customerOrderBaseService{}
	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("CustomerOrderService ")
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

	_, err = p.daoBusiness.GetDetails(businessId)
	if err != nil {
		err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "Invalid business_id", ErrorDetail: "Given app_business_id is not exist"}
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
func (p *customerOrderBaseService) EndService() {
	log.Printf("EndService ")
	p.CloseDatabaseService()
}

func (p *customerOrderBaseService) initializeService() {
	log.Printf("customerOrderBaseService:: GetBusinessDao ")
	p.daoCustomerOrder = customer_repository.NewCustomerOrderDao(p.GetClient(), p.businessId, p.customerId)
	p.daoBusiness = platform_repository.NewBusinessDao(p.GetClient())
}

// List - List All records
func (p *customerOrderBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("customerOrderBaseService::FindAll - Begin")

	listdata, err := p.daoCustomerOrder.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("customerOrderBaseService::FindAll - End ")
	return listdata, nil
}

// Get - Find By Code
func (p *customerOrderBaseService) Get(customerorderId string) (utils.Map, error) {
	log.Printf("customerOrderBaseService::Get::  Begin %v", customerorderId)

	data, err := p.daoCustomerOrder.Get(customerorderId)

	log.Println("customerOrderBaseService::Get:: End ", data, err)
	return data, err
}

func (p *customerOrderBaseService) Find(filter string) (utils.Map, error) {
	fmt.Println("customerOrderBaseService::FindByCode::  Begin ", filter)

	data, err := p.daoCustomerOrder.Find(filter)
	log.Println("customerOrderBaseService::FindByCode:: End ", data, err)
	return data, err
}

// Create - Create Service
func (p *customerOrderBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("customerOrderBaseService::Create - Begin")
	var customerorderId string

	dataval, dataok := indata[sales_common.FLD_CUSTOMER_ORDER_ID]
	if dataok {
		customerorderId = strings.ToLower(dataval.(string))
	} else {
		customerorderId = utils.GenerateUniqueId("c_order")
		log.Println("Unique customerOrder ID", customerorderId)
	}

	// Assign BusinessId
	indata[sales_common.FLD_BUSINESS_ID] = p.businessId
	indata[sales_common.FLD_CUSTOMER_ORDER_ID] = customerorderId

	data, err := p.daoCustomerOrder.Create(indata)
	if err != nil {
		return utils.Map{}, err
	}

	log.Println("customerOrderBaseService::Create - End ")
	return data, nil
}

// Update - Update Service
func (p *customerOrderBaseService) Update(customerorderId string, indata utils.Map) (utils.Map, error) {

	log.Println("customerOrderService::Update - Begin")

	data, err := p.daoCustomerOrder.Update(customerorderId, indata)

	log.Println("customerOrderService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *customerOrderBaseService) Delete(customerorderId string, delete_permanent bool) error {

	log.Println("customerOrderService::Delete - Begin", customerorderId)

	if delete_permanent {
		result, err := p.daoCustomerOrder.Delete(customerorderId)
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

	log.Printf("customerOrderService::Delete - End")
	return nil
}
