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

type CustomerTypeService interface {
	// List - List All records
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Get - Find By Code
	Get(TypeId string) (utils.Map, error)
	// Find - Find the item
	Find(filter string) (utils.Map, error)
	// Create - Create Service
	Create(indata utils.Map) (utils.Map, error)
	// Update - Update Service
	Update(TypeId string, indata utils.Map) (utils.Map, error)
	// Delete - Delete Service
	Delete(TypeId string, delete_permanent bool) error

	EndService()
}

type customerTypeBaseService struct {
	db_utils.DatabaseService
	daoCustomerType customer_repository.CustomerTypeDao
	daoBusiness     platform_repository.BusinessDao
	daoCustomer     sales_repository.CustomerDao

	child      CustomerTypeService
	businessId string
	customerId string
}

// NewCustomerTypeService - Construct CustomerType
func NewCustomerTypeService(props utils.Map) (CustomerTypeService, error) {
	funcode := sales_common.GetServiceModuleCode() + "M" + "01"

	p := customerTypeBaseService{}
	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("CustomerTypeService ")
	// Verify whether the business id data passed
	businessId, err := utils.GetMemberDataStr(props, sales_common.FLD_BUSINESS_ID)
	if err != nil {
		return p.errorReturn(err)
	}

	// Verify whether the User id data passed
	customerId, err := utils.GetMemberDataStr(props, sales_common.FLD_CUSTOMER_ID)
	if err != nil {
		return p.errorReturn(err)
	}

	// Assign the BusinessId
	p.businessId = businessId
	p.customerId = customerId
	p.initializeService()

	// Verify the Business Exists
	_, err = p.daoBusiness.Get(businessId)
	if err != nil {
		err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "Invalid BusinessId", ErrorDetail: "Given BusinessId is not exist"}
		return p.errorReturn(err)
	}

	// Verify the Customer Exist
	_, err = p.daoCustomer.Get(customerId)
	if err != nil {
		err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "Invalid CustomerId", ErrorDetail: "Given CustomerId is not exist"}
		return p.errorReturn(err)
	}

	p.child = &p

	return &p, err
}

// EndLoyaltyCardService - Close all the services
func (p *customerTypeBaseService) EndService() {
	log.Printf("EndService ")
	p.CloseDatabaseService()
}

func (p *customerTypeBaseService) initializeService() {
	log.Printf("CustomerTypeService:: GetBusinessDao ")
	p.daoCustomerType = customer_repository.NewCustomerTypeDao(p.GetClient(), p.businessId, p.customerId)
	p.daoBusiness = platform_repository.NewBusinessDao(p.GetClient())
	p.daoCustomer = sales_repository.NewCustomerDao(p.GetClient(), p.businessId)
}

// List - List All records
func (p *customerTypeBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("customerTypeBaseService::FindAll - Begin")

	listdata, err := p.daoCustomerType.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("customerTypeBaseService::FindAll - End ")
	return listdata, nil
}

// Get - Find By Code
func (p *customerTypeBaseService) Get(TypeId string) (utils.Map, error) {
	log.Printf("customerTypeBaseService::Get::  Begin %v", TypeId)

	data, err := p.daoCustomerType.Get(TypeId)

	log.Println("customerTypeBaseService::Get:: End ", err)
	return data, err
}

func (p *customerTypeBaseService) Find(filter string) (utils.Map, error) {
	fmt.Println("CustomerTypeService::FindByCode::  Begin ", filter)

	data, err := p.daoCustomerType.Find(filter)
	log.Println("CustomerTypeService::FindByCode:: End ", err)
	return data, err
}

// Create - Create Service
func (p *customerTypeBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("CustomerTypeService::Create - Begin")
	var TypeId string

	dataval, dataok := indata[sales_common.FLD_CUSTOMER_TYPE_ID]
	if dataok {
		TypeId = strings.ToLower(dataval.(string))
	} else {
		TypeId = utils.GenerateUniqueId("crt")
		log.Println("Unique CustomerType ID", TypeId)
	}

	// Assign BusinessId
	indata[sales_common.FLD_BUSINESS_ID] = p.businessId
	indata[sales_common.FLD_CUSTOMER_TYPE_ID] = TypeId

	data, err := p.daoCustomerType.Create(indata)
	if err != nil {
		return utils.Map{}, err
	}

	log.Println("CustomerTypeService::Create - End ")
	return data, nil
}

// Update - Update Service
func (p *customerTypeBaseService) Update(TypeId string, indata utils.Map) (utils.Map, error) {

	log.Println("CustomerTypeService::Update - Begin")

	// Delete Key values
	delete(indata, sales_common.FLD_BUSINESS_ID)
	delete(indata, sales_common.FLD_CUSTOMER_ID)
	delete(indata, sales_common.FLD_CUSTOMER_TYPE_ID)

	data, err := p.daoCustomerType.Update(TypeId, indata)

	log.Println("CustomerTypeService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *customerTypeBaseService) Delete(TypeId string, delete_permanent bool) error {

	log.Println("CustomerTypeService::Delete - Begin", TypeId)

	if delete_permanent {
		result, err := p.daoCustomerType.Delete(TypeId)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(TypeId, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("CustomerTypeService::Delete - End")
	return nil
}

func (p *customerTypeBaseService) errorReturn(err error) (CustomerTypeService, error) {
	// Close the Database Connection
	p.CloseDatabaseService()
	return nil, err
}
