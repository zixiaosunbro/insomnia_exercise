package rule

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"insomnia/src/pkg/config/middleware"
	"insomnia/src/pkg/utils"
	rs "insomnia/src/services/rule"
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

func TestRepositoryRule_SaveLintingRule(t *testing.T) {
	// case1: no data, add data
	rule := &LintingRule{
		OrganizationId: 1,
		ProjectId:      2,
		CreatorUid:     utils.RandomInt64Range(1, 1000),
		Rule:           "test",
		RuleCategory:   JsonRuleCategory,
	}
	err := RepoRule.SaveLintingRule(ctx, rule)
	assert.NoError(t, err)
	// case2: has data, update data
	rule.Rule = "test2"
	err = RepoRule.SaveLintingRule(ctx, rule)
	assert.NoError(t, err)
	equalFilter := map[string]interface{}{
		"organization_id": rule.OrganizationId,
		"project_id":      rule.ProjectId,
	}
	rules, err := RepoRule.GetLintingRule(ctx, equalFilter)
	assert.NoError(t, err)
	assert.Len(t, rules, 1)
	rule = rules[0]
	assert.Equal(t, "test2", rule.Rule)
	// del data
	err = DelLintingRule4Test(ctx, rule)
	assert.NoError(t, err)
}

func TestRepositoryRule_GetLintingRule(t *testing.T) {
	// case1: no data
	equalFilter := map[string]interface{}{
		"organization_id": 1,
		"project_id":      2,
	}
	rules, err := RepoRule.GetLintingRule(ctx, equalFilter)
	assert.NoError(t, err)
	assert.Empty(t, rules)
	// case2: has data
	rule := &LintingRule{
		OrganizationId: 1,
		ProjectId:      2,
		CreatorUid:     utils.RandomInt64Range(1, 1000),
		Rule:           "test",
		RuleCategory:   JsonRuleCategory,
	}
	err = RepoRule.SaveLintingRule(ctx, rule)
	assert.NoError(t, err)
	rules, err = RepoRule.GetLintingRule(ctx, equalFilter)
	assert.NoError(t, err)
	assert.Len(t, rules, 1)
	// del data
	err = DelLintingRule4Test(ctx, rule)
	assert.NoError(t, err)
}

func TestRepositoryRule_SaveProjectRule(t *testing.T) {
	// case1: no data, add data
	projectRule := &ProjectLintingRule{
		OrganizationId: 1,
		ProjectId:      2,
		RuleId:         3,
		OperatorUid:    utils.RandomInt64Range(1, 1000),
	}
	err := RepoRule.SaveProjectRule(ctx, projectRule)
	assert.NoError(t, err)
	// case2: has data, update data
	projectRule.RuleId = 4
	err = RepoRule.SaveProjectRule(ctx, projectRule)
	assert.NoError(t, err)
	equalFilter := map[string]interface{}{
		"organization_id": 1,
		"project_id":      2,
	}
	rules, err := RepoRule.GetProjectRule(ctx, equalFilter)
	assert.NoError(t, err)
	assert.Len(t, rules, 1)
	projectRule = rules[0]
	assert.Equal(t, int64(4), projectRule.RuleId)
	// del data
	err = DelProjectRule4Test(ctx, projectRule)
	assert.NoError(t, err)

}

func TestRepositoryRule_GetProjectRule(t *testing.T) {
	// case1: no data
	equalFilter := map[string]interface{}{
		"organization_id": 1,
		"project_id":      2,
	}
	rules, err := RepoRule.GetProjectRule(ctx, equalFilter)
	assert.NoError(t, err)
	assert.Empty(t, rules)
	// case2: has data
	projectRule := &ProjectLintingRule{
		OrganizationId: 1,
		ProjectId:      2,
		RuleId:         3,
		OperatorUid:    utils.RandomInt64Range(1, 1000),
	}
	err = RepoRule.SaveProjectRule(ctx, projectRule)
	assert.NoError(t, err)
	rules, err = RepoRule.GetProjectRule(ctx, equalFilter)
	assert.NoError(t, err)
	assert.Len(t, rules, 1)
	// del data
	err = DelProjectRule4Test(ctx, projectRule)
}

func TestRepositoryRule_GetProjectRuleCache(t *testing.T) {
	organId, projectId := int32(1), int32(2)
	// case1: no data
	cache, err := RepoRule.GetProjectRuleCache(ctx, organId, projectId)
	assert.Error(t, err)
	assert.Empty(t, cache)
	// case2: has data
	category, rule := JsonRuleCategory, "test"
	err = RepoRule.SaveProjectRuleCache(ctx, organId, projectId, category, rule)
	assert.NoError(t, err)
	cache, err = RepoRule.GetProjectRuleCache(ctx, organId, projectId)
	assert.NoError(t, err)
	assert.NotNil(t, cache)
	// del data
	_, err = rs.RedisCli.Del(ctx, cache.CacheTableName()).Result()
	assert.NoError(t, err)
}

func TestRepositoryRule_SaveProjectRuleCache(t *testing.T) {
	organId, projectId, category, rule := int32(1), int32(2), JsonRuleCategory, "test"
	err := RepoRule.SaveProjectRuleCache(ctx, organId, projectId, category, rule)
	assert.NoError(t, err)
	// del data
	cache := &ProjectRuleCache{
		OrganizationId: organId,
		ProjectId:      projectId,
	}
	_, err = rs.RedisCli.Del(ctx, cache.CacheTableName()).Result()
	assert.NoError(t, err)
}

func DelLintingRule4Test(ctx context.Context, rule ...*LintingRule) error {
	err := rs.DBMaster.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(rule).Error
		return err
	})
	return err
}

func DelProjectRule4Test(ctx context.Context, rule ...*ProjectLintingRule) error {
	err := rs.DBMaster.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(rule).Error
		return err
	})
	return err
}
