package main

import (
	"FenixGuiGrpcProxyServer/common_config"
	"crypto/tls"
	fenixGuiTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"time"
)

// ********************************************************************************************************************

// SetConnectionToFenixTestDataSyncServer - Set upp connection and Dial to FenixTestDataSyncServer
func (fenixGuiBuilderProxyServerObject *fenixGuiBuilderProxyServerObjectStruct) SetConnectionToFenixGuiBuilderServer() {

	var err error
	var opts []grpc.DialOption

	//When running on GCP then use credential otherwise not
	if common_config.ExecutionLocationForFenixGuiServer == common_config.GCP {
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})

		opts = []grpc.DialOption{
			grpc.WithTransportCredentials(creds),
		}
	}

	// Set up connection to FenixTestDataSyncServer
	// When run on GCP, use credentials
	if common_config.ExecutionLocationForFenixGuiServer == common_config.GCP {
		// Run on GCP
		remoteFenixGuiBuilderServerConnection, err = grpc.Dial(fenixGuiBuilderServerAddressToDial, opts...)
	} else {
		// Run Local
		remoteFenixGuiBuilderServerConnection, err = grpc.Dial(fenixGuiBuilderServerAddressToDial, grpc.WithInsecure())
	}
	if err != nil {
		fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
			"ID":                                 "50b59b1b-57ce-4c27-aa84-617f0cde3100",
			"fenixGuiBuilderServerAddressToDial": fenixGuiBuilderServerAddressToDial,
			"error message":                      err,
		}).Error("Did not connect to FenixGuiBuilderServer via gRPC")
		//os.Exit(0)
	} else {
		fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
			"ID": "0c650bbc-45d0-4029-bd25-4ced9925a059",
			"fenixGuiTestCaseBuilderServer_address_to_dial": fenixGuiBuilderServerAddressToDial,
		}).Info("gRPC connection OK to FenixTestDataSyncServer")

		// Creates a new Clients
		fenixGuiBuilderServerGrpcClient = fenixGuiTestCaseBuilderServerGrpcApi.NewFenixTestCaseBuilderServerGrpcServicesClient(remoteFenixGuiBuilderServerConnection)

	}
}

// ********************************************************************************************************************

// SendAreYouAliveToFenixGuiBuilderServer - Check if FenixGuiBuilderServer is alive
func (fenixGuiBuilderProxyServerObject *fenixGuiBuilderProxyServerObjectStruct) SendAreYouAliveToFenixGuiBuilderServer() (returnMessage *fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse) {

	var ctx context.Context
	var returnMessageAckNack bool
	var returnMessageString string
	var err error

	// Set up connection to Server
	fenixGuiBuilderProxyServerObject.SetConnectionToFenixGuiBuilderServer()

	// Create the request message
	emptyParameter := &fenixGuiTestCaseBuilderServerGrpcApi.EmptyParameter{

		ProtoFileVersionUsedByClient: fenixGuiTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(
			fenixGuiBuilderProxyServerObject.getHighestFenixTestDataProtoFileVersion()),
	}

	// Do gRPC-call
	//ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		//TODO Fixa så att denna inte görs som allt går bra
		fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
			"ID": "88bbfecb-b9a4-4e2f-92e0-cabbbdf75dc8",
		}).Error("Running Defer Cancel function")
		cancel()
	}()

	// Only add access token when run on GCP
	if common_config.ExecutionLocationForFenixGuiServer == common_config.GCP {

		// Add Access token
		ctx, returnMessageAckNack, returnMessageString = fenixGuiBuilderProxyServerObject.generateGCPAccessToken(ctx)
		if returnMessageAckNack == false {
			// When error
			returnMessage = &fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   returnMessageString,
				ErrorCodes: nil,
			}

			return returnMessage
		}

	}

	// Do the gRPC-call
	returnMessage, err = fenixGuiBuilderServerGrpcClient.AreYouAlive(ctx, emptyParameter)

	// Shouldn't happen
	if err != nil {
		fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
			"ID":    "818aaf0b-4112-4be4-97b9-21cc084c7b8b",
			"error": err,
		}).Error("Problem to do gRPC-call to FenixTestGuiBuilderServer for 'SendAreYouAliveToFenixGuiBuilderServer'")

	} else if returnMessage.AckNack == false {
		// FenixTestGuiBuilderServer couldn't handle gPRC call
		fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
			"ID":                                     "2ecbc800-2fb6-4e88-858d-a421b61c5529",
			"Message from FenixTestGuiBuilderServer": returnMessage.Comments,
		}).Error("Problem to do gRPC-call to FenixTestGuiBuilderServer for 'SendAreYouAliveToFenixGuiBuilderServer'")
	}

	return returnMessage

}

