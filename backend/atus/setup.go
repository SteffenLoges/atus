package atus

import (
	"atus/backend/config"
	"atus/backend/websocket"
)

type SetupStep string

const (
	SetupStepSourceAdded      SetupStep = "SOURCE_ADDED"
	SetupStepFileserverAdded  SetupStep = "FILESERVER_ADDED"
	SetupStepUploadConfigured SetupStep = "UPLOAD_CONFIGURED"
)

func GetSetupStatus() map[SetupStep]bool {
	status := make(map[SetupStep]bool)
	for _, prop := range []SetupStep{SetupStepSourceAdded, SetupStepFileserverAdded, SetupStepUploadConfigured} {
		status[prop] = config.GetBool("SETUP__" + string(prop))
	}

	return status
}

func SetSetupStepDone(clientHub *websocket.Hub, step SetupStep) {
	config.Set("SETUP__"+string(step), true)
	clientHub.MarshalAndBroadcast("SETUP_STATUS", GetSetupStatus())
}
