package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	awsSdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

func Setup(t *testing.T, workingDir string) *terraform.Options {
	randomName := fmt.Sprintf("cloud-nuke-test-%s", strings.ToLower(random.UniqueId()))
	awsRegion := "eu-west-1"

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

	test_structure.SaveString(t, workingDir, "randomName", randomName)
	test_structure.SaveString(t, workingDir, "awsRegion", awsRegion)
	test_structure.SaveTerraformOptions(t, workingDir, terraformOptions)

	return terraformOptions
}

func TestFargateCloudNuke(t *testing.T) {
	workingDir := filepath.Join(".terratest-working-dir", t.Name())

	defer test_structure.RunTestStage(t, "destroy", func() {
		options := test_structure.LoadTerraformOptions(t, workingDir)
		terraform.Destroy(t, options)
	})

	test_structure.RunTestStage(t, "init_apply", func() {
		options := Setup(t, workingDir)

		terraform.InitAndApply(t, options)
	})

	test_structure.RunTestStage(t, "validate", func() {
		options := test_structure.LoadTerraformOptions(t, workingDir)
		ValidateECR(t, options, workingDir)
	})
}

func ValidateECR(t *testing.T, terraformOptions *terraform.Options, workingDir string) {
	randomName := test_structure.LoadString(t, workingDir, "randomName")
	awsRegion := test_structure.LoadString(t, workingDir, "awsRegion")

	awsSession, err := session.NewSession(
		&awsSdk.Config{
			Region: awsSdk.String(awsRegion),
		},
	)

	if err != nil {
		t.Errorf("Couldn't create AWS session %s", err)
	}

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
		"TerraformTest": terraformOptions.Vars["name"].(string),
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
