package rule

import (
	"context"
	"insomnia/src/pkg/middleware_custom/mysql"
	"insomnia/src/pkg/utils"
	rs "insomnia/src/services/rule"
)

type ServiceRule struct {
}

func (s *ServiceRule) ProjectBindRule(ctx context.Context, organId, projectId int32, ruleInfo string, ruleCategory int, userId int64) error {
	tx := rs.DBMaster.Begin()
	err := tx.Error
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	// add rule
	rule := CreateLintingRule(organId, projectId, userId, ruleInfo, ruleCategory)
	if err = RepoRule.SaveLintingRule(ctx, rule, mysql.WithTx(tx)); err != nil {
		return err
	}
	// bind rule with project
	equalFilter := map[string]interface{}{
		"organization_id": organId,
		"project_id":      projectId,
	}
	projectRules, err := RepoRule.GetProjectRule(ctx, equalFilter, mysql.WithTx(tx))
	if err != nil {
		return err
	}
	var projectRule *ProjectLintingRule
	if len(projectRules) == 0 {
		// create project rule
		projectRule = CreateProjectLintingRule(organId, projectId, rule.ID, userId)
	} else {
		projectRule = projectRules[0]
		projectRule.RuleId = rule.ID
	}
	if err = RepoRule.SaveProjectRule(ctx, projectRule, mysql.WithTx(tx)); err != nil {
		return err
	}
	if err = tx.Commit().Error; err != nil {
		return err
	}
	// update redis
	_ = RepoRule.SaveProjectRuleCache(ctx, organId, projectId, ruleCategory, ruleInfo)
	return nil
}

func (s *ServiceRule) GetProjectInfo(ctx context.Context, organId, projectId int32) (*ProjectRuleCache, error) {
	// get rule from redis
	ruleCache, err := RepoRule.GetProjectRuleCache(ctx, organId, projectId)
	if err != nil && !utils.RedisReturnNil(err) {
		return nil, err
	}
	if utils.RedisReturnNil(err) {
		// get rule from mysql
		equalFilter := map[string]interface{}{
			"organization_id": organId,
			"project_id":      projectId,
		}
		var projectRules []*ProjectLintingRule
		if projectRules, err = RepoRule.GetProjectRule(ctx, equalFilter); err != nil {
			return nil, err
		}
		if len(projectRules) == 0 {
			return nil, nil
		}
		var ruleInfos []*LintingRule
		if ruleInfos, err = RepoRule.GetLintingRule(ctx, map[string]any{"id": projectRules[0].RuleId}); err != nil {
			return nil, err
		}
		if len(ruleInfos) == 0 {
			return nil, nil
		}
		ruleCache = &ProjectRuleCache{
			OrganizationId: organId,
			ProjectId:      projectId,
			Rule:           ruleInfos[0].Rule,
			Category:       ruleInfos[0].RuleCategory,
		}
		// save to redis
		_ = RepoRule.SaveProjectRuleCache(ctx, organId, projectId, ruleInfos[0].RuleCategory, ruleInfos[0].Rule)
	}
	return ruleCache, nil

}
