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

type AboutUsService interface {
	// List - List All records
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Get - Find By Code
	Get(aboutusId string) (utils.Map, error)
	// Find - Find the item
	Find(filter string) (utils.Map, error)
	// Create - Create Service
	Create(indata utils.Map) (utils.Map, error)
	// Update - Update Service
	Update(aboutusId string, indata utils.Map) (utils.Map, error)
	// Delete - Delete Service
	Delete(aboutusId string, delete_permanent bool) error

	EndService()
}

type aboutusBaseService struct {
	db_utils.DatabaseService
	daoAboutUs  sales_repository.AboutUsDao
	daoBusiness platform_repository.BusinessDao
	child       AboutUsService
	businessId  string
}

// NewAboutUsService - Construct AboutUs
func NewAboutUsService(props utils.Map) (AboutUsService, error) {
	funcode := sales_common.GetServiceModuleCode() + "M" + "01"

	p := aboutusBaseService{}
	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("AboutUsService ")
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
func (p *aboutusBaseService) EndService() {
	log.Printf("EndService ")
	p.CloseDatabaseService()
}

func (p *aboutusBaseService) initializeService() {
	log.Printf("AboutUsService:: GetBusinessDao ")
	p.daoAboutUs = sales_repository.NewAboutUsDao(p.GetClient(), p.businessId)
	p.daoBusiness = platform_repository.NewBusinessDao(p.GetClient())
}

// List - List All records
func (p *aboutusBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("aboutusBaseService::FindAll - Begin")

	listdata, err := p.daoAboutUs.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("aboutusBaseService::FindAll - End ")
	return listdata, nil
}

// Get - Find By Code
func (p *aboutusBaseService) Get(aboutusId string) (utils.Map, error) {
	log.Printf("aboutusBaseService::Get::  Begin %v", aboutusId)

	data, err := p.daoAboutUs.Get(aboutusId)

	log.Println("aboutusBaseService::Get:: End ", data, err)
	return data, err
}

func (p *aboutusBaseService) Find(filter string) (utils.Map, error) {
	fmt.Println("brandService::FindByCode::  Begin ", filter)

	data, err := p.daoAboutUs.Find(filter)
	log.Println("brandService::FindByCode:: End ", data, err)
	return data, err
}

// Create - Create Service
func (p *aboutusBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("AboutUsService::Create - Begin")
	var aboutusId string

	dataval, dataok := indata[sales_common.FLD_BRAND_ID]
	if dataok {
		aboutusId = strings.ToLower(dataval.(string))
	} else {
		aboutusId = utils.GenerateUniqueId("abtus")
		log.Println("Unique AboutUs ID", aboutusId)
	}

	// Assign BusinessId
	indata[sales_common.FLD_BUSINESS_ID] = p.businessId
	indata[sales_common.FLD_BRAND_ID] = aboutusId

	data, err := p.daoAboutUs.Create(indata)
	if err != nil {
		return utils.Map{}, err
	}

	log.Println("AboutUsService::Create - End ")
	return data, nil
}

// Update - Update Service
func (p *aboutusBaseService) Update(aboutusId string, indata utils.Map) (utils.Map, error) {

	log.Println("AboutUsService::Update - Begin")

	data, err := p.daoAboutUs.Update(aboutusId, indata)

	log.Println("AboutUsService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *aboutusBaseService) Delete(aboutusId string, delete_permanent bool) error {

	log.Println("AboutUsService::Delete - Begin", aboutusId)

	if delete_permanent {
		result, err := p.daoAboutUs.Delete(aboutusId)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(aboutusId, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("AboutUsService::Delete - End")
	return nil
}
