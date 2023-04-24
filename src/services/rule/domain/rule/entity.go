package rule

import (
	"fmt"
	"time"
)

const (
	JsonRuleCategory = 1
	YamlRuleCategory = 2
)

type LintingRule struct {
	ID             int64     `gorm:"column:id;primary_key"`
	OrganizationId int32     `gorm:"column:organization_id"`
	ProjectId      int32     `gorm:"column:project_id"`
	CreatorUid     int64     `gorm:"column:creator_uid"`
	Rule           string    `gorm:"column:rule"`
	RuleCategory   int       `gorm:"column:rule_category"`
	CreateAt       time.Time `gorm:"column:create_at;autoCreateTime"`
}

func (l LintingRule) TableName() string {
	return "insomnia_project_linting_rule"
}

func CreateLintingRule(organId, projectId int32, userId int64, rule string, ruleCategory int) *LintingRule {
	return &LintingRule{
		OrganizationId: organId,
		ProjectId:      projectId,
		CreatorUid:     userId,
		Rule:           rule,
		RuleCategory:   ruleCategory,
	}
}

type ProjectLintingRule struct {
	ID             int64     `gorm:"column:id;primary_key"`
	OrganizationId int32     `gorm:"column:organization_id"`
	ProjectId      int32     `gorm:"column:project_id"`
	RuleId         int64     `gorm:"column:rule_id"`
	OperatorUid    int64     `gorm:"column:opt_uid"`
	CreateAt       time.Time `gorm:"column:create_at;autoCreateTime"`
	UpdateAt       time.Time `gorm:"column:update_at;autoUpdateTime"`
}

func (p ProjectLintingRule) TableName() string {
	return "insomnia_project_apply_rule"
}

func CreateProjectLintingRule(organId, projectId int32, ruleId, userId int64) *ProjectLintingRule {
	return &ProjectLintingRule{
		OrganizationId: organId,
		ProjectId:      projectId,
		RuleId:         ruleId,
		OperatorUid:    userId,
	}
}

type ProjectRuleCache struct {
	OrganizationId int32
	ProjectId      int32
	Category       int
	Rule           string
}

func (p ProjectRuleCache) ExpireTime() time.Duration {
	return 3 * time.Minute
}

func (p *ProjectRuleCache) CacheTableName() string {
	return fmt.Sprintf("project_rule_%d_%d", p.OrganizationId, p.ProjectId)
}
