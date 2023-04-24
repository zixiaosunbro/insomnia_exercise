package application

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v3"
	"insomnia/src/services/rule/domain/rule"
	pb "insomnia/src/services/rule/proto/rule"
	"insomnia/src/services/rule/thirdparty/services"
)

func CheckRuleValid(ruleStr string, category int) bool {
	if category == rule.JsonRuleCategory {
		return jsoniter.Valid([]byte(ruleStr))
	} else if category == rule.YamlRuleCategory {
		var data any
		err := yaml.Unmarshal([]byte(ruleStr), &data)

		return err == nil
	}
	return true
}

func AssembleProjectRule(projects []*services.ProjectInfo, projectRule []*rule.ProjectRuleCache) []*pb.ProjectRule {
	projectRuleMap := make(map[int32]*rule.ProjectRuleCache, len(projectRule))
	for _, cache := range projectRule {
		projectRuleMap[cache.ProjectId] = cache
	}
	result := make([]*pb.ProjectRule, 0)
	for _, project := range projects {
		ruleInfo, ok := projectRuleMap[project.Id]
		if !ok {
			continue
		}
		result = append(result, &pb.ProjectRule{
			OrganId:      ruleInfo.OrganizationId,
			ProjectId:    project.Id,
			ProjectName:  project.Name,
			Rule:         ruleInfo.Rule,
			RuleCategory: int32(ruleInfo.Category),
		})
	}
	return result
}

func MGetProjectInfo(ctx context.Context, organId int32, projectIds []int32) ([]*rule.ProjectRuleCache, error) {
	if len(projectIds) == 0 {
		return nil, nil
	}
	ruleCaches := make([]*rule.ProjectRuleCache, 0)
	type ruleInfo struct {
		ruleCache *rule.ProjectRuleCache
		err       error
	}
	resultChan := make(chan *ruleInfo, len(projectIds))
	for _, pId := range projectIds {
		go func(projectId int32) {
			ruleCache, err := rule.SvrRule.GetProjectInfo(ctx, organId, projectId)
			resultChan <- &ruleInfo{
				ruleCache: ruleCache,
				err:       err,
			}
		}(pId)
	}
	for i := 0; i < len(projectIds); i++ {
		cache := <-resultChan
		if cache.err != nil || cache.ruleCache == nil {
			continue
		}
		ruleCaches = append(ruleCaches, cache.ruleCache)
	}

	return ruleCaches, nil
}
