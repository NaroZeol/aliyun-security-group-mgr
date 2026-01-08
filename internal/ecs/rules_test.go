package ecs

import (
	"reflect"
	"testing"

	ecs "github.com/alibabacloud-go/ecs-20140526/v7/client"
	"github.com/alibabacloud-go/tea/tea"
)

func Test_buildSecurityGroupRules(t *testing.T) {
	type args struct {
		response ecs.DescribeSecurityGroupAttributeResponse
	}
	tests := []struct {
		name    string
		args    args
		want    []SecurityGroupRule
		wantErr bool
	}{
		{
			name: "empty response",
			args: args{
				response: ecs.DescribeSecurityGroupAttributeResponse{},
			},
			want:    []SecurityGroupRule{},
			wantErr: false,
		},
		{
			name: "nil body",
			args: args{
				response: ecs.DescribeSecurityGroupAttributeResponse{
					Body: nil,
				},
			},
			want:    []SecurityGroupRule{},
			wantErr: false,
		},
		{
			name: "valid permissions",
			args: args{
				response: ecs.DescribeSecurityGroupAttributeResponse{
					Body: &ecs.DescribeSecurityGroupAttributeResponseBody{
						Permissions: &ecs.DescribeSecurityGroupAttributeResponseBodyPermissions{
							Permission: []*ecs.DescribeSecurityGroupAttributeResponseBodyPermissionsPermission{
								{
									SecurityGroupRuleId: tea.String("sg-rule-1"),
									Policy:              tea.String("Accept"),
									Priority:            tea.String("100"),
									Description:         tea.String("test rule"),
									SourceCidrIp:        tea.String("0.0.0.0/0"),
									PortRange:           tea.String("80/80"),
									IpProtocol:          tea.String("TCP"),
									Direction:           tea.String("ingress"),
								},
								{
									SecurityGroupRuleId: tea.String("sg-rule-2"),
									Policy:              tea.String("Drop"),
									Priority:            tea.String("200"),
									Description:         tea.String(""),
									SourceCidrIp:        tea.String("10.0.0.0/8"),
									PortRange:           tea.String("22/22"),
									IpProtocol:          tea.String("UDP"),
									Direction:           tea.String("egress"),
								},
							},
						},
					},
				},
			},
			want: []SecurityGroupRule{
				{
					Id:          "sg-rule-1",
					Policy:      "Accept",
					Priority:    "100",
					Description: "test rule",
					CidrIp:      "0.0.0.0/0",
					PortRange:   "80/80",
					IpProtocol:  "TCP",
					Direction:   "ingress",
				},
				{
					Id:          "sg-rule-2",
					Policy:      "Drop",
					Priority:    "200",
					Description: "",
					CidrIp:      "10.0.0.0/8",
					PortRange:   "22/22",
					IpProtocol:  "UDP",
					Direction:   "egress",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildSecurityGroupRules(tt.args.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildSecurityGroupRules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildSecurityGroupRules() = %v, want %v", got, tt.want)
			}
		})
	}
}
