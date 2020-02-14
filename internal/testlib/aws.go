package testlib

import mocks "github.com/mattermost/mattermost-cloud/internal/mocks/aws-sdk"

// AWSMockedAPI ..
type AWSMockedAPI struct {
	ACM            *mocks.ACMAPI
	RDS            *mocks.RDSAPI
	IAM            *mocks.IAMAPI
	EC2            *mocks.EC2API
	S3             *mocks.S3API
	Route53        *mocks.Route53API
	SecretsManager *mocks.SecretsManagerAPI
}

// NewAWSMockedAPI ...
func NewAWSMockedAPI() *AWSMockedAPI {
	return &AWSMockedAPI{
		ACM:            &mocks.ACMAPI{},
		RDS:            &mocks.RDSAPI{},
		IAM:            &mocks.IAMAPI{},
		EC2:            &mocks.EC2API{},
		S3:             &mocks.S3API{},
		Route53:        &mocks.Route53API{},
		SecretsManager: &mocks.SecretsManagerAPI{},
	}
}