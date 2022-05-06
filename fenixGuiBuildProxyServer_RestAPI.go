package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
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
	router.HandleFunc("/pinned-testinstructions-and-testinstructioncontainers/{userid}", fenixGuiBuilderProxyServerObject.RestSendGetPinnedInstructionsAndTestInstructionContainersToFenixGuiBuilderServer).Methods("GET")
	router.HandleFunc("/pinned-testinstructions-and-testinstructioncontainers/{userid}{pinnedTestInstructionsAndTestInstructionsContainers}", fenixGuiBuilderProxyServerObject.RestSendSavePinnedInstructionsAndTestInstructionContainersToFenixGuiBuilderServer).Methods("POST")

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
	response = fenixGuiBuilderProxyServerObject.SendGetTestInstructionsAndTestContainers(userId)

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

func (fenixGuiBuilderProxyServerObject *fenixGuiBuilderProxyServerObjectStruct) RestSendGetPinnedInstructionsAndTestInstructionContainersToFenixGuiBuilderServer(w http.ResponseWriter, r *http.Request) {
	// curl --request GET localhost:8080/pinned-testinstructions-and-testinstructioncontainers/s41797

	fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
		"id": "2472dda1-701d-4b23-8326-757e43df4af4",
	}).Debug("Incoming 'RestApi - /pinned-testinstructions-and-testinstructioncontainers'")

	defer fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
		"id": "db318ff4-ad36-43d4-a8d4-3e0ac4ff08c6",
	}).Debug("Outgoing 'RestApi - /pinned-testinstructions-and-testinstructioncontainers'")

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
	response = fenixGuiBuilderProxyServerObject.SendGetPinnedTestInstructionsAndTestContainers(userId)

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

func (fenixGuiBuilderProxyServerObject *fenixGuiBuilderProxyServerObjectStruct) RestSendSavePinnedInstructionsAndTestInstructionContainersToFenixGuiBuilderServer(w http.ResponseWriter, r *http.Request) {
	// curl --request POST localhost:8080/pinned-testinstructions-and-testinstructioncontainers/s41797

	fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
		"id": "2472dda1-701d-4b23-8326-757e43df4af4",
	}).Debug("Incoming 'RestApi - (POST) /pinned-testinstructions-and-testinstructioncontainers'")

	defer fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
		"id": "db318ff4-ad36-43d4-a8d4-3e0ac4ff08c6",
	}).Debug("Outgoing 'RestApi - (POST) /pinned-testinstructions-and-testinstructioncontainers'")

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

	// Cast pinned TestInstructions and TestInstructionContainer - json into struct
	pinnedTestInstructionsAndTestContainersMessage := fenixGuiTestCaseBuilderServerGrpcApi.PinnedTestInstructionsAndTestContainersMessage{}
	err := jsonpb.Unmarshal(r.Body, &pinnedTestInstructionsAndTestContainersMessage)

	// If casting json into proto-struct didn't succeed then return  error message
	if exit == false {
		fmt.Fprintf(w, "Couldn't convert json due not correct format. ")

		return
	}

	// Set the user in the message
	pinnedTestInstructionsAndTestContainersMessage.UserId = userId

	// Do gRPC-call
	response = fenixGuiBuilderProxyServerObject.SendGetPinnedTestInstructionsAndTestContainers(userId)

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
