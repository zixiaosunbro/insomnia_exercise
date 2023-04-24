package handler

import (
	"context"
	"insomnia/src/pkg/utils"
	"insomnia/src/services/rule/application"
	"insomnia/src/services/rule/domain/rule"
	error2 "insomnia/src/services/rule/error"
	pb "insomnia/src/services/rule/proto/rule"
)

type RuleDispatcher struct {
}

func (r *RuleDispatcher) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: req.Message}, nil
}

func (r *RuleDispatcher) AddProjectRule(ctx context.Context, req *pb.AddProjectRuleRequest) (*pb.AddProjectRuleReply, error) {
	if req.UserId == 0 || req.OrganizationId == 0 || req.ProjectId == 0 || req.Rule == "" ||
		!utils.ElementInSlice([]int32{rule.JsonRuleCategory, rule.YamlRuleCategory}, req.RuleCategory) {
		return nil, error2.ParamError
	}
	return application.SvrRule.AddProjectRule(ctx, req.UserId, req.Rule, req.OrganizationId, req.ProjectId, req.RuleCategory)
}

func (r *RuleDispatcher) UpdateProjectRule(ctx context.Context, req *pb.UpdateProjectRuleRequest) (*pb.UpdateProjectRuleReply, error) {
	if req.UserId == 0 || req.OrganizationId == 0 || req.ProjectId == 0 || req.Rule == "" ||
		!utils.ElementInSlice([]int32{rule.JsonRuleCategory, rule.YamlRuleCategory}, req.RuleCategory) {
		return nil, error2.ParamError
	}
	_, err := application.SvrRule.AddProjectRule(ctx, req.UserId, req.Rule, req.OrganizationId, req.ProjectId, req.RuleCategory)
	return &pb.UpdateProjectRuleReply{}, err
}

func (r *RuleDispatcher) GetProjectRule(ctx context.Context, req *pb.GetProjectRuleRequest) (*pb.GetProjectRuleReply, error) {
	if req.UserId == 0 || req.OrganizationId == 0 {
		return nil, error2.ParamError
	}
	return application.SvrRule.GetUserProjectRule(ctx, req.UserId, req.OrganizationId)
}

func NewRuleService() *RuleDispatcher {
	return &RuleDispatcher{}
}
