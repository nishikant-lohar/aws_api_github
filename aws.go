package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"aws_api/apis/ec2"
	"aws_api/apis/s3"
	"aws_api/apis/session"
)

func main() {

	router := mux.NewRouter()

	staticFileDirectory := http.Dir("./assets/")
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	router.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")

	// router URLs
	router.HandleFunc("/", session.GetStsTokenApi)
	// router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "OK")
	// })

	// ---------- EC2 URLS -------------------
	router.HandleFunc("/create-instance", ec2.CreateInstance)
	router.HandleFunc("/stop-instance", ec2.StopInstance)
	router.HandleFunc("/terminate-instance", ec2.TerminateInstance)
	router.HandleFunc("/describe-instance", ec2.DescribeInstance)
	router.HandleFunc("/list-instance", ec2.ListInstance).Methods("GET")

	// ---------- S3 URLS -------------------
	router.HandleFunc("/create-bucket", s3.CreateBucket)

	http.ListenAndServe(":8000", router)
}

// func init() {
// 	// dep.SessionTokenHandler = session.GetStsToken()
// 	fmt.Println("Token Created")
// 	fmt.Println("This is token session: ", dep.SessionTokenHandler)

// 	// ec2.CreateInstance()
// }
