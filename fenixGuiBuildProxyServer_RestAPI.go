package main

import (
	"encoding/json"
	"fmt"
	fenixGuiTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (fenixGuiBuilderProxyServerObject *fenixGuiBuilderProxyServerObjectStruct) restAPIServer() {
	log.Println("starting API server")
	//create a new router
	router := mux.NewRouter()
	log.Println("creating routes")
	//specify endpoints
	router.HandleFunc("/health-check", fenixGuiBuilderProxyServerObject.HealthCheck).Methods("GET")
	router.HandleFunc("/are-guibuilderserver-alive", fenixGuiBuilderProxyServerObject.RestSendAreYouAliveToFenixGuiBuilderServer).Methods("GET")
	router.HandleFunc("/testinstructions-and-testinstructioncontainers/{userid}", fenixGuiBuilderProxyServerObject.RestSendGetInstructionsAndTestInstructionContainersToFenixGuiBuilderServer).Methods("GET")

	http.Handle("/", router)

	//start and listen to requests
	http.ListenAndServe(":8080", router)

}

func (fenixGuiBuilderProxyServerObject *fenixGuiBuilderProxyServerObjectStruct) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// curl --request GET localhost:8080/health-check

	fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
		"id": "fb3c1ecb-3da8-4d27-b1c4-16d5120e7125",
	}).Debug("Incoming 'RestApi - /health-check'")

	defer fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
		"id": "fab7676d-c303-4b20-8980-397d7a59282e",
	}).Debug("Outgoing 'RestApi - /health-check'")

	// Set OK in Header
	w.WriteHeader(http.StatusOK)

	// Create Response message
	fmt.Fprintf(w, "API is up and running")
}

func (fenixGuiBuilderProxyServerObject *fenixGuiBuilderProxyServerObjectStruct) RestSendAreYouAliveToFenixGuiBuilderServer(w http.ResponseWriter, r *http.Request) {
	// curl --request GET localhost:8080/are-guibuilderserver-alive

	fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
		"id": "0645d30c-4479-49ab-bb72-9bc3fac329a5",
	}).Debug("Incoming 'RestApi - /are-guibuilderserver-alive'")

	defer fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
		"id": "cc168cfe-3544-4946-93d4-d2325893f8cd",
	}).Debug("Outgoing 'RestApi - /are-guibuilderserver-alive'")

	// gRPC -response
	var response *fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse

	// Do gRPC-call
	response = fenixGuiBuilderProxyServerObject.SendAreYouAliveToFenixGuiBuilderServer()

	// Create Header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Convert gRPC-response into json
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		// if error then just exit TODO Create correct response message
		return
	}

	// Create Response message
	w.Write(jsonResponse)
}

func (fenixGuiBuilderProxyServerObject *fenixGuiBuilderProxyServerObjectStruct) RestSendGetInstructionsAndTestInstructionContainersToFenixGuiBuilderServer(w http.ResponseWriter, r *http.Request) {
	// curl --request GET localhost:8080/testinstructions-and-testinstructioncontainers/s41797

	fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
		"id": "0645d30c-4479-49ab-bb72-9bc3fac329a5",
	}).Debug("Incoming 'RestApi - /testinstructions-and-testinstructioncontainers'")

	defer fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
		"id": "cc168cfe-3544-4946-93d4-d2325893f8cd",
	}).Debug("Outgoing 'RestApi - /testinstructions-and-testinstructioncontainers'")

	// gRPC -response
	var response *fenixGuiTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage

	// Extract UserId
	parameters := mux.Vars(r)
	userId, exit := parameters["userid"]

	// If parameter UserID is missing then return error message
	if exit == false {
		fmt.Fprintf(w, "Missing UserId")

		return
	}

	// Do gRPC-call
	response = fenixGuiBuilderProxyServerObject.SendAGetTestInstructionsAndTestContainers(userId)

	// Create Header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Convert gRPC-response into json
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		// if error then just exit TODO Create correct response message
		fmt.Fprintf(w, err.Error())

		return
	}

	// Create Response message
	w.Write(jsonResponse)
}
