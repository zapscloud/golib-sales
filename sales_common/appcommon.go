package sales_common

import (
	"log"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-platform/platform_common"
)

// Product Module tables
const (
	// Database Prefix
	DbPrefix = db_common.DB_COLLECTION_PREFIX
	// Collection Names
	DbRegions           = DbPrefix + "sales_regions"
	DbBanners           = DbPrefix + "sales_banners"
	DbBrands            = DbPrefix + "sales_brands"
	DbCatalogues        = DbPrefix + "sales_catalogues"
	DbCategories        = DbPrefix + "sales_categories"
	DbProducts          = DbPrefix + "sales_products"
	DbTestimonials      = DbPrefix + "sales_testimonials"
	DbBlogs             = DbPrefix + "sales_blogs"
	DbCustomers         = DbPrefix + "sales_customers"
	DbCustomerOrders    = DbPrefix + "sales_customer_orders"
	DbCustomerCarts     = DbPrefix + "sales_customer_carts"
	DbCustomerWishlists = DbPrefix + "sales_customer_wishlists"
	DbCustomerReviews   = DbPrefix + "sales_customer_reviews"
	DbPolicies          = DbPrefix + "sales_policies"
	DbPayments          = DbPrefix + "sales_payments"
	DbNavigations       = DbPrefix + "sales_navigations"
	DbPreferences       = DbPrefix + "sales_preferences"
	DbPages             = DbPrefix + "sales_pages"
	DbDealers           = DbPrefix + "sales_dealers"
	DbDiscounts         = DbPrefix + "sales_discounts"
	DbMedias            = DbPrefix + "sales_medias"
)

// Product Module table fields
const (
	// Common fields for all tables
	FLD_BUSINESS_ID = platform_common.FLD_BUSINESS_ID
	FLD_SEO_KEYID   = "seo_key_id"

	// Fields for Region
	FLD_REGION_ID           = "sales_region_id"
	FLD_REGION_NAME         = "sales_region_name"
	FLD_REGION_PINCODES     = "sales_region_pincodes"
	FLD_REGION_PINCODE_FROM = "pincode_from"
	FLD_REGION_PINCODE_TO   = "pincode_to"

	// Fields for Banner
	FLD_BANNER_ID   = "banner_id"
	FLD_BANNER_NAME = "banner_name"

	// Fields for Cart
	FLD_CART_ID = "cart_id"

	// Fields for Brand Table
	FLD_BRAND_ID   = "brand_id"
	FLD_BRAND_NAME = "brand_name"

	// Fields for Category Table
	FLD_CATALOGUE_ID   = "catalogue_id"
	FLD_CATALOGUE_NAME = "catalogue_name"

	// Fields for Category Table
	FLD_CATEGORY_ID   = "category_id"
	FLD_CATEGORY_NAME = "category_name"

	// Fields for Product Table
	FLD_PRODUCT_ID   = "product_id"
	FLD_PRODUCT_NAME = "product_name"

	// Fields for Testimonial
	FLD_TESTIMONIAL_ID   = "testimonial_id"
	FLD_TESTIMONIAL_NAME = "testimonial_name"

	//Fields for Blog
	FLD_BLOG_ID   = "blog_id"
	FLD_BLOG_NAME = "blog_name"

	//Fields for CustomerOrder
	FLD_CUSTOMER_ORDER_ID   = "customer_order_id"
	FLD_CUSTOMER_ORDER_NO   = "customer_order_no" // Auto generated Sequence Number
	FLD_CUSTOMER_ORDER_NAME = "customer_order_name"

	// Fields for Customer Table
	FLD_CUSTOMER_ID       = "customer_id"
	FLD_CUSTOMER_LOGIN_ID = "customer_loginid"
	FLD_CUSTOMER_PASSWORD = "customer_password"
	FLD_CUSTOMER_OTP      = "customer_otp"

	// Fields for Order
	FLD_ORDER_ID   = "order_id"
	FLD_ORDER_NAME = "order_name"

	// Fields for Payment
	FLD_PAYMENT_ID   = "payment_id"
	FLD_PAYMENT_NAME = "payment_name"

	// Fields for Policies
	FLD_POLICY_ID   = "policy_id"
	FLD_POLICY_TYPE = "policy_type"
	FLD_POLICY_NAME = "policy_name"

	// Fields for Navigation
	FLD_NAVIGATION_ID   = "navigation_id"
	FLD_NAVIGATION_NAME = "navigation_name"

	// Fields for Preference
	FLD_PREFERENCE_ID   = "preference_id"
	FLD_PREFERENCE_NAME = "preference_name"

	// Fields for Page
	FLD_PAGE_ID = "page_id"

	// Fields for Dealer
	FLD_DEALER_ID   = "dealer_id"
	FLD_DEALER_NAME = "dealer_name"

	// Fields for Review/Feedback
	FLD_REVIEW_ID = "review_id"

	// Fields for WhishList
	FLD_WISHLIST_ID = "wishlist_id"

	// Fields for Media
	FLD_MEDIA_ID = "media_id"

	// Fields for Discount
	FLD_DISCOUNT_ID = "discount_id"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)

}

func GetServiceModuleCode() string {
	return "SALES"
}

//
// Indexes
//
//
// db.zc_sales_region.createIndex({"sales_region_pincodes.pincode_from": 1}, {"sales_region_pincodes.pincode_to": 1})
