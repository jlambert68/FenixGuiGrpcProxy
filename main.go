package main

import (
	"FenixGuiGrpcProxyServer/common_config"
	_ "embed"
	"strconv"

	//"flag"
	"fmt"
	"log"
	"os"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
)

// Embedded resources into the binary
// The icon used
//go:embed resources/fenix_icon.png
var embededfenixIcon []byte

// mustGetEnv is a helper function for getting environment variables.
// Displays a warning if the environment variable is not set.
func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Warning: %s environment variable not set.\n", k)
	}
	return v
}

func main() {
	//time.Sleep(15 * time.Second)

	// Start up application as SysTray if environment variable says that
	if tempRunAsTrayApplication == true {
		systray.Run(onReady, onExit)
	}

	fenixGuiTestCaseBuilderServerMain()
}

func init() {
	//executionLocation := flag.String("startupType", "0", "The application should be started with one of the following: LOCALHOST_NODOCKER, LOCALHOST_DOCKER, GCP")
	//flag.Parse()

	var err error

	// Get Environment variable to tell how this program was started
	var executionLocation = mustGetenv("ExecutionLocation")

	switch executionLocation {
	case "LOCALHOST_NODOCKER":
		common_config.ExecutionLocation = common_config.LocalhostNoDocker

	case "LOCALHOST_DOCKER":
		common_config.ExecutionLocation = common_config.LocalhostDocker

	case "GCP":
		common_config.ExecutionLocation = common_config.GCP

	default:
		fmt.Println("Unknown Execution location for FenixGuiProxyServer: " + executionLocation + ". Expected one of the following: LOCALHOST_NODOCKER, LOCALHOST_DOCKER, GCP")
		os.Exit(0)

	}

	// Get Environment variable to tell how this program was started
	var executionLocationFenixGuiServer = mustGetenv("ExecutionLocationFenixGuiServer")

	switch executionLocationFenixGuiServer {
	case "LOCALHOST_NODOCKER":
		common_config.ExecutionLocationForFenixGuiServer = common_config.LocalhostNoDocker

	case "LOCALHOST_DOCKER":
		common_config.ExecutionLocationForFenixGuiServer = common_config.LocalhostDocker

	case "GCP":
		common_config.ExecutionLocationForFenixGuiServer = common_config.GCP

	default:
		fmt.Println("Unknown Execution location for FenixGuiServer: " + executionLocation + ". Expected one of the following: LOCALHOST_NODOCKER, LOCALHOST_DOCKER, GCP")
		os.Exit(0)

	}

	// Address to GuiBuilderProxyServer
	common_config.FenixGuiBuilderProxyServerAddress = mustGetenv("FenixGuiBuilderProxyServerAddress")

	// Port for GuiBuilderProxyServer
	common_config.FenixGuiBuilderProxyServerPort, err = strconv.Atoi(mustGetenv("FenixGuiBuilderProxyServerPort"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'FenixGuiBuilderProxyServerPort' to an integer, error: ", err)
		os.Exit(0)

	}

	// Address to GuiBuilderServer
	common_config.FenixGuiBuilderServerAddress = mustGetenv("FenixGuiBuilderServerAddress")

	// Port for GuiBuilderServer
	common_config.FenixGuiBuilderServerPort, err = strconv.Atoi(mustGetenv("FenixGuiBuilderServerPort"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'FenixGuiBuilderServerPort' to an integer, error: ", err)
		os.Exit(0)

	}

	// Create address for FenixGuiServer to call
	fenixGuiBuilderServerAddressToDial = common_config.FenixGuiBuilderServerAddress + ":" + strconv.Itoa(common_config.FenixGuiBuilderServerPort)

	// Get Environment variable to tell if the the application should run as a Tray Application or not
	var runAsTrayApplication = mustGetenv("RunAsTrayApplication")

	switch runAsTrayApplication {
	case "YES":
		tempRunAsTrayApplication = true

	case "NO":
		tempRunAsTrayApplication = false

	default:
		fmt.Println("Unknown value for 'RunAsTrayApplication': " + runAsTrayApplication + ". Expected one of the following: 'YES', 'NO'")
		os.Exit(0)

	}

}

// SysTray Application - StartUp
func onReady() {

	systray.SetIcon(embededfenixIcon)
	systray.SetTitle("Fenix-GUI REST -> gRPC Proxy")
	systray.SetTooltip("Fenix-GUI REST -> gRPC Proxy")
	mQuit := systray.AddMenuItem("Quit", "Quit the Fenix-GUI REST -> gRPC Proxy")

	// Sets the icon of a menu item. Only available on Mac and Windows.
	mQuit.SetIcon(icon.Data)

	// Run menu handles as go-routine
	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()

}

// SysTray Application - On exit
func onExit() {
	// clean up here, and exit the program
	os.Exit(0)

}
