package application

import (
	"bou.ke/monkey"
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"insomnia/src/services/rule/domain/rule"
	"insomnia/src/services/rule/thirdparty/services"
	"reflect"
	"testing"
)

func TestCheckRuleValid(t *testing.T) {
	// case1: valid json string
	jsonData := map[string]interface{}{
		"age":  26,
		"name": "wei",
	}
	jsonStr, _ := jsoniter.MarshalToString(jsonData)
	valid := CheckRuleValid(jsonStr, 1)
	assert.True(t, valid)
	// case2: invalid json string
	jsonStr = "invalid json string"
	valid = CheckRuleValid(jsonStr, 1)
	assert.False(t, valid)
	// case3: valid yaml string
	yamlStr := `
name: John Doe
age: 30
address:
  street: Main St
  city: Anytown
  state: CA
`
	valid = CheckRuleValid(yamlStr, 2)
	assert.True(t, valid)
	yamlStr = `
		name: John Doe
		age: 30
	`
	valid = CheckRuleValid(yamlStr, 2)
	assert.False(t, valid)
}

func TestAssembleProjectRule(t *testing.T) {
	projects := []*services.ProjectInfo{
		{
			Id:   1,
			Name: "project1",
		},
		{
			Id:   2,
			Name: "project2",
		},
		{
			Id:   3,
			Name: "project3",
		},
	}
	projectRule := []*rule.ProjectRuleCache{
		{
			ProjectId:      1,
			OrganizationId: 1,
			Rule:           "rule1",
			Category:       1,
		},
		{
			ProjectId:      2,
			OrganizationId: 1,
			Rule:           "rule2",
			Category:       2,
		},
	}
	results := AssembleProjectRule(projects, projectRule)
	assert.NotNil(t, results)
	assert.Len(t, results, 2)
}

func TestMGetProjectInfo(t *testing.T) {
	organId, projectIds := int32(1), []int32{2}
	var ruleSvr *rule.ServiceRule
	monkey.PatchInstanceMethod(reflect.TypeOf(ruleSvr), "GetProjectInfo", func(_ *rule.ServiceRule, ctx context.Context, organId, projectId int32) (*rule.ProjectRuleCache, error) {
		return &rule.ProjectRuleCache{
			OrganizationId: 1,
			ProjectId:      2,
			Category:       1,
			Rule:           "test",
		}, nil
	})
	defer monkey.UnpatchAll()
	result, err := MGetProjectInfo(context.Background(), organId, projectIds)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}
