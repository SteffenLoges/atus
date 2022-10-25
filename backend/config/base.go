package config

import (
	"atus/backend/helpers"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

var BaseFolder = GetEnv("BASE_FOLDER", "./atus_data")

var Base *baseConfig

type baseConfigFolders struct {
	WWW  string `json:"www"`
	Data string `json:"data"`
}

type baseConfigAuth struct {
	JWTSecret        string        `json:"jwt_secret"`
	JWTTTL           time.Duration `json:"jwt_ttl"`
	BCryptCost       int           `json:"bcrypt_cost"`
	MaxLoginAttempts int           `json:"max_login_attempts"`
	LockoutDuration  time.Duration `json:"lockout_duration"`
}

type baseConfigRequests struct {
	InsecureSkipVerify bool `json:"insecure_skip_verify"`
}

type baseConfigDependencies struct {
	FFMPEG      string `json:"ffmpeg"`
	FFProbe     string `json:"ffprobe"`
	MP4Fragment string `json:"mp4fragment"`
	MP4Dash     string `json:"mp4dash"`
}

type baseConfigSchedulers struct {
	ProcessPendingReleasesInterval time.Duration `json:"process_pending_releases_interval"`
	FileserverGetMetaFilesInterval time.Duration `json:"fileserver_get_meta_files_interval"`
}
type baseConfig struct {
	SQLiteDSN    string                  `json:"sqlite_dsn"`
	Folders      *baseConfigFolders      `json:"folders"`
	Auth         *baseConfigAuth         `json:"jwt"`
	Requests     *baseConfigRequests     `json:"requests"`
	Dependencies *baseConfigDependencies `json:"dependencies"`
	Schedulers   *baseConfigSchedulers   `json:"schedulers"`
}

func getBaseConfigDefaults() (*baseConfig, error) {

	jwtSecret, err := helpers.GenerateRandomString(32)
	if err != nil {
		return nil, err
	}

	return &baseConfig{
		SQLiteDSN: "./atus_data/db.sqlite3?cache=shared&mode=rwc",
		Folders: &baseConfigFolders{
			WWW:  "./www",
			Data: "./atus_data/data",
		},
		Auth: &baseConfigAuth{
			JWTSecret:        jwtSecret,
			JWTTTL:           time.Hour * 24 * 7,
			MaxLoginAttempts: 5,
			LockoutDuration:  time.Minute * 5,
		},
		Requests: &baseConfigRequests{
			InsecureSkipVerify: true,
		},
		Dependencies: &baseConfigDependencies{
			FFMPEG:      "ffmpeg",
			FFProbe:     "ffprobe",
			MP4Fragment: "mp4fragment",
			MP4Dash:     "mp4dash",
		},
		Schedulers: &baseConfigSchedulers{
			ProcessPendingReleasesInterval: time.Second * 5,
			FileserverGetMetaFilesInterval: time.Second * 20,
		},
	}, nil
}

func getBaseConfig() (*baseConfig, error) {

	configFile := path.Join(BaseFolder, "config.json")

	var baseConfig *baseConfig

	// check if config file exists, create with defaults if it doesn't
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Println("config file does not exist, creating with defaults")

		baseConfig, err = getBaseConfigDefaults()
		if err != nil {
			return nil, err
		}

		// save new config file
		b, err := json.MarshalIndent(baseConfig, "", "  ")
		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(configFile, b, 0644)
		if err != nil {
			return nil, err
		}

		return baseConfig, nil
	}

	// read config file
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &baseConfig)
	if err != nil {
		return nil, err
	}

	return baseConfig, nil

}

func init() {

	// check if base folder exists, if not create it
	if _, err := os.Stat(BaseFolder); os.IsNotExist(err) {
		if err := os.Mkdir(BaseFolder, 0755); err != nil {
			log.Fatalf("error creating base folder: %s", err.Error())
		}
	}

	baseConfig, err := getBaseConfig()
	if err != nil {
		log.Fatalf("error getting base config: %s", err.Error())
	}

	Base = baseConfig

	// create data folder if it doesn't exist
	if _, err := os.Stat(Base.Folders.Data); os.IsNotExist(err) {
		if err := os.MkdirAll(Base.Folders.Data, 0755); err != nil {
			log.Fatalf("could not create data folder: %s", err)
		}
	}

}
