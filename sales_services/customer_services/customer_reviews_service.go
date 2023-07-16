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

type CustomerReviewService interface {
	// List - List All records
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Get - Find By Code
	Get(reviewId string) (utils.Map, error)
	// Find - Find the item
	Find(filter string) (utils.Map, error)
	// Create - Create Service
	Create(indata utils.Map) (utils.Map, error)
	// Update - Update Service
	Update(reviewId string, indata utils.Map) (utils.Map, error)
	// Delete - Delete Service
	Delete(reviewId string, delete_permanent bool) error

	EndService()
}

type customerreviewBaseService struct {
	db_utils.DatabaseService
	daoCustomerReview customer_repository.CustomerReviewDao
	daoBusiness       platform_repository.BusinessDao
	daoCustomer       sales_repository.CustomerDao

	child      CustomerReviewService
	businessId string
	customerId string
}

// NewCustomerReviewService - Construct CustomerReview
func NewCustomerReviewService(props utils.Map) (CustomerReviewService, error) {
	funcode := sales_common.GetServiceModuleCode() + "M" + "01"

	p := customerreviewBaseService{}
	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("CustomerReviewService ")
	// Verify whether the business id data passed
	businessId, err := utils.GetMemberDataStr(props, sales_common.FLD_BUSINESS_ID)
	if err != nil {
		return p.errorReturn(err)
	}
	// Verify whether the User id data passed, this is optional parameter
	customerId, _ := utils.GetMemberDataStr(props, sales_common.FLD_CUSTOMER_ID)
	// if err != nil {
	// 	return p.errorReturn(err)
	// }

	// Assign the BusinessId
	p.businessId = businessId
	p.customerId = customerId
	p.initializeService()

	_, err = p.daoBusiness.Get(businessId)
	if err != nil {
		err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "Invalid business_id", ErrorDetail: "Given business_id is not exist"}
		return p.errorReturn(err)
	}

	// Verify the Customer Exist
	if len(customerId) > 0 {
		_, err = p.daoCustomer.Get(customerId)
		if err != nil {
			err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "Invalid CustomerId", ErrorDetail: "Given CustomerId is not exist"}
			return p.errorReturn(err)
		}
	}

	p.child = &p

	return &p, err
}

// EndLoyaltyCardService - Close all the services
func (p *customerreviewBaseService) EndService() {
	log.Printf("EndService ")
	p.CloseDatabaseService()
}

func (p *customerreviewBaseService) initializeService() {
	log.Printf("CustomerReviewService:: GetBusinessDao ")
	p.daoCustomerReview = customer_repository.NewCustomerReviewDao(p.GetClient(), p.businessId, p.customerId)
	p.daoBusiness = platform_repository.NewBusinessDao(p.GetClient())
	p.daoCustomer = sales_repository.NewCustomerDao(p.GetClient(), p.businessId)
}

// List - List All records
func (p *customerreviewBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("customerreviewBaseService::FindAll - Begin")

	listdata, err := p.daoCustomerReview.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("customerreviewBaseService::FindAll - End ")
	return listdata, nil
}

// Get - Find By Code
func (p *customerreviewBaseService) Get(reviewId string) (utils.Map, error) {
	log.Printf("customerreviewBaseService::Get::  Begin %v", reviewId)

	data, err := p.daoCustomerReview.Get(reviewId)

	log.Println("customerreviewBaseService::Get:: End ", err)
	return data, err
}

func (p *customerreviewBaseService) Find(filter string) (utils.Map, error) {
	fmt.Println("customerreviewBaseService::FindByCode::  Begin ", filter)

	data, err := p.daoCustomerReview.Find(filter)
	log.Println("customerreviewBaseService::FindByCode:: End ", err)
	return data, err
}

// Create - Create Service
func (p *customerreviewBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("CustomerReviewService::Create - Begin")
	var reviewId string

	dataval, dataok := indata[sales_common.FLD_REVIEW_ID]
	if dataok {
		reviewId = strings.ToLower(dataval.(string))
	} else {
		reviewId = utils.GenerateUniqueId("reviw")
		log.Println("Unique CustomerReview ID", reviewId)
	}

	// Assign BusinessId
	indata[sales_common.FLD_BUSINESS_ID] = p.businessId
	indata[sales_common.FLD_CUSTOMER_ID] = p.customerId
	indata[sales_common.FLD_REVIEW_ID] = reviewId

	data, err := p.daoCustomerReview.Create(indata)
	if err != nil {
		return utils.Map{}, err
	}

	log.Println("CustomerReviewService::Create - End ")
	return data, nil
}

// Update - Update Service
func (p *customerreviewBaseService) Update(reviewId string, indata utils.Map) (utils.Map, error) {

	log.Println("CustomerReviewService::Update - Begin")

	// Delete Key values
	delete(indata, sales_common.FLD_BUSINESS_ID)
	delete(indata, sales_common.FLD_CUSTOMER_ID)
	delete(indata, sales_common.FLD_REVIEW_ID)

	data, err := p.daoCustomerReview.Update(reviewId, indata)

	log.Println("CustomerReviewService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *customerreviewBaseService) Delete(reviewId string, delete_permanent bool) error {

	log.Println("CustomerReviewService::Delete - Begin", reviewId)

	if delete_permanent {
		result, err := p.daoCustomerReview.Delete(reviewId)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(reviewId, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("CustomerReviewService::Delete - End")
	return nil
}

func (p *customerreviewBaseService) errorReturn(err error) (CustomerReviewService, error) {
	// Close the Database Connection
	p.CloseDatabaseService()
	return nil, err
}
