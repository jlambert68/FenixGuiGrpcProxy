package main

import (
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/api/idtoken"
	grpcMetadata "google.golang.org/grpc/metadata"
	"time"
)

// Generate Google access token. Used when running in GCP
func (fenixGuiBuilderProxyServerObject *fenixGuiBuilderProxyServerObjectStruct) generateGCPAccessToken(ctx context.Context) (appendedCtx context.Context, returnAckNack bool, returnMessage string) {

	// Only create the token if there is none, or it has expired
	if fenixGuiBuilderProxyServerObject.gcpAccessToken == nil || fenixGuiBuilderProxyServerObject.gcpAccessToken.Expiry.Before(time.Now()) {

		// Create an identity token.
		// With a global TokenSource tokens would be reused and auto-refreshed at need.
		// A given TokenSource is specific to the audience.
		/*
			tokenSource, err := idtoken.NewTokenSource(ctx, "https://"+common_config.FenixGuiBuilderServerAddress)
			if err != nil {
				fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
					"ID":  "8ba622d8-b4cd-46c7-9f81-d9ade2568eca",
					"err": err,
				}).Error("Couldn't generate access token")

				return nil, false, "Couldn't generate access token"
			}

			token, err := tokenSource.Token()
		*/
		/*
			var eMailAndPrivateKey = struct {
				Email      string `json:"client_email"`
				PrivateKey string `json:"private_key"`
			}{}
			json.Unmarshal(serviceAccountKeyJson, &eMailAndPrivateKey)
			config := &jwt.Config{
				Email:      eMailAndPrivateKey.Email,
				PrivateKey: []byte(eMailAndPrivateKey.PrivateKey),
				Scopes: []string{
					gcp_scope,
				},
				TokenURL:   google.JWTTokenURL,
				UseIDToken: false,
			}

		*/

		tokenSource, err := idtoken.NewTokenSource(ctx, gcp_scope, idtoken.WithCredentialsJSON(serviceAccountKeyJson))

		if err != nil {
			fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
				"ID":  "8ba622d8-b4cd-46c7-9f81-d9ade2568eca",
				"err": err,
			}).Error("Couldn't generate access token")

			return nil, false, "Couldn't generate access token"
		}

		token, err := tokenSource.Token()
		//token, err := config.TokenSource(oauth2.NoContext).Token()

		if err != nil {
			fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
				"ID":  "0cf31da5-9e6b-41bc-96f1-6b78fb446194",
				"err": err,
			}).Error("Problem getting the token")

			return nil, false, "Problem getting the token"
		} else {
			fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
				"ID":    "8b1ca089-0797-4ee6-bf9d-f9b06f606ae9",
				"token": "Nothing to see", //token,
			}).Debug("Got Bearer Token")
		}

		fenixGuiBuilderProxyServerObject.gcpAccessToken = token

	}

	fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
		"ID": "cd124ca3-87bb-431b-9e7f-e044c52b4960",
		"fenixGuiBuilderProxyServerObject.gcpAccessToken": "Nothing to see", //fenixGuiBuilderProxyServerObject.gcpAccessToken,
	}).Debug("Will use Bearer Token")

	// Add token to gRPC Request.
	appendedCtx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+fenixGuiBuilderProxyServerObject.gcpAccessToken.AccessToken)

	return appendedCtx, true, ""

}

// ********************************************************************************************************************
// Check if Calling Client is using correct proto-file version
func (fenixGuiBuilderProxyServerObject *fenixGuiBuilderProxyServerObjectStruct) isClientUsingCorrectTestDataProtoFileVersion(callingClientUuid string, usedProtoFileVersion fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum) (returnMessage *fenixTestCaseBuilderServerGrpcApi.AckNackResponse) {

	var clientUseCorrectProtoFileVersion bool
	var protoFileExpected fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum
	var protoFileUsed fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum

	protoFileUsed = usedProtoFileVersion
	protoFileExpected = fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(fenixGuiBuilderProxyServerObject.getHighestFenixTestDataProtoFileVersion())

	// Check if correct proto files is used
	if protoFileExpected == protoFileUsed {
		clientUseCorrectProtoFileVersion = true
	} else {
		clientUseCorrectProtoFileVersion = false
	}

	// Check if Client is using correct proto files version
	if clientUseCorrectProtoFileVersion == false {
		// Not correct proto-file version is used

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_WRONG_PROTO_FILE_VERSION
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    false,
			Comments:   "Wrong proto file used. Expected: '" + protoFileExpected.String() + "', but got: '" + protoFileUsed.String() + "'",
			ErrorCodes: errorCodes,
		}

		fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
			"id": "513dd8fb-a0bb-4738-9a0b-b7eaf7bb8adb",
		}).Debug("Wrong proto file used. Expected: '" + protoFileExpected.String() + "', but got: '" + protoFileUsed.String() + "' for Client: " + callingClientUuid)

		return returnMessage

	} else {
		return nil
	}

}

// ********************************************************************************************************************
// Get the highest FenixProtoFileVersionEnumeration
func (fenixGuiBuilderProxyServerObject *fenixGuiBuilderProxyServerObjectStruct) getHighestFenixTestDataProtoFileVersion() int32 {

	// Check if there already is a 'highestFenixProtoFileVersion' saved, if so use that one
	if highestFenixProtoFileVersion != -1 {
		return highestFenixProtoFileVersion
	}

	// Find the highest value for proto-file version
	var maxValue int32
	maxValue = 0

	for _, v := range fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum_value {
		if v > maxValue {
			maxValue = v
		}
	}

	highestFenixProtoFileVersion = maxValue

	return highestFenixProtoFileVersion
}
