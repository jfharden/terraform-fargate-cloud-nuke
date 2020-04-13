package main

import (
	"fmt"
	"strings"
	"testing"

	awsSdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	//awsSDK "github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
)

func TestFargateCloudNuke(t *testing.T) {
	awsRegion := "eu-west-1"

	awsSession, err := session.NewSession(
		&awsSdk.Config{
			Region: awsSdk.String(awsRegion),
		},
	)

	if err != nil {
		t.Errorf("Couldn't create AWS session %s", err)
	}

	randomName := fmt.Sprintf("cloud-nuke-test-%s", strings.ToLower(random.UniqueId()))

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

	ValidateECR(t, terraformOptions, randomName, awsSession)
}

func ValidateECR(t *testing.T, terraformOptions *terraform.Options, randomName string, awsSession *session.Session) {
	accountId := aws.GetAccountId(t)
	expectedEcrArn := fmt.Sprintf("arn:aws:ecr:eu-west-1:%s:repository/%s", accountId, randomName)
	ecrArn := terraform.Output(t, terraformOptions, "ecs_repository_arn")

	assert.Equal(t, expectedEcrArn, ecrArn)

	ecrService := ecr.New(awsSession)
	result, err := ecrService.DescribeRepositories(
		&ecr.DescribeRepositoriesInput{
			RepositoryNames: []*string{
				awsSdk.String(randomName),
			},
		},
	)

	if err != nil {
		t.Errorf("Error trying to describe ECR repository %s. Error was %s", randomName, err)
	}

	assert.Equal(t, len(result.Repositories), 1, "Exactly 1 ECR repository")

	repository := result.Repositories[0]

	assert.True(t, awsSdk.BoolValue(repository.ImageScanningConfiguration.ScanOnPush), "Image scanning enabled on ECR repo")
}
