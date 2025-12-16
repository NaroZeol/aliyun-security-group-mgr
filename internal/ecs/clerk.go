package ecs

import (
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ecs "github.com/alibabacloud-go/ecs-20140526/v7/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	credential "github.com/aliyun/credentials-go/credentials"

	"aliyun-security-group-mgr/internal/conf"
)

type Clerk struct {
	ecsClient *ecs.Client
	config    *conf.GlobalConfiguration
}

func NewClerk(globalConfig *conf.GlobalConfiguration) (*Clerk, error) {
	client, err := createClient(globalConfig)
	if err != nil {
		return nil, err
	}
	return &Clerk{
		ecsClient: client,
		config:    globalConfig,
	}, nil
}

func createClient(globalConfig *conf.GlobalConfiguration) (*ecs.Client, error) {
	credentialConfig := &credential.Config{
		Type:            tea.String(*globalConfig.Credential.Type),
		AccessKeyId:     tea.String(*globalConfig.Credential.AccessKeyId),
		AccessKeySecret: tea.String(*globalConfig.Credential.AccessKeySecret),
	}
	credential, err := credential.NewCredential(credentialConfig)
	if err != nil {
		return nil, err
	}

	config := &openapi.Config{
		Credential: credential,
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Ecs
	config.Endpoint = tea.String(*globalConfig.ECS.Endpoint)

	return ecs.NewClient(config)
}

func (e *Clerk) DescribeSecurityGroupAttribute() ([]SecurityGroupRule, error) {
	describeSecurityGroupAttributeRequest := &ecs.DescribeSecurityGroupAttributeRequest{
		RegionId:        e.config.ECS.RegionId,
		SecurityGroupId: e.config.SecurityGroup.Id,
	}
	runtime := &util.RuntimeOptions{}
	response, err := e.ecsClient.DescribeSecurityGroupAttributeWithOptions(describeSecurityGroupAttributeRequest, runtime)
	if err != nil {
		return nil, err
	}
	return buildSecurityGroupRules(*response)
}

// GetIpRules gets all rules for the given cidrIp
func (e *Clerk) GetIpRules(cidrIp string) ([]SecurityGroupRule, error) {
	rules, err := e.DescribeSecurityGroupAttribute()
	if err != nil {
		return nil, err
	}

	var filteredRules []SecurityGroupRule
	for _, rule := range rules {
		if rule.CidrIp == cidrIp {
			filteredRules = append(filteredRules, rule)
		}
	}

	return filteredRules, nil
}

func (e *Clerk) AddSecurityGroupRule(rule SecurityGroupRule) error {
	if rule.Policy != PolicyAccept {
		return e.addIngressSecurityGroupRule(rule)
	}

	if rule.Policy == PolicyAccept {
		return e.addEgressSecurityGroupRule(rule)
	}

	return fmt.Errorf("unsupported policy: %s for rule: %v", rule.Policy, rule)
}

func (e *Clerk) addIngressSecurityGroupRule(rule SecurityGroupRule) error {
	authorizeSecurityGroupRequest := &ecs.AuthorizeSecurityGroupRequest{
		RegionId:        e.config.ECS.RegionId,
		SecurityGroupId: e.config.SecurityGroup.Id,

		IpProtocol:   &rule.IpProtocol,
		PortRange:    &rule.PortRange,
		SourceCidrIp: &rule.CidrIp,
		Description:  &rule.Description,
		Priority:     tea.String("1"),
	}
	runtime := &util.RuntimeOptions{}
	_, err := e.ecsClient.AuthorizeSecurityGroupWithOptions(authorizeSecurityGroupRequest, runtime)
	if err != nil {
		return err
	}
	return nil
}

func (e *Clerk) addEgressSecurityGroupRule(rule SecurityGroupRule) error {
	authorizeSecurityGroupEgressRequest := &ecs.AuthorizeSecurityGroupEgressRequest{
		RegionId:        e.config.ECS.RegionId,
		SecurityGroupId: e.config.SecurityGroup.Id,

		IpProtocol:  &rule.IpProtocol,
		PortRange:   &rule.PortRange,
		DestCidrIp:  &rule.CidrIp,
		Description: &rule.Description,
		Priority:    tea.String("1"),
	}

	runtime := &util.RuntimeOptions{}
	_, err := e.ecsClient.AuthorizeSecurityGroupEgressWithOptions(authorizeSecurityGroupEgressRequest, runtime)
	if err != nil {
		return err
	}
	return nil
}

func (e *Clerk) RemoveSecurityGroupRule(rule SecurityGroupRule) error {
	if rule.Policy != PolicyAccept {
		return e.removeIngressSecurityGroupRule(rule)
	}

	if rule.Policy == PolicyAccept {
		return e.removeEgressSecurityGroupRule(rule)
	}

	return fmt.Errorf("unsupported policy: %s for rule: %v", rule.Policy, rule)
}

func (e *Clerk) removeIngressSecurityGroupRule(rule SecurityGroupRule) error {
	revokeSecurityGroupRequest := &ecs.RevokeSecurityGroupRequest{
		RegionId:        e.config.ECS.RegionId,
		SecurityGroupId: e.config.SecurityGroup.Id,

		SecurityGroupRuleId: []*string{tea.String(rule.Id)},
	}
	runtime := &util.RuntimeOptions{}
	_, err := e.ecsClient.RevokeSecurityGroupWithOptions(revokeSecurityGroupRequest, runtime)
	if err != nil {
		return err
	}
	return nil
}

func (e *Clerk) removeEgressSecurityGroupRule(rule SecurityGroupRule) error {
	revokeSecurityGroupEgressRequest := &ecs.RevokeSecurityGroupEgressRequest{
		RegionId:        e.config.ECS.RegionId,
		SecurityGroupId: e.config.SecurityGroup.Id,

		SecurityGroupRuleId: []*string{tea.String(rule.Id)},
	}

	runtime := &util.RuntimeOptions{}
	_, err := e.ecsClient.RevokeSecurityGroupEgressWithOptions(revokeSecurityGroupEgressRequest, runtime)
	if err != nil {
		return err
	}
	return nil
}

func (e *Clerk) ModifySecurityGroupRule(ruleId string, newRule SecurityGroupRule) error {
	if newRule.Policy != PolicyAccept {
		return e.modifyIngressSecurityRule(ruleId, newRule)
	}

	if newRule.Policy == PolicyAccept {
		return e.modifyEgressSecurityRule(ruleId, newRule)
	}

	return fmt.Errorf("unsupported policy: %s for rule: %v", newRule.Policy, newRule)
}

func (e *Clerk) modifyIngressSecurityRule(ruleId string, newRule SecurityGroupRule) error {
	modifySecurityGroupRuleRequest := &ecs.ModifySecurityGroupRuleRequest{
		RegionId:        e.config.ECS.RegionId,
		SecurityGroupId: e.config.SecurityGroup.Id,

		SecurityGroupRuleId: tea.String(ruleId),
		IpProtocol:          &newRule.IpProtocol,
		PortRange:           &newRule.PortRange,
		Description:         &newRule.Description,
		Priority:            &newRule.Priority,
	}

	runtime := &util.RuntimeOptions{}
	_, err := e.ecsClient.ModifySecurityGroupRuleWithOptions(modifySecurityGroupRuleRequest, runtime)
	if err != nil {
		return err
	}

	return nil
}

func (e *Clerk) modifyEgressSecurityRule(ruleId string, newRule SecurityGroupRule) error {
	modifySecurityGroupEgressRuleRequest := &ecs.ModifySecurityGroupEgressRuleRequest{
		RegionId:        e.config.ECS.RegionId,
		SecurityGroupId: e.config.SecurityGroup.Id,

		SecurityGroupRuleId: tea.String(ruleId),
		IpProtocol:          &newRule.IpProtocol,
		PortRange:           &newRule.PortRange,
		Description:         &newRule.Description,
		Priority:            &newRule.Priority,
	}

	runtime := &util.RuntimeOptions{}
	_, err := e.ecsClient.ModifySecurityGroupEgressRuleWithOptions(modifySecurityGroupEgressRuleRequest, runtime)
	if err != nil {
		return err
	}

	return nil
}
