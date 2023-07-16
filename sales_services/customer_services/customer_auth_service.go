package customer_services

import (
	"log"

	"github.com/zapscloud/golib-auth/auth_common"
	"github.com/zapscloud/golib-auth/auth_services"
	"github.com/zapscloud/golib-platform/platform_common"
	"github.com/zapscloud/golib-sales/sales_common"
	"github.com/zapscloud/golib-sales/sales_services"
	"github.com/zapscloud/golib-utils/utils"
)

func ValidateAuthCredential(dbProps utils.Map, dataAuth utils.Map) (utils.Map, error) {

	log.Printf("ValidateAppAuth %v", dataAuth)

	// Authenticate with Clients tables
	_, _, clientData, err := auth_services.AuthenticateClient(dbProps, dataAuth)
	if err != nil {
		return nil, err
	}
	log.Println("Auth Client Record ", clientData, err)

	// Update Client Data in AuthData
	dataAuth[platform_common.FLD_CLIENT_TYPE] = clientData[platform_common.FLD_CLIENT_TYPE].(string)
	dataAuth[platform_common.FLD_CLIENT_SCOPE] = clientData[platform_common.FLD_CLIENT_SCOPE].(string)

	// Get the GrantType
	grantType := dataAuth[auth_common.GRANT_TYPE].(string)

	// Get Scope values if anything passed
	mapScopes := auth_services.ParseScope(dataAuth)
	log.Println("Scopes ", mapScopes)
	switch grantType {
	//
	// ============[ Grant_Type: Client Credentials ] ========================================
	case auth_common.GRANT_TYPE_CLIENT_CREDENTIALS:
		// Client Credentials not support for Customers
		err = &utils.AppError{ErrorStatus: 417, ErrorMsg: "Status Expectation Failed", ErrorDetail: "Authentication Failure"}
		return utils.Map{}, err

	//
	// ============[ Grant_Type: Password Credentials ] ======================================
	case auth_common.GRANT_TYPE_PASSWORD:

		// For all other cases like WebApp, MobileApp and etc
		businessId, err := utils.GetMemberDataStr(mapScopes, platform_common.FLD_BUSINESS_ID)
		if err != nil {
			return nil, err
		}

		// Validate BusinessId is exist
		if !utils.IsEmpty(businessId) {
			_, err = auth_services.IsBusinessExist(dbProps, businessId)
			if err != nil {
				return nil, err
			}
		}

		// Assign BusinessId in AuthData
		dataAuth[platform_common.FLD_BUSINESS_ID] = businessId

		// Authenticate Customer
		custData, err := authenticateCustomer(dbProps, businessId, dataAuth)
		if err != nil {
			return utils.Map{}, err
		}

		dataAuth[auth_common.USER_ID] = custData[sales_common.FLD_CUSTOMER_ID].(string)

	//
	// ============[ Grant_Type: REFRESH ] ========================================
	case auth_common.GRANT_TYPE_REFRESH:
		/* Need to Implement Refersh Token */
		//dataAuth.RefreshToken = ctx.FormValue("refresh_token")

	}

	log.Printf("Auth Values %v", dataAuth)
	return dataAuth, nil
}

func authenticateCustomer(dbProps utils.Map, businessId string, dataAuth utils.Map) (utils.Map, error) {

	// Append Business Id
	dbProps[sales_common.FLD_BUSINESS_ID] = businessId

	// User Validation
	svcCustomer, err := sales_services.NewCustomerService(dbProps)
	if err != nil {
		err := &utils.AppError{ErrorStatus: 417, ErrorMsg: "Status Expectation Failed", ErrorDetail: "Authentication Failure"}
		return utils.Map{}, err
	}
	defer svcCustomer.EndService()

	// Set default authKey as Customer LoginId
	authKey := sales_common.FLD_CUSTOMER_LOGIN_ID
	// if loginType, loginTypeOK := mapScopes[auth_common.LOGIN_TYPE]; loginTypeOK {

	// 	loginType = loginType.(string)

	// 	if loginType == auth_common.LOGIN_TYPE_PHONE {
	// 		authKey = platform_common.FLD_APP_USER_PHONE
	// 	} else if loginType == auth_common.LOGIN_TYPE_EMAIL {
	// 		authKey = platform_common.FLD_APP_USER_EMAILID
	// 	}
	// }

	authKeyValue := dataAuth[auth_common.USERNAME].(string)
	authPassword := dataAuth[auth_common.PASSWORD].(string)

	log.Println("Business::Auth:: Parameter Value ", authKey, authKeyValue)
	appUserData, err := svcCustomer.Authenticate(authKey, authKeyValue, authPassword)
	if err != nil {

		err := &utils.AppError{ErrorStatus: 401, ErrorMsg: "Status Unauthorized", ErrorDetail: "Authentication Failure"}
		return utils.Map{}, err
	}

	return appUserData, nil
}
