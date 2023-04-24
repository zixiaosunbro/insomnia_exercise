package services

import "context"

// IAuthenticationService assumption that exists an authentication service can use to verify the identity and role of a user via session token
type IAuthenticationService interface {
	JudgeOrganizationAdmin(ctx context.Context, userId uint64, organizationId int32) (bool, error)
}

type AuthenticationService struct {
}

func (a *AuthenticationService) JudgeOrganizationAdmin(ctx context.Context, userId uint64, organizationId int32) (bool, error) {
	// to support api test, mock data.
	var isAdmin bool
	if organizationId == 101 {
		isAdmin = userId == 1010
	}
	return isAdmin, nil
}
