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

	clientId := dataAuth[auth_common.CLIENT_ID].(string)
	clientSecret := dataAuth[auth_common.CLIENT_SECRET].(string)

	// Authenticate with AppClient tables
	appClientData, err := auth_services.AuthenticateAppClient(dbProps, clientId, clientSecret)
	if err != nil {
		return nil, err
	}
	businessId := appClientData[auth_common.CLIENT_SCOPE].(string)
	log.Println("Auth Client Record ", appClientData, err)

	dataAuth[platform_common.FLD_CLIENT_TYPE] = appClientData[platform_common.FLD_CLIENT_TYPE].(string)
	dataAuth[platform_common.FLD_CLIENT_SCOPE] = appClientData[platform_common.FLD_CLIENT_SCOPE].(string)

	mapScopes := utils.Map{}
	if scopeValue, scopeOk := dataAuth[auth_common.SCOPE]; scopeOk && scopeValue.(string) != "" {
		mapScopes = auth_services.ParseScope(scopeValue.(string))
	}

	log.Println("Scopes ", mapScopes)

	if dataAuth[auth_common.GRANT_TYPE] == auth_common.GRANT_TYPE_PASSWORD {

		// userType with "App" or "Business"
		if userType, userTypeOk := mapScopes[auth_common.USER_TYPE]; userTypeOk &&
			(userType.(string) == auth_common.USER_TYPE_CUSTOMER) {

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
			// Authenticate AppUser
			custData, err := authenticateCustomer(dbProps, businessId, authKey, authKeyValue, authPassword)
			if err != nil {
				return utils.Map{}, err
			}

			dataAuth[auth_common.USER_ID] = custData[sales_common.FLD_CUSTOMER_ID].(string)

		} else {
			// No UserType or Other UserTypes
			err := &utils.AppError{ErrorStatus: 401, ErrorMsg: "Invalid UserType", ErrorDetail: "UserType is invalid"}
			return utils.Map{}, err
		}
	} /*else if dataAuth[GRANT_TYPE] == GRANT_TYPE_REFRESH {
		dataAuth.RefreshToken = ctx.FormValue("refresh_token")
		err := &utils.AppError{ErrorStatus: 401, ErrorMsg: "Client DB Connection Error", ErrorDetail: "Client DB Connection Error"}
		return utils.Map{}, err
	}*/

	log.Printf("Auth Values %v", dataAuth)
	return dataAuth, nil
}

func authenticateCustomer(dbProps utils.Map, businessId string, auth_key string, auth_key_value string, auth_password string) (utils.Map, error) {

	// Append Business Idl
	dbProps[sales_common.FLD_BUSINESS_ID] = businessId

	// User Validation
	svcCustomer, err := sales_services.NewCustomerService(dbProps)
	if err != nil {
		err := &utils.AppError{ErrorStatus: 417, ErrorMsg: "Status Expectation Failed", ErrorDetail: "Authentication Failure"}
		return utils.Map{}, err
	}
	defer svcCustomer.EndService()

	log.Println("Business::Auth:: Parameter Value ", auth_key, auth_key_value)
	appUserData, err := svcCustomer.Authenticate(auth_key, auth_key_value, auth_password)
	if err != nil {

		err := &utils.AppError{ErrorStatus: 401, ErrorMsg: "Status Unauthorized", ErrorDetail: "Authentication Failure"}
		return utils.Map{}, err
	}

	return appUserData, nil
}
