package iam

import (
	"testing"

	defsecTypes "github.com/aquasecurity/defsec/pkg/types"

	"github.com/aquasecurity/defsec/pkg/providers/aws/iam"
	"github.com/liamg/iamgo"

	"github.com/aquasecurity/defsec/internal/adapters/terraform/tftestutil"
	"github.com/aquasecurity/defsec/test/testutil"
)

func Test_adaptGroups(t *testing.T) {
	tests := []struct {
		name      string
		terraform string
		expected  []iam.Group
	}{
		{
			name: "basic",
			terraform: `
			resource "aws_iam_group_policy" "my_developer_policy" {
				name  = "my_developer_policy"
				group = aws_iam_group.my_developers.name

				policy = <<EOF
				{
				  "Version": "2012-10-17",
				  "Statement": [
				  {
					"Sid": "new policy",
					"Effect": "Allow",
					"Resource": "*",
					"Action": [
						"ec2:Describe*"
					]
				  }
				  ]
				}
				EOF
			  }
			  
			  resource "aws_iam_group" "my_developers" {
				name = "developers"
				path = "/users/"
			  }
			  
			  `,
			expected: []iam.Group{
				{
					Metadata: defsecTypes.NewTestMetadata(),
					Name:     defsecTypes.String("developers", defsecTypes.NewTestMetadata()),
					Policies: []iam.Policy{
						{
							Metadata: defsecTypes.NewTestMetadata(),
							Name:     defsecTypes.String("my_developer_policy", defsecTypes.NewTestMetadata()),
							Document: func() iam.Document {

								builder := iamgo.NewPolicyBuilder()
								builder.WithVersion("2012-10-17")

								sb := iamgo.NewStatementBuilder()
								sb.WithEffect(iamgo.EffectAllow)
								sb.WithSid("new policy")
								sb.WithActions([]string{"ec2:Describe*"})
								sb.WithResources([]string{"*"})

								builder.WithStatement(sb.Build())

								return iam.Document{
									Parsed:   builder.Build(),
									Metadata: defsecTypes.NewTestMetadata(),
									IsOffset: false,
									HasRefs:  false,
								}
							}(),
							Builtin: defsecTypes.Bool(false, defsecTypes.NewTestMetadata()),
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			modules := tftestutil.CreateModulesFromSource(t, test.terraform, ".tf")
			adapted := adaptGroups(modules)
			testutil.AssertDefsecEqual(t, test.expected, adapted)
		})
	}
}
