package ec2

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gorilla/mux"

	dep "aws_api/dependencies"
)

// CreateInstance This function creates instance which gets the instance details from the
// global variable
func CreateInstance(respWriter http.ResponseWriter, request *http.Request) {
	tmpl := template.Must(template.ParseFiles("assets/create_instance.html"))
	if request.Method != http.MethodPost {
		// response.Form = false
		tmpl.Execute(respWriter, nil)
		return
	} else if request.Method == http.MethodPost {
		fmt.Println("In post method")
		var response CreateInstanceResponse
		// respWriter.Header().Set("Content-Type", "application/json")

		instanceName := request.FormValue("instance_name")
		imageID := request.FormValue("image_id")
		region := request.FormValue("region")
		instanceType := request.FormValue("instance_type")

		session, errInSessionCreation := session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(*dep.SessionTokenHandler.AccessKeyId, *dep.SessionTokenHandler.SecretAccessKey, *dep.SessionTokenHandler.SessionToken),
		})

		// tmpl.Execute(respWriter, struct{ Success bool }{true})

		if errInSessionCreation != nil {
			response.Info.Success = false
			response.Info.StatusCode = 400
			response.Info.Message = fmt.Sprintf("Error while creating instance handler: %s", errInSessionCreation)

			fmt.Println("Some error while creatung instance handler: ")
			if err := json.NewEncoder(respWriter).Encode(response); err != nil {
				panic(err)
			}
			tmpl.Execute(respWriter, response)
			return
		}

		instanceHandler := ec2.New(session, &aws.Config{
			// Region: aws.String(dep.Region),
			Region: aws.String(region),
		})

		instanceResult, errInInstanceHandler := instanceHandler.RunInstances(&ec2.RunInstancesInput{
			ImageId:      aws.String(imageID),
			InstanceType: aws.String(instanceType),
			// InstanceType: aws.String(dep.InstanceName),
			MinCount: aws.Int64(dep.MinInstanceCount),
			MaxCount: aws.Int64(dep.MaxInstanceCount),
		})

		if errInInstanceHandler != nil {
			response.Info.Success = false
			response.Info.StatusCode = 400
			response.Info.Message = fmt.Sprintf("Error while creating instance: %s", errInInstanceHandler)

			// fmt.Printf(">>>>>>>>>>>>>>>>>>>>>%+v == %T\n", response, response)
			// temp, _ := json.Marshal(response)
			// json.NewEncoder(respWriter).Encode(response)
			// fmt.Printf(">>> %T\n", temp)
			// fmt.Println(">>>\n", string(temp))
			// json.NewEncoder(respWriter).Encode(response)
			// if err := json.NewEncoder(respWriter).Encode(response); err != nil {
			// 	panic(err)
			// }
			fmt.Println("Returning")
			tmpl.Execute(respWriter, response)
			return
		}

		fmt.Println("Instance Created successfully1", *instanceResult)
		fmt.Println("Instance Created successfully2", *instanceResult.Instances[0])

		_, errInTagCreation := instanceHandler.CreateTags(&ec2.CreateTagsInput{
			Resources: []*string{
				instanceResult.Instances[0].InstanceId,
			},
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(instanceName),
				},
			},
		})
		if errInTagCreation != nil {
			response.Info.Success = false
			response.Info.StatusCode = 400
			response.Info.Message = fmt.Sprintf("Error in creating instance tag: %s", errInTagCreation)

			// if err := json.NewEncoder(respWriter).Encode(response); err != nil {
			// 	panic(err)
			// }
			// fmt.Println("Error in creating instance tag: ", errInTagCreation)
			tmpl.Execute(respWriter, response)
			return

		}
		response.Info.Success = true
		response.Info.StatusCode = 200
		response.Info.Message = fmt.Sprint("Successfully created instance")
		response.Data.InstanceID = *instanceResult.Instances[0].InstanceId
		response.Data.InstanceType = *instanceResult.Instances[0].InstanceType

		// if err := json.NewEncoder(respWriter).Encode(response); err != nil {
		// 	panic(err)
		// }
		tmpl.Execute(respWriter, response)
		return
	}
}

