package sales_services

import (
	"fmt"
	"log"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/db_utils"
	"github.com/zapscloud/golib-platform/platform_repository"
	"github.com/zapscloud/golib-platform/platform_services"
	"github.com/zapscloud/golib-sales/sales_common"
	"github.com/zapscloud/golib-sales/sales_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// BrandService - Brand Service structure
type CompareProdsService interface {
	// List - List All records
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Get - Find By Code
	Get(materialTypeId string) (utils.Map, error)
	// Find - Find the item
	Find(filter string) (utils.Map, error)
	// Create - Create Service
	Create(indata utils.Map) (utils.Map, error)
	// Update - Update Service
	Update(materialTypeId string, indata utils.Map) (utils.Map, error)
	// Delete - Delete Service
	Delete(materialTypeId string, delete_permanent bool) error

	EndService()
}

// BrandService - Brand Service structure
type compareProdsBaseService struct {
	db_utils.DatabaseService
	dbRegion        db_utils.DatabaseService
	daoCompareProds sales_repository.CompareProdsDao
	daoBusiness     platform_repository.BusinessDao
	child           CompareProdsService
	businessId      string
}

// NewCompareProdsService - Construct CompareProds
func NewCompareProdsService(props utils.Map) (CompareProdsService, error) {
	funcode := sales_common.GetServiceModuleCode() + "M" + "01"

	log.Printf("CompareProdsService::Start ")
	// Verify whether the business id data passed
	businessId, err := utils.GetMemberDataStr(props, sales_common.FLD_BUSINESS_ID)
	if err != nil {
		return nil, err
	}

	p := compareProdsBaseService{}
	// Open Database Service
	err = p.OpenDatabaseService(props)
	if err != nil {
		return nil, err
	}

	// Open RegionDB Service
	p.dbRegion, err = platform_services.OpenRegionDatabaseService(props)
	if err != nil {
		p.CloseDatabaseService()
		return nil, err
	}

	// Assign the BusinessId
	p.businessId = businessId
	p.initializeService()

	_, err = p.daoBusiness.Get(businessId)
	if err != nil {
		err := &utils.AppError{
			ErrorCode:   funcode + "01",
			ErrorMsg:    "Invalid BusinessId",
			ErrorDetail: "Given BusinessId is not exist"}
		return p.errorReturn(err)
	}

	p.child = &p

	return &p, err
}

// compareProdsBaseService - Close all the services
func (p *compareProdsBaseService) EndService() {
	log.Printf("EndService ")
	p.CloseDatabaseService()
	p.dbRegion.CloseDatabaseService()
}

func (p *compareProdsBaseService) initializeService() {
	log.Printf("CompareProdsService:: GetBusinessDao ")
	p.daoBusiness = platform_repository.NewBusinessDao(p.GetClient())
	p.daoCompareProds = sales_repository.NewCompareProdsDao(p.dbRegion.GetClient(), p.businessId)
}

// List - List All records
func (p *compareProdsBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("compareProdsBaseService::FindAll - Begin")

	listdata, err := p.daoCompareProds.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("compareProdsBaseService::FindAll - End ")
	return listdata, nil
}

// Get - Find By Code
func (p *compareProdsBaseService) Get(materialTypeId string) (utils.Map, error) {
	log.Printf("compareProdsBaseService::Get::  Begin %v", materialTypeId)

	data, err := p.daoCompareProds.Get(materialTypeId)

	log.Println("BrandService::Get:: End ", err)
	return data, err
}

func (p *compareProdsBaseService) Find(filter string) (utils.Map, error) {
	fmt.Println("compareProdsBaseService::FindByCode::  Begin ", filter)

	data, err := p.daoCompareProds.Find(filter)
	log.Println("compareProdsBaseService::FindByCode:: End ", err)
	return data, err
}

// Create - Create Service
func (p *compareProdsBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("CompareProdsService::Create - Begin")

	// var materialTypeId string

	// dataval, dataok := indata[sales_common.FLD_MATERIAL_TYPE_ID]
	// if dataok {
	// 	materialTypeId = strings.ToLower(dataval.(string))
	// } else {
	// 	materialTypeId = utils.GenerateUniqueId("mate")
	// 	log.Println("Unique CompareProds ID", materialTypeId)
	// }

	// Assign Business Id
	indata[sales_common.FLD_BUSINESS_ID] = p.businessId
	//indata[sales_common.FLD_MATERIAL_TYPE_ID] = materialTypeId

	data, err := p.daoCompareProds.Create(indata)
	if err != nil {
		return utils.Map{}, err
	}

	log.Println("CompareProdsService::Create - End ")
	return data, nil
}

// Update - Update Service
func (p *compareProdsBaseService) Update(materialTypeId string, indata utils.Map) (utils.Map, error) {

	log.Println("CompareProdsService::Update - Begin")

	delete(indata, sales_common.FLD_MATERIAL_TYPE_ID)

	data, err := p.daoCompareProds.Update(materialTypeId, indata)

	log.Println("CompareProdsService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *compareProdsBaseService) Delete(materialTypeId string, delete_permanent bool) error {

	log.Println("BrandService::Delete - Begin", materialTypeId)

	if delete_permanent {
		result, err := p.daoCompareProds.Delete(materialTypeId)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(materialTypeId, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("BrandService::Delete - End")
	return nil
}

func (p *compareProdsBaseService) errorReturn(err error) (CompareProdsService, error) {
	// Close the Database Connection
	p.EndService()
	return nil, err
}
