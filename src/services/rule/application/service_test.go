package application

import (
	"bou.ke/monkey"
	"context"
	"github.com/stretchr/testify/assert"
	"insomnia/src/pkg/config/middleware"
	"insomnia/src/pkg/utils"
	rs "insomnia/src/services/rule"
	domainRule "insomnia/src/services/rule/domain/rule"
	"insomnia/src/services/rule/thirdparty/services"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	middleware.Init()
	rs.ResourceInit()
	m.Run()
}

var (
	ctx = context.Background()
)

func TestRuleService_AddProjectRule(t *testing.T) {
	userId, ruleStr, organId, projectId, category := uint64(1), "{}", int32(1), int32(1), int32(1)
	monkey.Patch(AllowToAddRule, func(ctx context.Context, userId uint64, organizationId int32) (bool, error) {
		return true, nil
	})
	defer monkey.UnpatchAll()
	var ruleService *domainRule.ServiceRule
	monkey.PatchInstanceMethod(reflect.TypeOf(ruleService), "ProjectBindRule", func(_ *domainRule.ServiceRule, ctx context.Context, organId, projectId int32, ruleStr string, category int, userId int64) error {
		return nil
	})
	_, err := SvrRule.AddProjectRule(ctx, userId, ruleStr, organId, projectId, category)
	assert.NoError(t, err)
}

func TestRuleService_GetUserProjectRule(t *testing.T) {
	userId, organId := uint64(1), int32(1)
	var organSvc *services.OrganizationService
	monkey.PatchInstanceMethod(reflect.TypeOf(organSvc), "GetUserProjects", func(_ *services.OrganizationService, ctx context.Context, userId uint64, organId int32) ([]*services.ProjectInfo, error) {
		return []*services.ProjectInfo{
			{
				Id:   1,
				Name: "test",
			},
		}, nil
	})
	defer monkey.UnpatchAll()
	monkey.Patch(MGetProjectInfo, func(ctx context.Context, organId int32, projectIds []int32) ([]*domainRule.ProjectRuleCache, error) {
		return []*domainRule.ProjectRuleCache{
			{
				OrganizationId: organId,
				ProjectId:      1,
				Rule:           utils.RandomString(10),
				Category:       1,
			},
		}, nil
	})
	res, err := SvrRule.GetUserProjectRule(ctx, userId, organId)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