//StopInstance : This function stops the instance when given an instance id
func StopInstance(respWriter http.ResponseWriter, request *http.Request) {
	// var inputParamsStruct StopInstanceInputParams

	// _ = json.NewDecoder(request.Body).Decode(&inputParamsStruct)

	tmpl := template.Must(template.ParseFiles("assets/stop_instance.html"))
	var response StopInstanceResponse
	if request.Method != http.MethodPost {
		tmpl.Execute(respWriter, nil)
		return
	} else if request.Method == http.MethodPost {
		fmt.Println("In stop instance post method")

		// fmt.Printf("input param: %+v", inputParamsStruct)

		// inputInstanceID := inputParamsStruct.InstanceID
		inputInstanceID := request.FormValue("instance_id")
		region := request.FormValue("region")
		fmt.Println(request)
		fmt.Println("inputInstanceId: ", inputInstanceID)

		session, errInSessionCreation := session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(*dep.SessionTokenHandler.AccessKeyId, *dep.SessionTokenHandler.SecretAccessKey, *dep.SessionTokenHandler.SessionToken),
		})

		if errInSessionCreation != nil {
			response.Info.Success = false
			response.Info.StatusCode = 400
			response.Info.Message = fmt.Sprintf("Error while creating instance handler: %s", errInSessionCreation)

			fmt.Println("Some error while creatung instance handler: ")
			// if err := json.NewEncoder(respWriter).Encode(response); err != nil {
			// 	panic(err)
			// }
			tmpl.Execute(respWriter, response)
			return
		}
		instanceHandler := ec2.New(session, &aws.Config{
			Region: aws.String(region),
		})
		fmt.Println("instanceHandler: ", instanceHandler)
		fmt.Printf("Type: %T", instanceHandler)
		stopInstanceResult, errInStopInstance := instanceHandler.StopInstances(&ec2.StopInstancesInput{
			InstanceIds: []*string{
				aws.String(inputInstanceID),
			},
		})
		fmt.Println("stopInstanceResult: ", stopInstanceResult)
		fmt.Println("errInStopInstance: ", errInStopInstance)

		if errInStopInstance != nil {
			response.Info.Success = false
			response.Info.StatusCode = 400
			response.Info.Message = fmt.Sprintf("%s", errInStopInstance)
			response.Info.Status = fmt.Sprint("Error in stopping instance")

			respWriter.WriteHeader(http.StatusBadRequest)

			// if err := json.NewEncoder(respWriter).Encode(response); err != nil {
			// 	panic(err)
			// }
			tmpl.Execute(respWriter, response)
			return
		}

		response.Info.Success = true
		response.Info.Message = fmt.Sprintf("Stopped instance : %s", inputInstanceID)
		response.Info.StatusCode = 200
		response.Info.Status = "Successfylly stopped instance"

		// if err := json.NewEncoder(respWriter).Encode(response); err != nil {
		// 	panic(err)
		// }
		tmpl.Execute(respWriter, response)
		return
	}
}

