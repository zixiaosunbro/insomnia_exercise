package application

import (
	"context"
	error2 "insomnia/src/services/rule/error"
	"insomnia/src/services/rule/thirdparty"
)

func AllowToAddRule(ctx context.Context, userId uint64, organizationId int32) (bool, error) {
	isAdmin, err := thirdparty.SvcAuthentication.JudgeOrganizationAdmin(ctx, userId, organizationId)
	if err != nil {
		return false, err
	}
	if !isAdmin {
		return false, error2.PermissionError
	}
	return true, nil
}
