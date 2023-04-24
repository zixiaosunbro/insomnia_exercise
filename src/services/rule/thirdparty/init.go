package thirdparty

import "insomnia/src/services/rule/thirdparty/services"

var (
	SvcAuthentication services.IAuthenticationService
	SvcOrganization   services.IOrganizationService
)

func init() {
	SvcAuthentication = &services.AuthenticationService{}
	SvcOrganization = &services.OrganizationService{}
}