// ********************************************************************************************************************

// SendGetTestInstructionsAndTestContainers - Get available TestInstructions and TestInstructionContainers
func (fenixGuiBuilderProxyServerObject *fenixGuiBuilderProxyServerObjectStruct) SendGetTestInstructionsAndTestContainers(userId string) (returnMessage *fenixGuiTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage) {

	var ctx context.Context
	var returnMessageAckNack bool
	var returnMessageString string
	var err error

	// Set up connection to Server
	fenixGuiBuilderProxyServerObject.SetConnectionToFenixGuiBuilderServer()

	// Create the request message
	userIdentificationMessage := &fenixGuiTestCaseBuilderServerGrpcApi.UserIdentificationMessage{
		UserId: userId,
		ProtoFileVersionUsedByClient: fenixGuiTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(
			fenixGuiBuilderProxyServerObject.getHighestFenixTestDataProtoFileVersion()),
	}

	// Do gRPC-call
	//ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		//TODO Fixa så att denna inte görs som allt går bra
		fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
			"ID": "797c00e1-510d-4cbe-a48b-dc63828ecd7e",
		}).Error("Running Defer Cancel function")
		cancel()
	}()

	// Only add access token when run on GCP
	if common_config.ExecutionLocationForFenixGuiServer == common_config.GCP {

		// Add Access token
		ctx, returnMessageAckNack, returnMessageString = fenixGuiBuilderProxyServerObject.generateGCPAccessToken(ctx)
		if returnMessageAckNack == false {
			// When error
			ackNackResponse := &fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   returnMessageString,
				ErrorCodes: nil,
			}

			returnMessage = &fenixGuiTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage{
				TestInstructionMessages:          nil,
				TestInstructionContainerMessages: nil,
				AckNackResponse:                  ackNackResponse,
			}

			return returnMessage
		}

	}

	// Do the gRPC-call
	returnMessage, err = fenixGuiBuilderServerGrpcClient.GetTestInstructionsAndTestContainers(ctx, userIdentificationMessage)

	// Shouldn't happen
	if err != nil {
		fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
			"ID":    "d7235084-33e5-43a2-9fa7-dfb05ec6869e",
			"error": err,
		}).Error("Problem to do gRPC-call to FenixTestGuiBuilderServer for 'SendGetTestInstructionsAndTestContainers'")

	} else if returnMessage.AckNackResponse.AckNack == false {
		// FenixTestGuiBuilderServer couldn't handle gPRC call
		fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
			"ID":                                     "30e6f1ee-202a-47bf-a2c4-5066d0f8cf75",
			"Message from FenixTestGuiBuilderServer": returnMessage.AckNackResponse.Comments,
		}).Error("Problem to do gRPC-call to FenixTestGuiBuilderServer for 'SendGetTestInstructionsAndTestContainers'")
	}

	return returnMessage

}

