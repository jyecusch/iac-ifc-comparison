package deploy

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"

	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func md5Hash(b []byte) string {
	hasher := md5.New() //#nosec G401 -- md5 used only to produce a unique ID from non-sensistive information (policy IDs)
	hasher.Write(b)

	return hex.EncodeToString(hasher.Sum(nil))
}

func (a *AwsExtendedProvider) Policy(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Policy) error {
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	topicConfig := &deploymentspb.Policy{
		Principals: config.Principals,
	}

	remainingConfig := &deploymentspb.Policy{
		Principals: config.Principals,
	}

	for _, res := range config.Resources {
		if res.Id.Type == resourcespb.ResourceType_Topic {
			topicConfig.Resources = append(topicConfig.Resources, res)
		} else {
			remainingConfig.Resources = append(remainingConfig.Resources, res)
		}
	}

	for _, action := range config.Actions {
		if action == resourcespb.Action_TopicPublish {
			topicConfig.Actions = append(topicConfig.Actions, action)
		} else {
			remainingConfig.Actions = append(remainingConfig.Actions, action)
		}
	}

	if len(topicConfig.Actions) > 0 {
		targetArns := make([]interface{}, 0, len(topicConfig.Resources))

		for _, res := range topicConfig.Resources {
			targetArns = append(targetArns, a.Topics[res.Id.Name].bus.Arn)
		}

		principalRoles := make(map[string]*iam.Role)

		for _, princ := range config.Principals {
			if role, err := a.roleForPrincipal(princ); err == nil {
				principalRoles[princ.Id.Name] = role
			} else {
				return err
			}
		}

		serialPolicy, err := json.Marshal(config)
		if err != nil {
			return err
		}

		policyJson := pulumi.All(targetArns...).ApplyT(func(args []interface{}) (string, error) {
			arns := make([]string, 0, len(args))

			for _, iArn := range args {
				arn, ok := iArn.(string)
				if !ok {
					return "", fmt.Errorf("input not a string: %T %v", arn, arn)
				}

				arns = append(arns, arn)
			}

			jsonb, err := json.Marshal(map[string]interface{}{
				"Version": "2012-10-17",
				"Statement": []map[string]interface{}{
					{
						"Action":   []string{"events:PutEvents"},
						"Effect":   "Allow",
						"Resource": arns,
					},
				},
			})
			if err != nil {
				return "", err
			}

			return string(jsonb), nil
		})

		// create role policy for each role
		for k, r := range principalRoles {
			// Role policies require a unique name
			// Use a hash of the policy document to help create a unique name
			policyName := fmt.Sprintf("%s-%s", k, md5Hash(serialPolicy))

			_, err := iam.NewRolePolicy(ctx, policyName, &iam.RolePolicyArgs{
				Role:   r.ID(),
				Policy: policyJson,
			}, opts...)
			if err != nil {
				return err
			}
		}
	}

	if len(remainingConfig.Actions) == 0 {
		return nil
	}

	return a.NitricAwsPulumiProvider.Policy(ctx, parent, name, remainingConfig)
}

func (a *AwsExtendedProvider) roleForPrincipal(resource *deploymentspb.Resource) (*iam.Role, error) {
	switch resource.Id.Type {
	case resourcespb.ResourceType_Service:
		if f, ok := a.LambdaRoles[resource.Id.Name]; ok {
			return f, nil
		}
	default:
		return nil, fmt.Errorf("could not find role for principal: %+v", resource)
	}

	return nil, fmt.Errorf("could not find role for principal: %+v", resource)
}
