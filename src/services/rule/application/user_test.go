package application

import (
	"bou.ke/monkey"
	"context"
	"github.com/stretchr/testify/assert"
	"insomnia/src/services/rule/thirdparty/services"
	"reflect"
	"testing"
)

func TestAllowToAddRule(t *testing.T) {
	userId, organizationId := uint64(1), int32(1)
	var authSvc *services.AuthenticationService
	monkey.PatchInstanceMethod(reflect.TypeOf(authSvc), "JudgeOrganizationAdmin", func(_ *services.AuthenticationService, ctx context.Context, userId uint64, organizationId int32) (bool, error) {
		return true, nil
	})
	defer monkey.UnpatchAll()
	isAdmin, err := AllowToAddRule(context.Background(), userId, organizationId)
	assert.NoError(t, err)
	assert.True(t, isAdmin)

}