// SendGetPinnedTestInstructionsAndTestContainers - Get pinned TestInstructions and TestInstructionContainers
func (fenixGuiBuilderProxyServerObject *fenixGuiBuilderProxyServerObjectStruct) SendGetPinnedTestInstructionsAndTestContainers(userId string) (returnMessage *fenixGuiTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage) {

	var ctx context.Context
	var returnMessageAckNack bool
	var returnMessageString string
	var err error

	// Set up connection to Server
	fenixGuiBuilderProxyServerObject.SetConnectionToFenixGuiBuilderServer()

	// Create the request message
	userIdentificationMessage := &fenixGuiTestCaseBuilderServerGrpcApi.UserIdentificationMessage{
		UserId: userId,
		ProtoFileVersionUsedByClient: fenixGuiTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(
			fenixGuiBuilderProxyServerObject.getHighestFenixTestDataProtoFileVersion()),
	}

	// Do gRPC-call
	//ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		//TODO Fixa så att denna inte görs som allt går bra
		fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
			"ID": "c5ba19bd-75ff-4366-818d-745d4d7f1a52",
		}).Error("Running Defer Cancel function")
		cancel()
	}()

	// Only add access token when run on GCP
	if common_config.ExecutionLocationForFenixGuiServer == common_config.GCP {

		// Add Access token
		ctx, returnMessageAckNack, returnMessageString = fenixGuiBuilderProxyServerObject.generateGCPAccessToken(ctx)
		if returnMessageAckNack == false {
			// When error
			ackNackResponse := &fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   returnMessageString,
				ErrorCodes: nil,
			}

			returnMessage = &fenixGuiTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage{
				TestInstructionMessages:          nil,
				TestInstructionContainerMessages: nil,
				AckNackResponse:                  ackNackResponse,
			}

			return returnMessage
		}

	}

	// Do the gRPC-call
	returnMessage, err = fenixGuiBuilderServerGrpcClient.GetPinnedTestInstructionsAndTestContainers(ctx, userIdentificationMessage)

	// Shouldn't happen
	if err != nil {
		fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
			"ID":    "7bff3257-a193-4d07-83aa-f106f6f734a0",
			"error": err,
		}).Error("Problem to do gRPC-call to FenixTestGuiBuilderServer for 'SendGetPinnedTestInstructionsAndTestContainers'")

	} else if returnMessage.AckNackResponse.AckNack == false {
		// FenixTestGuiBuilderServer couldn't handle gPRC call
		fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
			"ID":                                     "72f79764-2549-4ce7-867e-16cd0f414dff",
			"Message from FenixTestGuiBuilderServer": returnMessage.AckNackResponse.Comments,
		}).Error("Problem to do gRPC-call to FenixTestGuiBuilderServer for 'SendGetPinnedTestInstructionsAndTestContainers'")
	}

	return returnMessage

}

// SendSavePinnedTestInstructionsAndTestContainers - Save pinned TestInstructions and TestInstructionContainers
func (fenixGuiBuilderProxyServerObject *fenixGuiBuilderProxyServerObjectStruct) SendSavePinnedTestInstructionsAndTestContainers(pinnedTestInstructionsAndTestContainersMessage *fenixGuiTestCaseBuilderServerGrpcApi.PinnedTestInstructionsAndTestContainersMessage) (returnMessage *fenixGuiTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage) {

	var ctx context.Context
	var returnMessageAckNack bool
	var returnMessageString string
	var err error

	// Set up connection to Server
	fenixGuiBuilderProxyServerObject.SetConnectionToFenixGuiBuilderServer()

	// Create the request message
	userIdentificationMessage := &fenixGuiTestCaseBuilderServerGrpcApi.UserIdentificationMessage{
		UserId: pinnedTestInstructionsAndTestContainersMessage.UserId,
		ProtoFileVersionUsedByClient: fenixGuiTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(
			fenixGuiBuilderProxyServerObject.getHighestFenixTestDataProtoFileVersion()),
	}

	// Do gRPC-call
	//ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		//TODO Fixa så att denna inte görs som allt går bra
		fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
			"ID": "c5ba19bd-75ff-4366-818d-745d4d7f1a52",
		}).Error("Running Defer Cancel function")
		cancel()
	}()

	// Only add access token when run on GCP
	if common_config.ExecutionLocationForFenixGuiServer == common_config.GCP {

		// Add Access token
		ctx, returnMessageAckNack, returnMessageString = fenixGuiBuilderProxyServerObject.generateGCPAccessToken(ctx)
		if returnMessageAckNack == false {
			// When error
			ackNackResponse := &fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   returnMessageString,
				ErrorCodes: nil,
			}

			returnMessage = &fenixGuiTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage{
				TestInstructionMessages:          nil,
				TestInstructionContainerMessages: nil,
				AckNackResponse:                  ackNackResponse,
			}

			return returnMessage
		}

	}

	// Do the gRPC-call
	returnMessage, err = fenixGuiBuilderServerGrpcClient.GetPinnedTestInstructionsAndTestContainers(ctx, userIdentificationMessage)

	// Shouldn't happen
	if err != nil {
		fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
			"ID":    "7bff3257-a193-4d07-83aa-f106f6f734a0",
			"error": err,
		}).Error("Problem to do gRPC-call to FenixTestGuiBuilderServer for 'SendGetPinnedTestInstructionsAndTestContainers'")

	} else if returnMessage.AckNackResponse.AckNack == false {
		// FenixTestGuiBuilderServer couldn't handle gPRC call
		fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
			"ID":                                     "72f79764-2549-4ce7-867e-16cd0f414dff",
			"Message from FenixTestGuiBuilderServer": returnMessage.AckNackResponse.Comments,
		}).Error("Problem to do gRPC-call to FenixTestGuiBuilderServer for 'SendGetPinnedTestInstructionsAndTestContainers'")
	}

	return returnMessage

}
