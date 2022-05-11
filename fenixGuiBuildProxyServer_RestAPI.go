package main

import (
	"encoding/json"
	"fmt"
	fenixGuiTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"strings"

	"errors"
	"github.com/golang/gddo/httputil/header"
	"github.com/gorilla/mux"
)

// Structs used when converting json messages in RestAPI

// RestUserMessageStruct
type RestUserMessageStruct struct {
	UserId string `protobuf:"bytes,1,opt,name=UserId,proto3" json:"UserId,omitempty"`
}

type RestSavePinnedInstructionsAndTestInstructionContainersToFenixGuiBuilderServerStruct struct {
	UserId                                 string                                                                        `protobuf:"bytes,1,opt,name=UserId,proto3" json:"UserId,omitempty"`
	PinnedTestInstructionMessages          []*fenixGuiTestCaseBuilderServerGrpcApi.PinnedTestInstructionMessage          `protobuf:"bytes,3,rep,name=PinnedTestInstructionMessages,proto3" json:"PinnedTestInstructionMessages,omitempty"`
	PinnedTestInstructionContainerMessages []*fenixGuiTestCaseBuilderServerGrpcApi.PinnedTestInstructionContainerMessage `protobuf:"bytes,4,rep,name=PinnedTestInstructionContainerMessages,proto3" json:"PinnedTestInstructionContainerMessages,omitempty"`
}

func (fenixGuiBuilderProxyServerObject *fenixGuiBuilderProxyServerObjectStruct) restAPIServer() {
	log.Println("starting API server")
	//create a new router
	router := mux.NewRouter()
	log.Println("creating routes")
	//specify endpoints
	router.HandleFunc("/health-check", fenixGuiBuilderProxyServerObject.HealthCheck).Methods("GET")
	router.HandleFunc("/are-guibuilderserver-alive", fenixGuiBuilderProxyServerObject.RestSendAreYouAliveToFenixGuiBuilderServer).Methods("GET")
	router.HandleFunc("/testinstructions-and-testinstructioncontainers", fenixGuiBuilderProxyServerObject.RestSendGetInstructionsAndTestInstructionContainersToFenixGuiBuilderServer).Methods("GET")
	router.HandleFunc("/pinned-testinstructions-and-testinstructioncontainers", fenixGuiBuilderProxyServerObject.RestSendGetPinnedInstructionsAndTestInstructionContainersToFenixGuiBuilderServer).Methods("GET")
	router.HandleFunc("/pinned-testinstructions-and-testinstructioncontainers", fenixGuiBuilderProxyServerObject.RestSendSavePinnedInstructionsAndTestInstructionContainersToFenixGuiBuilderServer).Methods("POST")

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
	/*
		curl -X GET \
		localhost:8080/testinstructions-and-testinstructioncontainers \
		-H 'Content-Type: application/json' \
		-d '{"UserId":"s41797"}'
	*/
	fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
		"id": "0645d30c-4479-49ab-bb72-9bc3fac329a5",
	}).Debug("Incoming 'RestApi - (GET) /testinstructions-and-testinstructioncontainers'")

	defer fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
		"id": "cc168cfe-3544-4946-93d4-d2325893f8cd",
	}).Debug("Outgoing 'RestApi - (GET) /testinstructions-and-testinstructioncontainers'")

	// gRPC -response
	var response *fenixGuiTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage

	// Variable where Rest-json-payload will end up in
	jsonData := &RestUserMessageStruct{}

	// Extract and Validate json body
	err := extractAndValidateJsonBody(&w, r, jsonData)
	if err != nil {
		// If something went wrong then just exit
		return
	}

	// Do gRPC-call
	response = fenixGuiBuilderProxyServerObject.SendGetTestInstructionsAndTestContainers(jsonData.UserId)

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
	/*
		curl -X GET \
		localhost:8080/pinned-testinstructions-and-testinstructioncontainers \
		-H 'Content-Type: application/json' \
		-d '{"UserId":"s41797"}'
	*/

	fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
		"id": "2472dda1-701d-4b23-8326-757e43df4af4",
	}).Debug("Incoming 'RestApi - /pinned-testinstructions-and-testinstructioncontainers'")

	defer fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
		"id": "db318ff4-ad36-43d4-a8d4-3e0ac4ff08c6",
	}).Debug("Outgoing 'RestApi - /pinned-testinstructions-and-testinstructioncontainers'")

	// gRPC -response
	var response *fenixGuiTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage

	// Variable where Rest-json-payload will end up in
	jsonData := &RestUserMessageStruct{}

	// Extract and Validate json body
	err := extractAndValidateJsonBody(&w, r, jsonData)
	if err != nil {
		// If something went wrong then just exit
		return
	}

	// Do gRPC-call
	response = fenixGuiBuilderProxyServerObject.SendGetPinnedTestInstructionsAndTestContainers(jsonData.UserId)

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
	/*
		curl -X POST localhost:8080/pinned-testinstructions-and-testinstructioncontainers \
		-H 'Content-Type: application/json' \
		-d '{"UserId":"s41797","PinnedTestInstructionMessages":[{"TestInstructionUuid":"2f130d7e-f8aa-466f-b29d-0fb63608c1a6","TestInstructionName":"TestInstructionName 1"}],"PinnedTestInstructionContainerMessages":[{"TestInstructionContainerUuid":"b107bdd9-4152-4020-b3f0-fc750b45885e","TestInstructionContainerName":"TestInstructionContainerName 1"},{"TestInstructionContainerUuid":"e81b9734-5dce-43c9-8d77-3368940cf126","TestInstructionContainerName":"TestInstructionContainerName"}]}'
	*/
	// curl -X POST localhost:8080/pinned-testinstructions-and-testinstructioncontainers -H 'Content-Type: application/json' -d '{"UserId":"s41797","PinnedTestInstructionMessages":[{"TestInstructionUuid":"myUuid", "TestInstructionName":"myName"}],"PinnedTestInstructionContainerMessages":[{"TestInstructionContainerUuid":"myUuid2", "TestInstructionContainerName":"myName2"}]}'
	fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
		"id": "2472dda1-701d-4b23-8326-757e43df4af4",
	}).Debug("Incoming 'RestApi - (POST) /pinned-testinstructions-and-testinstructioncontainers'")

	defer fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
		"id": "db318ff4-ad36-43d4-a8d4-3e0ac4ff08c6",
	}).Debug("Outgoing 'RestApi - (POST) /pinned-testinstructions-and-testinstructioncontainers'")

	// Variable where Rest-json-payload will end up in
	jsonData := &RestSavePinnedInstructionsAndTestInstructionContainersToFenixGuiBuilderServerStruct{}

	// Extract and Validate json body
	err := extractAndValidateJsonBody(&w, r, jsonData)
	if err != nil {
		// If something went wrong then just exit
		return
	}

	// Create gRPC -response variable
	var response *fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse

	// Create input message for gRPC-call
	pinnedTestInstructionsAndTestContainersMessage := &fenixGuiTestCaseBuilderServerGrpcApi.PinnedTestInstructionsAndTestContainersMessage{
		UserId: jsonData.UserId,
		ProtoFileVersionUsedByClient: fenixGuiTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(
			fenixGuiBuilderProxyServerObject.getHighestFenixTestDataProtoFileVersion()),
		PinnedTestInstructionMessages:          jsonData.PinnedTestInstructionMessages,
		PinnedTestInstructionContainerMessages: jsonData.PinnedTestInstructionContainerMessages,
	}

	// Do gRPC-call
	response = fenixGuiBuilderProxyServerObject.SendSavePinnedTestInstructionsAndTestContainers(pinnedTestInstructionsAndTestContainersMessage)

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