//TerminateInstance : This function terminates the instance when given an instance id
func TerminateInstance(respWriter http.ResponseWriter, request *http.Request) {
	// var inputParamsStruct StopInstanceInputParams

	// _ = json.NewDecoder(request.Body).Decode(&inputParamsStruct)

	// fmt.Printf("input param: %+v", inputParamsStruct)

	// inputInstanceID := inputParamsStruct.InstanceID
	// fmt.Println("inputInstanceId: ", inputInstanceID)

	var response TerminateInstanceResponse

	tmpl := template.Must(template.ParseFiles("assets/terminate_instance.html"))
	if request.Method != http.MethodPost {
		response.Info.Success = true
		tmpl.Execute(respWriter, nil)
		return
	} else if request.Method == http.MethodPost {
		session, errInSessionCreation := session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(*dep.SessionTokenHandler.AccessKeyId, *dep.SessionTokenHandler.SecretAccessKey, *dep.SessionTokenHandler.SessionToken),
		})

		if errInSessionCreation != nil {
			response.Info.Success = false
			response.Info.StatusCode = 400
			response.Info.Message = fmt.Sprintf("Error while creating instance handler: %s", errInSessionCreation)

			fmt.Println("Some error while creatung instance handler: ")
			// if err := json.NewEncoder(respWriter).Encode(response); err != nil {
			// 	panic(err)
			// }
			tmpl.Execute(respWriter, response)
			return
		}

		inputInstanceID := request.FormValue("instance_id")
		region := request.FormValue("region")
		instanceHandler := ec2.New(session, &aws.Config{
			Region: aws.String(region),
		})

		terminateInstanceResult, errInTerminateInstance := instanceHandler.TerminateInstances(&ec2.TerminateInstancesInput{
			InstanceIds: []*string{
				aws.String(inputInstanceID),
			},
		})
		fmt.Println("terminateInstanceResult: ", terminateInstanceResult)
		fmt.Println("errInTerminateInstance: ", errInTerminateInstance)

		if errInTerminateInstance != nil {
			response.Info.Success = false
			response.Info.StatusCode = 400
			response.Info.Message = fmt.Sprintf("%s", errInTerminateInstance)
			response.Info.Status = fmt.Sprint("Error in terminate instance")

			respWriter.WriteHeader(http.StatusBadRequest)

			// if err := json.NewEncoder(respWriter).Encode(response); err != nil {
			// 	panic(err)
			// }
			tmpl.Execute(respWriter, response)
			return
		}

		response.Info.Success = true
		response.Info.Message = fmt.Sprintf("Terminated instance : %s", inputInstanceID)
		response.Info.StatusCode = 200
		response.Info.Status = "Successfylly terminated instance"

		// if err := json.NewEncoder(respWriter).Encode(response); err != nil {
		// 	panic(err)
		// }
		tmpl.Execute(respWriter, response)
		return
	}
}

// DescribeInstance : This function sends the necessary info/stats of instance when given
// an instance id as an input
func DescribeInstance(respWriter http.ResponseWriter, request *http.Request) {
	// inputParams := mux.Vars(request)

	// inputInstanceID := inputParams["instance_id"]
	fmt.Println("In desc inst")
	var response DescribeInstanceResponse

	tmpl := template.Must(template.ParseFiles("assets/describe_instance.html"))
	if request.Method != http.MethodPost {
		response.Info.Success = true
		tmpl.Execute(respWriter, nil)
		return
	} else if request.Method == http.MethodPost {
		session, errInSessionCreation := session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(*dep.SessionTokenHandler.AccessKeyId, *dep.SessionTokenHandler.SecretAccessKey, *dep.SessionTokenHandler.SessionToken),
		})
		inputInstanceID := request.FormValue("instance_id")
		region := request.FormValue("region")
		fmt.Println(inputInstanceID, region)

		if errInSessionCreation != nil {
			response.Info.Success = false
			response.Info.StatusCode = 400
			response.Info.Message = fmt.Sprintf("Error while creating instance handler: %s", errInSessionCreation)

			fmt.Println("Some error while creatung instance handler: ")
			// if err := json.NewEncoder(respWriter).Encode(response); err != nil {
			// 	panic(err)
			// }
			tmpl.Execute(respWriter, response)
			return
		}
		instanceHandler := ec2.New(session, &aws.Config{
			Region: aws.String(region),
		})

		describeInstanceResult, errInDescribeInstance := instanceHandler.DescribeInstanceStatus(&ec2.DescribeInstanceStatusInput{
			InstanceIds: []*string{
				aws.String(inputInstanceID),
			},
		})
		fmt.Println("describeInstanceResult: ", describeInstanceResult)
		fmt.Println("errInDescribeInstance: ", errInDescribeInstance)

		if errInDescribeInstance != nil {
			response.Info.Success = false
			response.Info.StatusCode = 400
			response.Info.Message = fmt.Sprintf("%s", errInDescribeInstance)
			response.Info.Status = fmt.Sprint("Error in describe instance")

			// respWriter.WriteHeader(http.StatusBadRequest)

			// if err := json.NewEncoder(respWriter).Encode(response); err != nil {
			// 	panic(err)
			// }
			tmpl.Execute(respWriter, response)
			return
		}

		response.Info.Success = true
		response.Info.Message = fmt.Sprintf("Fetched instance data : %s", inputInstanceID)
		response.Info.StatusCode = 200
		response.Info.Status = "Successfylly fetched instance data"
		response.Data.InstanceID = *describeInstanceResult.InstanceStatuses[0].InstanceId
		response.Data.InstanceStatus = *describeInstanceResult.InstanceStatuses[0].InstanceState.Name

		// if err := json.NewEncoder(respWriter).Encode(response); err != nil {
		// 	panic(err)
		// }
		tmpl.Execute(respWriter, response)
		return
	}
}

