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

type ReviewService interface {
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

type reviewBaseService struct {
	db_utils.DatabaseService
	daoReview   sales_repository.ReviewDao
	daoBusiness platform_repository.BusinessDao
	child       ReviewService
	businessId  string
}

// NewReviewService - Construct Review
func NewReviewService(props utils.Map) (ReviewService, error) {
	funcode := sales_common.GetServiceModuleCode() + "M" + "01"

	p := reviewBaseService{}
	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ReviewService ")
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
func (p *reviewBaseService) EndService() {
	log.Printf("EndService ")
	p.CloseDatabaseService()
}

func (p *reviewBaseService) initializeService() {
	log.Printf("ReviewService:: GetBusinessDao ")
	p.daoReview = sales_repository.NewReviewDao(p.GetClient(), p.businessId)
	p.daoBusiness = platform_repository.NewBusinessDao(p.GetClient())
}

// List - List All records
func (p *reviewBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("reviewBaseService::FindAll - Begin")

	listdata, err := p.daoReview.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("reviewBaseService::FindAll - End ")
	return listdata, nil
}

// Get - Find By Code
func (p *reviewBaseService) Get(reviewId string) (utils.Map, error) {
	log.Printf("reviewBaseService::Get::  Begin %v", reviewId)

	data, err := p.daoReview.Get(reviewId)

	log.Println("reviewBaseService::Get:: End ", data, err)
	return data, err
}

func (p *reviewBaseService) Find(filter string) (utils.Map, error) {
	fmt.Println("reviewBaseService::FindByCode::  Begin ", filter)

	data, err := p.daoReview.Find(filter)
	log.Println("reviewBaseService::FindByCode:: End ", data, err)
	return data, err
}

// Create - Create Service
func (p *reviewBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("ReviewService::Create - Begin")
	var reviewId string

	dataval, dataok := indata[sales_common.FLD_REVIEW_ID]
	if dataok {
		reviewId = strings.ToLower(dataval.(string))
	} else {
		reviewId = utils.GenerateUniqueId("rev")
		log.Println("Unique Review ID", reviewId)
	}

	// Assign BusinessId
	indata[sales_common.FLD_BUSINESS_ID] = p.businessId
	indata[sales_common.FLD_REVIEW_ID] = reviewId

	data, err := p.daoReview.Create(indata)
	if err != nil {
		return utils.Map{}, err
	}

	log.Println("ReviewService::Create - End ")
	return data, nil
}

// Update - Update Service
func (p *reviewBaseService) Update(reviewId string, indata utils.Map) (utils.Map, error) {

	log.Println("ReviewService::Update - Begin")

	data, err := p.daoReview.Update(reviewId, indata)

	log.Println("ReviewService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *reviewBaseService) Delete(reviewId string, delete_permanent bool) error {

	log.Println("ReviewService::Delete - Begin", reviewId)

	if delete_permanent {
		result, err := p.daoReview.Delete(reviewId)
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

	log.Printf("ReviewService::Delete - End")
	return nil
}