func extractAndValidateJsonBody(responseWriterPointer *http.ResponseWriter, httpRequest *http.Request, myInputTypeVariable interface{}) (err error) {
	// If the Content-Type header is present, check that it has the value
	// application/json. Note that we are using the gddo/httputil/header
	// package to parse and extract the value here, so the check works
	// even if the client includes additional charset or boundary
	// information in the header.
	responseWriter := *responseWriterPointer
	if httpRequest.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(httpRequest.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(responseWriter, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	// Use http.MaxBytesReader to enforce a maximum read of 1MB from the
	// response body. A request body larger than that will now result in
	// Decode() returning a "http: request body too large" error.
	httpRequest.Body = http.MaxBytesReader(responseWriter, httpRequest.Body, 1048576)

	// Setup the decoder and call the DisallowUnknownFields() method on it.
	// This will cause Decode() to return a "json: unknown field ..." error
	// if it encounters any extra unexpected fields in the JSON. Strictly
	// speaking, it returns an error for "keys which do not match any
	// non-ignored, exported fields in the destination".
	dec := json.NewDecoder(httpRequest.Body)
	dec.DisallowUnknownFields()

	var p = myInputTypeVariable //RestSavePinnedInstructionsAndTestInstructionContainersToFenixGuiBuilderServerStruct
	err = dec.Decode(&p)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// Catch any syntax errors in the JSON and send an error message
		// which interpolates the location of the problem to make it
		// easier for the client to fix.
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			http.Error(responseWriter, msg, http.StatusBadRequest)

		// In some circumstances Decode() may also return an
		// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
		// is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			http.Error(responseWriter, msg, http.StatusBadRequest)

		// Catch any type errors, like trying to assign a string in the
		// JSON request body to a int field in our RestSavePinnedInstructionsAndTestInstructionContainersToFenixGuiBuilderServerStruct struct. We can
		// interpolate the relevant field name and position into the error
		// message to make it easier for the client to fix.
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			http.Error(responseWriter, msg, http.StatusBadRequest)

		// Catch the error caused by extra unexpected fields in the request
		// body. We extract the field name from the error message and
		// interpolate it in our custom error message. There is an open
		// issue at https://github.com/golang/go/issues/29035 regarding
		// turning this into a sentinel error.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			http.Error(responseWriter, msg, http.StatusBadRequest)

		// An io.EOF error is returned by Decode() if the request body is
		// empty.
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			http.Error(responseWriter, msg, http.StatusBadRequest)

		// Catch the error caused by the request body being too large. Again
		// there is an open issue regarding turning this into a sentinel
		// error at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			http.Error(responseWriter, msg, http.StatusRequestEntityTooLarge)

		// Otherwise default to logging the error and sending a 500 Internal
		// Server Error response.
		default:
			log.Println(err.Error())
			http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return err
	}

	// Call decode again, using a pointer to an empty anonymous struct as
	// the destination. If the request body only contained a single JSON
	// object this will return an io.EOF error. So if we get anything else,
	// we know that there is additional data in the request body.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		http.Error(responseWriter, msg, http.StatusBadRequest)
		return
	}

	//fmt.Fprintf(responseWriter, "RestSavePinnedInstructionsAndTestInstructionContainersToFenixGuiBuilderServerStruct: %+v", p)

	return nil
}
