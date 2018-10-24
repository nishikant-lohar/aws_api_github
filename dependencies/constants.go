package dependencies

import "github.com/aws/aws-sdk-go/service/sts"

const (
	// AwsSessionDuration time for the session
	AwsSessionDuration = 3600

	// AwsSerialArnNo MFA token
	AwsSerialArnNo = "arn:aws:iam::294069028655:mfa/nishikant.lohar"

	// MfaTokenCode from authenticator app
	MfaTokenCode = "235500"
	// InstanceName var
	InstanceName = "N_MyInstance"

	// ImageID var
	ImageID = "ami-0b33d91d"

	// MinInstanceCount var
	MinInstanceCount = 1

	// MaxInstanceCount var
	MaxInstanceCount = 1

	// Region var
	Region = "eu-west-1"

	// InstanceType var
	InstanceType = "t2.micro"
)

// SessionTokenHandler this will contain the session token handler which will be initialized once the init func gets executed
var SessionTokenHandler *sts.Credentials
