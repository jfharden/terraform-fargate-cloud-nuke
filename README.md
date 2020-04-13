# terraform-fargate-cloud-nuke

This terraform module provides cloud-nuke https://github.com/gruntwork-io/cloud-nuke running on a schedule in AWS  
fargate.

WARNING: This is extremely destructive and will wipe out most resources in the aws account in which it runs. It  
also creates a role with admin permissions, so you should manage access to that role with extreme caution.

The following things will be created:

1. An ECR repo to store the docker image  
2. A Codebuild to build the docker image, run tests, and push to the ECR repo  
3. An ECS Fargate task to run the docker container with associated task role and security group  
4. Optionally an ECS cluster to run the task in (if none is specified then one will be created)  
5. Optionally a cloudwatch schedule to rebuild the docker container

See https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/ScheduledEvents.html for the specification for the  
cloudwatch events.

## Example usage

```
module "cloud-nuke" {
  source = "git::https://github.com/jfharden/terraform-fargate-cloud-nukevpc.git?ref=v1.0.0"
  name = "cloud-nuke"
  cloudwatch_schedule = "cron(0 3 * * ? *)" // 3am UTC every day
  cloud_nuke_version = "0.1.17"
}

```

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 0.12 |

## Providers

| Name | Version |
|------|---------|
| aws | n/a |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| cloud\_nuke\_version | Version number of cloud nuke to install. | `string` | `"v0.1.17"` | no |
| cloudwatch\_build\_schedule | Cloudwatch Schedule expression to specify when to run a build of the container. If blank then no build schdule will be created. | `string` | `""` | no |
| cloudwatch\_run\_schedule | Cloudwatch Schedule expression to specify when to run the cloud nuke task. | `string` | `"cron(0 3 * * ? *)"` | no |
| name | Name to apply to all the resources. | `string` | `"fargate-cloud-nuke"` | no |

## Outputs

| Name | Description |
|------|-------------|
| cloudwatch\_schedule\_arn | ARN for the cloudwatch schedule, if it was created |
| ecs\_repository\_arn | ARN for the ecs repository |
| task\_arn | ARN for the ecs task |

