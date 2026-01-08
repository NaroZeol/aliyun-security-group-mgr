package ecs

import (
	ecs "github.com/alibabacloud-go/ecs-20140526/v7/client"
)

type InnerDescribeSecurityGroupAttributeResponse ecs.DescribeSecurityGroupAttributeResponse

type Policy string

const (
	PolicyAccept = "Accept"
	PolicyDrop   = "Drop"

	DirectionIngress = "ingress"
	DirectionEgress  = "egress"
)

type SecurityGroupRule struct {
	Id         string
	CidrIp     string
	PortRange  string
	IpProtocol string

	Policy      string
	Priority    string
	Direction   string
	Description string
}

func buildSecurityGroupRules(response ecs.DescribeSecurityGroupAttributeResponse) ([]SecurityGroupRule, error) {
	rules := []SecurityGroupRule{}

	if response.Body == nil || response.Body.Permissions == nil || response.Body.Permissions.Permission == nil {
		return rules, nil
	}

	for _, perm := range response.Body.Permissions.Permission {
		rule := SecurityGroupRule{
			Id:          *perm.SecurityGroupRuleId,
			Policy:      *perm.Policy,
			Priority:    *perm.Priority,
			Description: *perm.Description,

			CidrIp:     *perm.SourceCidrIp,
			PortRange:  *perm.PortRange,
			IpProtocol: *perm.IpProtocol,
			Direction:  *perm.Direction,
		}
		rules = append(rules, rule)
	}

	return rules, nil
}
