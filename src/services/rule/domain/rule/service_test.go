package rule

import (
	"bou.ke/monkey"
	"context"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"insomnia/src/pkg/middleware_custom/mysql"
	"insomnia/src/pkg/utils"
	rs "insomnia/src/services/rule"
	"reflect"
	"testing"
)

func TestServiceRule_ProjectBindRule(t *testing.T) {
	organId, projectId, ruleInfo, ruleCategory, userId := int32(1), int32(1), "ruleInfo", 1, utils.RandomInt64Range(1, 2000)
	// case1: project not bind rule
	err := SvrRule.ProjectBindRule(ctx, organId, projectId, ruleInfo, ruleCategory, userId)
	assert.NoError(t, err)
	// case2: project already bind rule, update
	err = SvrRule.ProjectBindRule(ctx, organId, projectId, utils.RandomString(30), YamlRuleCategory, userId)
	assert.NoError(t, err)
	equalFilter := map[string]interface{}{
		"organization_id": organId,
		"project_id":      projectId,
	}
	projectRules, err := RepoRule.GetProjectRule(ctx, equalFilter)
	assert.NoError(t, err)
	assert.Len(t, projectRules, 1)
	// del data
	cache := &ProjectRuleCache{
		OrganizationId: organId,
		ProjectId:      projectId,
	}
	_, err = rs.RedisCli.Del(ctx, cache.CacheTableName()).Result()
	assert.NoError(t, err)
	err = rs.DBMaster.Transaction(func(tx *gorm.DB) error {
		err = tx.Where("organization_id = ? and project_id = ?", organId, projectId).Delete(&ProjectLintingRule{}).Error
		if err != nil {
			return err
		}
		err = tx.Where("organization_id = ? and project_id = ?", organId, projectId).Delete(&LintingRule{}).Error
		if err != nil {
			return err
		}
		return nil
	})
	assert.NoError(t, err)
}

func TestServiceRule_GetProjectInfo(t *testing.T) {
	organId, projectId := int32(1), int32(1)
	var repoRule *RepositoryRule
	monkey.PatchInstanceMethod(reflect.TypeOf(repoRule), "GetProjectRule", func(_ *RepositoryRule, ctx context.Context, equalFilter map[string]any, opts ...mysql.DBOption) ([]*ProjectLintingRule, error) {
		return []*ProjectLintingRule{
			{
				RuleId: 10,
			},
		}, nil
	})
	defer monkey.UnpatchAll()
	monkey.PatchInstanceMethod(reflect.TypeOf(repoRule), "GetLintingRule", func(_ *RepositoryRule, ctx context.Context, equalFilter map[string]any, opts ...mysql.DBOption) ([]*LintingRule, error) {
		return []*LintingRule{
			{
				Rule:         "ruleInfo",
				RuleCategory: JsonRuleCategory,
			},
		}, nil
	})
	ruleCache, err := SvrRule.GetProjectInfo(ctx, organId, projectId)
	assert.NoError(t, err)
	assert.NotNil(t, ruleCache)
	// del data
	cache := &ProjectRuleCache{
		OrganizationId: organId,
		ProjectId:      projectId,
	}
	_, err = rs.RedisCli.Del(ctx, cache.CacheTableName()).Result()
	assert.NoError(t, err)

}
