package application

import (
	"context"
	"fmt"
	rs "insomnia/src/services/rule"
	domainRule "insomnia/src/services/rule/domain/rule"
	error2 "insomnia/src/services/rule/error"
	pb "insomnia/src/services/rule/proto/rule"
	"insomnia/src/services/rule/thirdparty"
	"time"
)

var (
	lockKeyAddProjectRule = "lock_key_add_project_rule_%d_%d"
)

type RuleService struct {
}

func (r *RuleService) AddProjectRule(ctx context.Context, userId uint64, ruleStr string, organizationId, projectId, ruleCategory int32) (*pb.AddProjectRuleReply, error) {
	// judge organization admin
	isAdmin, err := AllowToAddRule(ctx, userId, organizationId)
	if err != nil || !isAdmin {
		return nil, err
	}
	// avoid concurrent add rule
	lock, _ := rs.RedisCli.LockWait(ctx, fmt.Sprintf(lockKeyAddProjectRule, organizationId, projectId), time.Second*3, time.Minute*1)
	if lock == nil {
		return nil, error2.ConcurrentError
	}
	defer lock.Unlock(ctx)
	// check rule info valid
	if !CheckRuleValid(ruleStr, int(ruleCategory)) {
		return nil, error2.RuleError
	}
	err = domainRule.SvrRule.ProjectBindRule(ctx, organizationId, projectId, ruleStr, int(ruleCategory), int64(userId))
	if err != nil {
		return nil, err
	}
	return &pb.AddProjectRuleReply{}, nil
}

func (r *RuleService) GetUserProjectRule(ctx context.Context, userId uint64, organId int32) (*pb.GetProjectRuleReply, error) {
	projects, err := thirdparty.SvcOrganization.GetUserProjects(ctx, userId, organId)
	if err != nil {
		return nil, err
	}
	projectIds := make([]int32, len(projects))
	for idx, project := range projects {
		projectIds[idx] = project.Id
	}
	rules, err := MGetProjectInfo(ctx, organId, projectIds)
	if err != nil {
		return nil, err
	}
	projectRule := AssembleProjectRule(projects, rules)
	return &pb.GetProjectRuleReply{
		ProjectRules: projectRule,
	}, nil
}

var SvrRule *RuleService

func init() {
	SvrRule = &RuleService{}
}
