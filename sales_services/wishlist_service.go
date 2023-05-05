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

type WishlistService interface {
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

type wishlistBaseService struct {
	db_utils.DatabaseService
	daoWishlist sales_repository.WishlistDao
	daoBusiness platform_repository.BusinessDao
	child       WishlistService
	businessId  string
}

// NewWishlistService - Construct Wishlist
func NewWishlistService(props utils.Map) (WishlistService, error) {
	funcode := sales_common.GetServiceModuleCode() + "M" + "01"

	p := wishlistBaseService{}
	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("WishlistService ")
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
func (p *wishlistBaseService) EndService() {
	log.Printf("EndService ")
	p.CloseDatabaseService()
}

func (p *wishlistBaseService) initializeService() {
	log.Printf("WishlistService:: GetBusinessDao ")
	p.daoWishlist = sales_repository.NewWishlistDao(p.GetClient(), p.businessId)
	p.daoBusiness = platform_repository.NewBusinessDao(p.GetClient())
}

// List - List All records
func (p *wishlistBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("wishlistBaseService::FindAll - Begin")

	listdata, err := p.daoWishlist.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("wishlistBaseService::FindAll - End ")
	return listdata, nil
}

// Get - Find By Code
func (p *wishlistBaseService) Get(wishlistId string) (utils.Map, error) {
	log.Printf("wishlistBaseService::Get::  Begin %v", wishlistId)

	data, err := p.daoWishlist.Get(wishlistId)

	log.Println("wishlistBaseService::Get:: End ", data, err)
	return data, err
}

func (p *wishlistBaseService) Find(filter string) (utils.Map, error) {
	fmt.Println("wishlistBaseService::FindByCode::  Begin ", filter)

	data, err := p.daoWishlist.Find(filter)
	log.Println("wishlistBaseService::FindByCode:: End ", data, err)
	return data, err
}

// Create - Create Service
func (p *wishlistBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("WishlistService::Create - Begin")
	var wishlistId string

	dataval, dataok := indata[sales_common.FLD_WISHLIST_ID]
	if dataok {
		wishlistId = strings.ToLower(dataval.(string))
	} else {
		wishlistId = utils.GenerateUniqueId("wish")
		log.Println("Unique Wishlist ID", wishlistId)
	}

	// Assign BusinessId
	indata[sales_common.FLD_BUSINESS_ID] = p.businessId
	indata[sales_common.FLD_WISHLIST_ID] = wishlistId

	data, err := p.daoWishlist.Create(indata)
	if err != nil {
		return utils.Map{}, err
	}

	log.Println("WishlistService::Create - End ")
	return data, nil
}

// Update - Update Service
func (p *wishlistBaseService) Update(wishlistId string, indata utils.Map) (utils.Map, error) {

	log.Println("WishlistService::Update - Begin")

	data, err := p.daoWishlist.Update(wishlistId, indata)

	log.Println("WishlistService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *wishlistBaseService) Delete(wishlistId string, delete_permanent bool) error {

	log.Println("WishlistService::Delete - Begin", wishlistId)

	if delete_permanent {
		result, err := p.daoWishlist.Delete(wishlistId)
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

	log.Printf("WishlistService::Delete - End")
	return nil
}
