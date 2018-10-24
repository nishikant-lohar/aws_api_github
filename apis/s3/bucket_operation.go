package s3

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	dep "aws_api/dependencies"
)

// CreateBucket func
func CreateBucket(respWriter http.ResponseWriter, request *http.Request) {
	print("In create bucket")
	session, errInSessionCreation := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(*dep.SessionTokenHandler.AccessKeyId, *dep.SessionTokenHandler.SecretAccessKey, *dep.SessionTokenHandler.SessionToken),
		Region:      aws.String(dep.Region),
	})

	fmt.Println("session: ", session)

	// tmpl.Execute(respWriter, struct{ Success bool }{true})

	if errInSessionCreation != nil {
		// response.Info.Success = false
		// response.Info.StatusCode = 400
		// response.Info.Message = fmt.Sprintf("Error while creating instance handler: %s", errInSessionCreation)

		fmt.Println("Some error while creatung instance handler: ")
		// if err := json.NewEncoder(respWriter).Encode(response); err != nil {
		// 	panic(err)
		// }
		// tmpl.Execute(respWriter, response)
		return
	}
	fmt.Println("instnace handler")
	instanceHandler := s3.New(session)
	fmt.Println("Bucket creation")
	bucketCreationResult, bucketCreationErr := instanceHandler.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String("testgolangbucketnishi"),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(dep.Region),
		},
	})
	fmt.Println("bucker result: ", bucketCreationResult)
	fmt.Println("Bucket error", bucketCreationErr)

	if bucketCreationErr != nil {
		if aerr, ok := bucketCreationErr.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				fmt.Println(s3.ErrCodeBucketAlreadyExists, aerr.Error())
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				fmt.Println(s3.ErrCodeBucketAlreadyOwnedByYou, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(bucketCreationErr.Error())
		}
		return
	}
}
