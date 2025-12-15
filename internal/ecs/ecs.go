package ecs

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ecs "github.com/alibabacloud-go/ecs-20140526/v7/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	credential "github.com/aliyun/credentials-go/credentials"

	"aliyun-security-group-mgr/internal/conf"
)

type Clerk struct {
	ecsClient *ecs.Client
}

func NewEcsClient(globalConfig *conf.GlobalConfiguration) (*Clerk, error) {
	client, err := createClient(globalConfig)
	if err != nil {
		return nil, err
	}
	return &Clerk{
		ecsClient: client,
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

func (e *Clerk) DescribeSecurityGroupAttribute(globalConfig *conf.GlobalConfiguration) (*ecs.DescribeSecurityGroupAttributeResponse, error) {
	describeSecurityGroupAttributeRequest := &ecs.DescribeSecurityGroupAttributeRequest{
		RegionId:        globalConfig.ECS.RegionId,
		SecurityGroupId: globalConfig.SecurityGroup.Id,
		NicType:         tea.String("internet"),
	}
	runtime := &util.RuntimeOptions{}
	return e.ecsClient.DescribeSecurityGroupAttributeWithOptions(describeSecurityGroupAttributeRequest, runtime)
}