// DescribeInstance : This function sends the necessary info/stats of instance when given
// an instance id as an input
func ListInstance(respWriter http.ResponseWriter, request *http.Request) {
	inputParams := mux.Vars(request)

	inputInstanceID := inputParams["instance_id"]

	var response DescribeInstanceResponse

	session, errInSessionCreation := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(*dep.SessionTokenHandler.AccessKeyId, *dep.SessionTokenHandler.SecretAccessKey, *dep.SessionTokenHandler.SessionToken),
	})

	if errInSessionCreation != nil {
		response.Info.Success = false
		response.Info.StatusCode = 400
		response.Info.Message = fmt.Sprintf("Error while creating instance handler: %s", errInSessionCreation)

		fmt.Println("Some error while creatung instance handler: ")
		if err := json.NewEncoder(respWriter).Encode(response); err != nil {
			panic(err)
		}
		return
	}
	instanceHandler := ec2.New(session, &aws.Config{
		Region: aws.String(dep.Region),
	})

	describeInstanceResult, errInDescribeInstance := instanceHandler.DescribeInstanceStatus(&ec2.DescribeInstanceStatusInput{
		InstanceIds: []*string{
			aws.String(inputInstanceID),
		},
		// Region: []*string{},
	})
	fmt.Println("describeInstanceResult: ", describeInstanceResult)
	fmt.Println("errInDescribeInstance: ", errInDescribeInstance)

	if errInDescribeInstance != nil {
		response.Info.Success = false
		response.Info.StatusCode = 400
		response.Info.Message = fmt.Sprintf("%s", errInDescribeInstance)
		response.Info.Status = fmt.Sprint("Error in describe instance")

		respWriter.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(respWriter).Encode(response); err != nil {
			panic(err)
		}
		return
	}

	response.Info.Success = true
	response.Info.Message = fmt.Sprintf("Fetched instance data : %s", inputInstanceID)
	response.Info.StatusCode = 200
	response.Info.Status = "Successfylly fetched instance data"
	response.Data.InstanceID = *describeInstanceResult.InstanceStatuses[0].InstanceId
	response.Data.InstanceStatus = *describeInstanceResult.InstanceStatuses[0].InstanceState.Name

	if err := json.NewEncoder(respWriter).Encode(response); err != nil {
		panic(err)
	}
	return
}

type PageVariables struct {
	Date string
	Time string
}

func HomePage(w http.ResponseWriter, r *http.Request) {

	// if r.Method == "GET" {
	// 	t := template.Must("assets/login.html")
	// 	t.Execute(w, nil)
	// 	return
	// }
	r.ParseForm()
	// logic part of log in
	fmt.Println("username:", r.Form["username"])
	fmt.Println("password:", r.Form["password"])
	// ----------------------
	// r.ParseForm()
	// fmt.Println(r.Form) // print information on server side.

	now := time.Now()              // find the time right now
	HomePageVars := PageVariables{ //store the date and time in a struct
		Date: now.Format("02-01-2006"),
		Time: now.Format("15:04:05"),
	}

	// fmt.Println("path", r.URL.Path)
	// fmt.Println("scheme", r.URL.Scheme)
	// fmt.Println(r.Form["url_long"])
	// for k, v := range r.Form {
	// 	fmt.Println("key:", k)
	// 	fmt.Println("val:", strings.Join(v, ""))
	// }
	// fmt.Fprintf(w, "Hello astaxie!") // write data to response

	t, err := template.ParseFiles("assets/homepage.html") //parse the html file homepage.html
	if err != nil {                                       // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                  // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func HomePage1(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("assets/generate_token.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
	if r.Method == http.MethodGet {
		fmt.Println("username: ", r.FormValue("email"))
		fmt.Println("username: ", r.FormValue("subject"))
		tmpl.Execute(w, struct{ Success bool }{true})
	}
}
