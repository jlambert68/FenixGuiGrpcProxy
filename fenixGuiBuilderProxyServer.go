package main

import (
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
)

// Used for only process cleanup once
var cleanupProcessed = false

func cleanup() {

	if cleanupProcessed == false {

		cleanupProcessed = true

		// Cleanup before close down application
		fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{}).Info("Clean up and shut down servers")

		// Stop Backend gRPC Server
		fenixGuiBuilderProxyServerObject.StopGrpcServer()

		//log.Println("Close DB_session: %v", DB_session)
		//DB_session.Close()
	}
}

func fenixGuiTestCaseBuilderServerMain() {

	// Connect to CloudDB
	fenixSyncShared.ConnectToDB()

	// Set up BackendObject
	fenixGuiBuilderProxyServerObject = &fenixGuiBuilderProxyServerObjectStruct{}

	// Init logger
	fenixGuiBuilderProxyServerObject.InitLogger("")

	// Clean up when leaving. Is placed after logger because shutdown logs information
	defer cleanup()

	// Start RestApi-server
	go fenixGuiBuilderProxyServerObject.restAPIServer()

	// Start Backend gRPC-server
	fenixGuiBuilderProxyServerObject.InitGrpcServer()

}
