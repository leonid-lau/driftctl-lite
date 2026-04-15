package cloud

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ResourceAttributes holds key-value pairs describing a live cloud resource.
type ResourceAttributes map[string]string

// Provider defines the interface for fetching live cloud resources.
type Provider interface {
	FetchResource(ctx context.Context, resourceType, resourceID string) (ResourceAttributes, error)
}

// AWSProvider implements Provider using the AWS SDK.
type AWSProvider struct {
	s3Client *s3.Client
	region   string
}

// NewAWSProvider creates a new AWSProvider using the default AWS config chain.
func NewAWSProvider(ctx context.Context, region string) (*AWSProvider, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("loading AWS config: %w", err)
	}
	return &AWSProvider{
		s3Client: s3.NewFromConfig(cfg),
		region:   region,
	}, nil
}

// FetchResource retrieves live attributes for a given resource type and ID.
func (p *AWSProvider) FetchResource(ctx context.Context, resourceType, resourceID string) (ResourceAttributes, error) {
	switch resourceType {
	case "aws_s3_bucket":
		return p.fetchS3Bucket(ctx, resourceID)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

// fetchS3Bucket fetches live attributes for an S3 bucket.
func (p *AWSProvider) fetchS3Bucket(ctx context.Context, bucketName string) (ResourceAttributes, error) {
	out, err := p.s3Client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, fmt.Errorf("fetching S3 bucket %q: %w", bucketName, err)
	}
	attrs := ResourceAttributes{
		"bucket": bucketName,
		"region": string(out.LocationConstraint),
	}
	return attrs, nil
}
