package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	//awsSDK "github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
)

func TestFargateCloudNuke(t *testing.T) {
	randomName := fmt.Sprintf("cloud-nuke-test-%s", strings.ToLower(random.UniqueId()))

	awsRegion := "eu-west-1"

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/terraform-fargate-cloud-nuke-basic/",

		Vars: map[string]interface{}{
			"name": randomName,
		},

		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	ValidateECR(t, terraformOptions, randomName)
}

func ValidateECR(t *testing.T, terraformOptions *terraform.Options, randomName string) {
	accountId := aws.GetAccountId(t)
	expectedEcrArn := fmt.Sprintf("arn:aws:ecr:eu-west-1:%s:repository/%s", accountId, randomName)
	ecrArn := terraform.Output(t, terraformOptions, "ecs_repository_arn")

	assert.Equal(t, expectedEcrArn, ecrArn)
}
