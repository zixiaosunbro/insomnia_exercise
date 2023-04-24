package rule

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
	"insomnia/src/pkg/middleware_custom/mysql"
	rs "insomnia/src/services/rule"
)

type RepositoryRule struct {
}

func (r *RepositoryRule) GetProjectRule(ctx context.Context, equalFilter map[string]any, opts ...mysql.DBOption) ([]*ProjectLintingRule, error) {
	opt := mysql.NewDBOptions()
	opt.Apply(opts...)
	var session = rs.DBSlave
	if opt.UserMaster {
		session = rs.DBMaster
	}
	if opt.Tx != nil {
		session = opt.Tx
	}
	var projectRules []*ProjectLintingRule
	result := session.WithContext(ctx).Where(equalFilter).Find(&projectRules)
	return projectRules, result.Error
}

func (r *RepositoryRule) SaveProjectRule(ctx context.Context, projectRule *ProjectLintingRule, opts ...mysql.DBOption) error {
	opt := mysql.NewDBOptions()
	opt.Apply(opts...)
	var session *gorm.DB
	var err error
	if opt.Tx != nil {
		session = opt.Tx
	} else {
		session = rs.DBMaster.Begin()
		if session.Error != nil {
			return session.Error
		}
		defer func() {
			if err != nil {
				_ = session.Rollback()
			}
		}()
	}
	result := session.WithContext(ctx).Save(projectRule)
	if result.Error != nil {
		return result.Error
	}
	if opt.Tx == nil {
		if err = session.Commit().Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *RepositoryRule) GetLintingRule(ctx context.Context, equalFilter map[string]any, opts ...mysql.DBOption) ([]*LintingRule, error) {
	opt := mysql.NewDBOptions()
	opt.Apply(opts...)
	var session = rs.DBSlave
	if opt.UserMaster {
		session = rs.DBMaster
	}
	if opt.Tx != nil {
		session = opt.Tx
	}
	var rules []*LintingRule
	result := session.WithContext(ctx).Where(equalFilter).Find(&rules)
	return rules, result.Error
}

// SaveLintingRule support commit outside
func (r *RepositoryRule) SaveLintingRule(ctx context.Context, rule *LintingRule, opts ...mysql.DBOption) error {
	opt := mysql.NewDBOptions()
	opt.Apply(opts...)
	var session *gorm.DB
	var err error
	if opt.Tx != nil {
		session = opt.Tx
	} else {
		session = rs.DBMaster.Begin()
		if session.Error != nil {
			return session.Error
		}
		defer func() {
			if err != nil {
				_ = session.Rollback()
			}
		}()
	}
	result := session.WithContext(ctx).Save(rule)
	if result.Error != nil {
		return result.Error
	}
	if opt.Tx == nil {
		if err = session.Commit().Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *RepositoryRule) GetProjectRuleCache(ctx context.Context, organId, projectId int32) (*ProjectRuleCache, error) {
	ruleCache := &ProjectRuleCache{
		OrganizationId: organId,
		ProjectId:      projectId,
	}
	info, err := rs.RedisCli.Get(ctx, ruleCache.CacheTableName()).Result()
	if err != nil {
		return nil, err
	}
	if err = jsoniter.UnmarshalFromString(info, ruleCache); err != nil {
		return nil, err
	}
	return ruleCache, nil
}

// SaveProjectRuleCache using redis cache, avoid too many db queries
func (r *RepositoryRule) SaveProjectRuleCache(ctx context.Context, organId, projectId int32, category int, rule string) error {
	cacheIns := &ProjectRuleCache{
		OrganizationId: organId,
		ProjectId:      projectId,
		Category:       category,
		Rule:           rule,
	}
	info, err := jsoniter.MarshalToString(cacheIns)
	if err != nil {
		return err
	}
	if _, err = rs.RedisCli.Set(ctx, cacheIns.CacheTableName(), info, cacheIns.ExpireTime()).Result(); err != nil {
		return err
	}

	return nil
}
