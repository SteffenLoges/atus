package video

import (
	"atus/backend/config"
	"fmt"
	"os/exec"
)

func getFFMPEG() string {
	if config.Base.Dependencies.FFMPEG != "" {
		return config.Base.Dependencies.FFMPEG
	}
	return "ffmpeg"
}

func getFFProbe() string {
	if config.Base.Dependencies.FFProbe != "" {
		return config.Base.Dependencies.FFProbe
	}
	return "ffprobe"
}

func getMP4Fragment() string {
	if config.Base.Dependencies.MP4Fragment != "" {
		return config.Base.Dependencies.MP4Fragment
	}
	return "mp4fragment"
}

func getMP4Dash() string {
	if config.Base.Dependencies.MP4Dash != "" {
		return config.Base.Dependencies.MP4Dash
	}
	return "mp4dash"
}

func isInstalled(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func CheckDependencies() error {
	if !isInstalled(getFFMPEG()) {
		return fmt.Errorf("ffmpeg is not installed")
	}

	if !isInstalled(getFFProbe()) {
		return fmt.Errorf("ffprobe is not installed")
	}

	if !isInstalled(getMP4Fragment()) {
		return fmt.Errorf("mp4fragment is not installed")
	}

	if !isInstalled(getMP4Dash()) {
		return fmt.Errorf("mp4dash is not installed")
	}

	return nil
}
