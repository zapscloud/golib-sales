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

// ProductService - Business Product Service structure
type ProductService interface {
	// List - List All records
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Get - Find By Code
	Get(productId string) (utils.Map, error)
	// Find - Find the item
	Find(filter string) (utils.Map, error)
	// Create - Create Service
	Create(indata utils.Map) (utils.Map, error)
	// Update - Update Service
	Update(productId string, indata utils.Map) (utils.Map, error)
	// Delete - Delete Service
	Delete(productId string, delete_permanent bool) error

	EndService()
}

// ProductService - Business Product Service structure
type productBaseService struct {
	db_utils.DatabaseService
	daoProduct  sales_repository.ProductDao
	daoBusiness platform_repository.BusinessDao
	child       ProductService
	businessId  string
}

// NewProductService - Construct Product
func NewProductService(props utils.Map) (ProductService, error) {
	funcode := sales_common.GetServiceModuleCode() + "M" + "01"

	p := productBaseService{}
	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ProductService ")
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
		err := &utils.AppError{
			ErrorCode:   funcode + "01",
			ErrorMsg:    "Invalid business_id",
			ErrorDetail: "Given app_business_id is not exist"}
		return nil, err
	}

	p.child = &p

	return &p, err
}

// EndLoyaltyCardService - Close all the services
func (p *productBaseService) EndService() {
	log.Printf("EndProductService ")
	p.CloseDatabaseService()
}

func (p *productBaseService) initializeService() {
	log.Printf("ProductMongoService:: GetBusinessDao ")
	p.daoProduct = sales_repository.NewProductDao(p.GetClient(), p.businessId)
	p.daoBusiness = platform_repository.NewBusinessDao(p.GetClient())
}

// List - List All records
func (p *productBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("ProductService::FindAll - Begin")

	listdata, err := p.daoProduct.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("ProductService::FindAll - End ")
	return listdata, nil
}

// Get - Find By Code
func (p *productBaseService) Get(productId string) (utils.Map, error) {
	log.Printf("ProductService::Get::  Begin %v", productId)

	data, err := p.daoProduct.Get(productId)

	log.Println("ProductService::Get:: End ", data, err)
	return data, err
}

func (p *productBaseService) Find(filter string) (utils.Map, error) {
	fmt.Println("productBaseService::FindByCode::  Begin ", filter)

	data, err := p.daoProduct.Find(filter)
	log.Println("productBaseService::FindByCode:: End ", data, err)
	return data, err
}

// Create - Create Service
func (p *productBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("ProductService::Create - Begin")
	var productId string

	dataval, dataok := indata[sales_common.FLD_PRODUCT_ID]
	if dataok {
		productId = strings.ToLower(dataval.(string))
	} else {
		productId = utils.GenerateUniqueId("prod")
		log.Println("Unique Product ID", productId)
	}

	//BusinessProduct
	indata[sales_common.FLD_BUSINESS_ID] = p.businessId
	indata[sales_common.FLD_PRODUCT_ID] = productId

	data, err := p.daoProduct.Create(indata)
	if err != nil {
		return utils.Map{}, err
	}

	log.Println("ProductService::Create - End ")
	return data, nil
}

// Update - Update Service
func (p *productBaseService) Update(productId string, indata utils.Map) (utils.Map, error) {

	log.Println("BusinessProdcutService::Update - Begin")

	data, err := p.daoProduct.Update(productId, indata)

	log.Println("ProductService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *productBaseService) Delete(productId string, delete_permanent bool) error {

	log.Println("ProductService::Delete - Begin", productId)

	if delete_permanent {
		result, err := p.daoProduct.Delete(productId)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(productId, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("ProductService::Delete - End")
	return nil
}
