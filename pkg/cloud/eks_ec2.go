package cloud

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type EKSSplitRoleEC2API struct {
	// https://quip-amazon.com/RtxJAub9BSlM/EKS-Tachyon-IAM-Technical-Approach
	describeAndDelete EC2API
	createAndMutate   EC2API
}

const assumeRoleSessionDuration = 1 * time.Hour

func NewEKSSplitRoleEC2API(region, roleForDescribeAndDelete, roleForCreateAndMutate string, awsSdkDebugLog bool) (*EKSSplitRoleEC2API, error) {

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		panic(err)
	}

	if awsSdkDebugLog {
		cfg.ClientLogMode = aws.LogRequestWithBody | aws.LogResponseWithBody
	}

	return &EKSSplitRoleEC2API{
		describeAndDelete: createEC2API(cfg, roleForDescribeAndDelete),
		createAndMutate:   createEC2API(cfg, roleForCreateAndMutate),
	}, nil
}

func createEC2API(cfg aws.Config, role string) EC2API {
	creds := stscreds.NewAssumeRoleProvider(sts.NewFromConfig(cfg), role, func(aro *stscreds.AssumeRoleOptions) {
		aro.Duration = assumeRoleSessionDuration
	})
	ec2Config := aws.Config{
		Region:       cfg.Region,
		DefaultsMode: aws.DefaultsModeStandard,
		Credentials:  aws.NewCredentialsCache(creds),
	}
	return ec2.NewFromConfig(ec2Config, func(o *ec2.Options) {
		o.APIOptions = append(o.APIOptions,
			RecordRequestsMiddleware(),
		)

		endpoint := os.Getenv("AWS_EC2_ENDPOINT")
		if endpoint != "" {
			o.BaseEndpoint = &endpoint
		}

		o.RetryMaxAttempts = retryMaxAttempt
	})
}

func (e *EKSSplitRoleEC2API) DescribeVolumes(ctx context.Context, params *ec2.DescribeVolumesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeVolumesOutput, error) {
	return e.describeAndDelete.DescribeVolumes(ctx, params, optFns...)
}
func (e *EKSSplitRoleEC2API) CreateVolume(ctx context.Context, params *ec2.CreateVolumeInput, optFns ...func(*ec2.Options)) (*ec2.CreateVolumeOutput, error) {
	return e.createAndMutate.CreateVolume(ctx, params, optFns...)
}
func (e *EKSSplitRoleEC2API) DeleteVolume(ctx context.Context, params *ec2.DeleteVolumeInput, optFns ...func(*ec2.Options)) (*ec2.DeleteVolumeOutput, error) {
	return e.describeAndDelete.DeleteVolume(ctx, params, optFns...)
}
func (e *EKSSplitRoleEC2API) AttachVolume(ctx context.Context, params *ec2.AttachVolumeInput, optFns ...func(*ec2.Options)) (*ec2.AttachVolumeOutput, error) {
	return e.createAndMutate.AttachVolume(ctx, params, optFns...)
}
func (e *EKSSplitRoleEC2API) DetachVolume(ctx context.Context, params *ec2.DetachVolumeInput, optFns ...func(*ec2.Options)) (*ec2.DetachVolumeOutput, error) {
	return e.createAndMutate.DetachVolume(ctx, params, optFns...)
}
func (e *EKSSplitRoleEC2API) DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	return e.describeAndDelete.DescribeInstances(ctx, params, optFns...)
}
func (e *EKSSplitRoleEC2API) DescribeAvailabilityZones(ctx context.Context, params *ec2.DescribeAvailabilityZonesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeAvailabilityZonesOutput, error) {
	return e.describeAndDelete.DescribeAvailabilityZones(ctx, params, optFns...)
}
func (e *EKSSplitRoleEC2API) CreateSnapshot(ctx context.Context, params *ec2.CreateSnapshotInput, optFns ...func(*ec2.Options)) (*ec2.CreateSnapshotOutput, error) {
	return e.createAndMutate.CreateSnapshot(ctx, params, optFns...)
}
func (e *EKSSplitRoleEC2API) DeleteSnapshot(ctx context.Context, params *ec2.DeleteSnapshotInput, optFns ...func(*ec2.Options)) (*ec2.DeleteSnapshotOutput, error) {
	return e.describeAndDelete.DeleteSnapshot(ctx, params, optFns...)
}
func (e *EKSSplitRoleEC2API) DescribeSnapshots(ctx context.Context, params *ec2.DescribeSnapshotsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSnapshotsOutput, error) {
	return e.describeAndDelete.DescribeSnapshots(ctx, params, optFns...)
}
func (e *EKSSplitRoleEC2API) ModifyVolume(ctx context.Context, params *ec2.ModifyVolumeInput, optFns ...func(*ec2.Options)) (*ec2.ModifyVolumeOutput, error) {
	return e.createAndMutate.ModifyVolume(ctx, params, optFns...)
}
func (e *EKSSplitRoleEC2API) DescribeVolumesModifications(ctx context.Context, params *ec2.DescribeVolumesModificationsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeVolumesModificationsOutput, error) {
	return e.describeAndDelete.DescribeVolumesModifications(ctx, params, optFns...)
}
func (e *EKSSplitRoleEC2API) DescribeTags(ctx context.Context, params *ec2.DescribeTagsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeTagsOutput, error) {
	return e.describeAndDelete.DescribeTags(ctx, params, optFns...)
}
func (e *EKSSplitRoleEC2API) CreateTags(ctx context.Context, params *ec2.CreateTagsInput, optFns ...func(*ec2.Options)) (*ec2.CreateTagsOutput, error) {
	return e.createAndMutate.CreateTags(ctx, params, optFns...)
}
func (e *EKSSplitRoleEC2API) DeleteTags(ctx context.Context, params *ec2.DeleteTagsInput, optFns ...func(*ec2.Options)) (*ec2.DeleteTagsOutput, error) {
	return e.describeAndDelete.DeleteTags(ctx, params, optFns...)
}
func (e *EKSSplitRoleEC2API) EnableFastSnapshotRestores(ctx context.Context, params *ec2.EnableFastSnapshotRestoresInput, optFns ...func(*ec2.Options)) (*ec2.EnableFastSnapshotRestoresOutput, error) {
	return e.createAndMutate.EnableFastSnapshotRestores(ctx, params, optFns...)
}
