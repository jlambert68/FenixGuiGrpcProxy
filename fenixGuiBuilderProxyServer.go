package main

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path/filepath"
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
	//fenixSyncShared.ConnectToDB()

	// Set up BackendObject
	fenixGuiBuilderProxyServerObject = &fenixGuiBuilderProxyServerObjectStruct{runAsTrayApplication: tempRunAsTrayApplication}

	// Init logger
	// When application is run as tray application then use text file as log
	var filePathName = ""
	var err error

	if fenixGuiBuilderProxyServerObject.runAsTrayApplication == true {
		// Get path for this application

		logfilename := "mylog.log"
		filePathName, err = filepath.Abs(logfilename)
		if err != nil {
			log.Println("Couldn't generate filePathName for log: ", err)
			os.Exit(0)
		}
	}

	fenixGuiBuilderProxyServerObject.InitLogger(filePathName)

	// Clean up when leaving. Is placed after logger because shutdown logs information
	defer cleanup()

	// Start RestApi-server
	go fenixGuiBuilderProxyServerObject.restAPIServer()

	// Start Backend gRPC-server
	fenixGuiBuilderProxyServerObject.InitGrpcServer()

}
