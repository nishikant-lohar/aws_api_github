package session

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"

	dep "aws_api/dependencies"
)

// GetStsToken func
func GetStsToken() *sts.Credentials {

	fmt.Println("Enter MFA code")
	var mfaToken string
	fmt.Scanf("%s", &mfaToken)
	sessionHandler := sts.New(session.New())
	inputCredentials := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(dep.AwsSessionDuration),
		SerialNumber:    aws.String(dep.AwsSerialArnNo),
		TokenCode:       aws.String(mfaToken),
	}

	result, err := sessionHandler.GetSessionToken(inputCredentials)
	if err != nil {
		if aerror, ok := err.(awserr.Error); ok {
			switch aerror.Code() {
			case sts.ErrCodeRegionDisabledException:
				fmt.Printf("ErrCodeRegionDisabledException: %s", aerror.Error())
			default:
				fmt.Printf("Unexpected Error 1: %s", aerror.Error())
			}
		} else {
			fmt.Printf("Unexpected Error 2: %s", err.Error())
		}
	}
	return result.Credentials
}

// GetStsTokenApi func
func GetStsTokenApi(respWriter http.ResponseWriter, request *http.Request) {
	// fmt.Println("in GetStsTokenApi")
	type responseStruct struct {
		Form    bool
		Success bool
		Message string
	}
	var response responseStruct

	tmpl := template.Must(template.ParseFiles("assets/generate_token.html"))
	if request.Method != http.MethodPost {
		response.Form = false
		tmpl.Execute(respWriter, nil)
		return
	} else if request.Method == http.MethodPost {
		response.Form = true
		// fmt.Println("mfa_token: ", request.FormValue("mfa_token"))

		mfaToken := request.FormValue("mfa_token")

		sessionHandler := sts.New(session.New())
		inputCredentials := &sts.GetSessionTokenInput{
			DurationSeconds: aws.Int64(dep.AwsSessionDuration),
			SerialNumber:    aws.String(dep.AwsSerialArnNo),
			TokenCode:       aws.String(mfaToken),
		}

		result, err := sessionHandler.GetSessionToken(inputCredentials)
		if err != nil {
			if aerror, ok := err.(awserr.Error); ok {
				switch aerror.Code() {
				case sts.ErrCodeRegionDisabledException:
					fmt.Printf("ErrCodeRegionDisabledException: %s", aerror.Error())
					response.Success = false
					response.Message = "ErrCodeRegionDisabledException"
					tmpl.Execute(respWriter, response)
				default:
					fmt.Printf("Unexpected Error 1: %s", aerror.Error())
					response.Success = false
					response.Message = "Unexpected Error 1"
					tmpl.Execute(respWriter, response)
				}
			} else {
				fmt.Printf("Unexpected Error 2: %s", err.Error())
				response.Success = false
				response.Message = "Unexpected Error 2"
				tmpl.Execute(respWriter, response)
			}
		} else {
			dep.SessionTokenHandler = result.Credentials
			response.Success = true
			response.Message = "Successfully Created Token"
			log.Println("Successfully generated token")
			tmpl.Execute(respWriter, response)
		}
	}
}
