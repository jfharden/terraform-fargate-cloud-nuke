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
	"github.com/stretchr/testify/assert"
)

func TestFargateCloudNuke(t *testing.T) {
	awsRegion := "eu-west-1"

	uniqueId := strings.ToLower(random.UniqueId())

	awsSession, err := session.NewSession(
		&awsSdk.Config{
			Region: awsSdk.String(awsRegion),
		},
	)

	if err != nil {
		t.Errorf("Couldn't create AWS session %s", err)
	}

	randomName := fmt.Sprintf("cloud-nuke-test-%s", uniqueId)

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/terraform-fargate-cloud-nuke-basic/",

		Vars: map[string]interface{}{
			"name": randomName,
			"tags": map[string]string{
				"TerraformTest": randomName,
				"Stage":         "test",
			},
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
	// Check the ARN is as expected
	accountId := aws.GetAccountId(t)
	expectedEcrArn := fmt.Sprintf("arn:aws:ecr:eu-west-1:%s:repository/%s", accountId, randomName)
	ecrArn := terraform.Output(t, terraformOptions, "ecs_repository_arn")

	assert.Equal(t, expectedEcrArn, ecrArn)

	// Check the repo has image scanning enabled
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

	assert.Equal(t, 1, len(result.Repositories), "Exactly 1 ECR repository")

	repository := result.Repositories[0]

	assert.True(t,
		awsSdk.BoolValue(repository.ImageScanningConfiguration.ScanOnPush),
		"Image scanning enabled on ECR repo",
	)

	// Check the repo has the tags we set
	tagRequest, tagResponse := ecrService.ListTagsForResourceRequest(
		&ecr.ListTagsForResourceInput{
			ResourceArn: awsSdk.String(ecrArn),
		},
	)

	err = tagRequest.Send()
	if err != nil {
		t.Errorf("Couldn't request tags for ECR repository: %s", err)
	}

	expectedTags := map[string]string{
		"TerraformTest": randomName,
		"Stage":         "test",
	}

	AssertEcrTagsMatch(t, expectedTags, tagResponse.Tags)
}

func AssertEcrTagsMatch(t *testing.T, expected map[string]string, actual []*ecr.Tag) {
	equalNumberOfTags := assert.Equal(t,
		len(expected),
		len(actual),
		fmt.Sprintf(
			"Number of tags expected [%v] does match actual number of tags [%v]. Expected: %v Actual: %v",
			len(expected),
			len(actual),
			expected,
			actual,
		),
	)

	if !equalNumberOfTags {
		return
	}

	for _, tag := range actual {
		actualKey := awsSdk.StringValue(tag.Key)
		actualValue := awsSdk.StringValue(tag.Value)

		expectedValue, tagInExpected := expected[actualKey]
		if !tagInExpected {
			t.Errorf("Error: Tag key %s in actual, but not expected\n", actualKey)
			return
		}

		if !(expectedValue == actualValue) {
			t.Errorf("Error: Tag with key %s has value %s but we expected %s", actualKey, actualValue, expectedValue)
			return
		}
	}
}
