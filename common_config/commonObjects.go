package common_config

// ***********************************************************************************************************
// The following variables receives their values from environment variables

// Where is the client running
var ExecutionLocation ExecutionLocationTypeType

// Where is the Fenix TestDataSync server running
// LocationForFenixTestDataServer
var ExecutionLocationForFenixGuiServer ExecutionLocationTypeType

// Definitions for where client and Fenix Server is running
type ExecutionLocationTypeType int

// Constants used for where stuff is running
const (
	LocalhostNoDocker ExecutionLocationTypeType = iota
	LocalhostDocker
	GCP
)

// FenixGuiBuilderProxyServer
var LocationForFenixGuiBuilderProxyServerTypeMapping = map[ExecutionLocationTypeType]string{
	LocalhostNoDocker: "LOCALHOST_NODOCKER",
	LocalhostDocker:   "LOCALHOST_DOCKER",
	GCP:               "GCP",
}

// Address to Fenix TestData Server & Client, will have their values from Environment variables at startup
var (
	FenixGuiBuilderServerAddress      string //
	FenixGuiBuilderServerPort         int    // TODO remove,
	FenixGuiBuilderProxyServerAddress string
	FenixGuiBuilderProxyServerPort    int
)

// ***********************************************************************************************************
