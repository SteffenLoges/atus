package main

import (
	"atus/backend/atus"
	"atus/backend/config"
	"atus/backend/fileserver"
	"atus/backend/logger"
	"atus/backend/release"
	"atus/backend/routes"
	"atus/backend/source"
	"atus/backend/sqlite"
	"atus/backend/user"
	"atus/backend/websocket"
	"atus/backend/websocketEvents"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func cleanup() {
	for i := 1; i <= 10; i++ {
		logger.ForceConsole().Debug("Cleaning up, please wait...")
	}

	source.DeleteTempFavicons()
	sqlite.Conn.Close()

	logger.ForceConsole().Debug("Cleanup done, bye!")
}

func main() {

	if err := sqlite.Connect(config.Base.SQLiteDSN); err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	if err := sqlite.Prepare(); err != nil {
		log.Fatalf("could not prepare database: %v", err)
	}

	atusInstance, err := atus.New()
	if err != nil {
		logger.Errorf("error while creating atus instance: %v", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cleanup()
		os.Exit(0)
	}()

	ctx := context.WithValue(context.Background(), atus.ContextKey, atusInstance)
	clientHub := websocket.NewHub(ctx)

	// -- setup event handlers ----------------------------------------------------------------------

	atusInstance.OnReleaseAdded = func(r *release.Release) {
		user.BroadcastNotification(clientHub, "NEW_RELEASE", fmt.Sprintf("Found new release on %s", r.Source.Name), r.Name)
	}

	atusInstance.OnFileserversUpdated = func(f *atus.Fileserver) {
		clientHub.MarshalAndBroadcast("FILESERVER_STATISTICS", atusInstance.GetFileserverStatistics())
	}

	getClientsForReleaseUpadte := func(rlsUID string) []*websocket.Client {
		return append(
			// browsing clients
			clientHub.GetClientsByPage("releases_browse", nil),

			// clients on the releases details page
			clientHub.GetClientsByPage("releases_details", &map[string]string{
				"uid": rlsUID,
			})...,
		)
	}

	atusInstance.OnMetaFilesUpdated = func(r *atus.Release) {
		for _, c := range getClientsForReleaseUpadte(r.UID) {
			c.MarshalAndSend("RELEASE__DETAILS__META_FILES", map[string]interface{}{
				"uid":  r.UID,
				"data": r.MetaFiles,
			})
		}
	}

	atusInstance.OnReleaseStateUpdated = func(r *atus.Release, uploadDate time.Time) {
		for _, c := range getClientsForReleaseUpadte(r.UID) {
			c.MarshalAndSend("RELEASE_DETAILS__STATE", map[string]interface{}{
				"uid": r.UID,
				"data": map[string]interface{}{
					"state":      r.State,
					"uploadDate": uploadDate,
				},
			})
		}
	}

	atusInstance.OnDownloadStateChanged = func(f *fileserver.ListFile) {
		pr := atusInstance.GetPendingReleaseByHash(f.Hash)
		if pr == nil {
			return
		}

		for _, c := range getClientsForReleaseUpadte(pr.UID) {
			c.MarshalAndSend("RELEASE__DETAILS__DOWNLOAD_STATE", map[string]interface{}{
				"uid":  pr.UID,
				"data": f,
			})
		}
	}

	logger.OnLog = func(le logger.LogEntry, severity logger.LogSeverity, message ...string) {
		for _, client := range clientHub.GetClientsByPage("debug", nil) {
			client.MarshalAndSend("DEBUG_ENTRY", le)
		}
	}

	clientHub.OnClientConnected = func(c *websocket.Client) {
		logger.Debugf("Websocket-Client connected - UUID: %s, IP: %s", c.UserUID, c.IP)
		c.MarshalAndSend("SETUP_STATUS", atus.GetSetupStatus())
		c.MarshalAndSend("FILESERVER_STATISTICS", atusInstance.GetFileserverStatistics())
	}

	// -- websocket event handlers ------------------------------------------------------------------

	// -- global ----------------------------------

	clientHub.SetEventHandler("GLOBAL__FORCE_UPDATE_FILESERVER_STATISTICS", websocketEvents.Global_ForceUpdateFileserverStatistics)

	// -- sources ---------------------------------

	// add source
	clientHub.SetEventHandler("SETTINGS__SOURCES_ADD__SET_RSS_URL", websocketEvents.Settings__SourcesAdd_SetRSSURL)
	clientHub.SetEventHandler("SETTINGS__SOURCES_ADD__SET_META_PATH", websocketEvents.Settings__SourcesAdd_SetMetaPath)
	clientHub.SetEventHandler("SETTINGS__SOURCES_ADD__SET_IMAGE_PATH", websocketEvents.Settings__SourcesAdd_SetImagePath)
	clientHub.SetEventHandler("SETTINGS__SOURCES_ADD__SET_SETTINGS", websocketEvents.Settings__SourcesAdd_SetSettings)

	// manage sources
	clientHub.SetEventHandler("SETTINGS__SOURCES_MANAGE__GET_ALL", websocketEvents.Settings__SourcesManage_GetAll)
	clientHub.SetEventHandler("SETTINGS__SOURCES_MANAGE__TOGGLE", websocketEvents.Settings__SourcesManage_Toggle)
	clientHub.SetEventHandler("SETTINGS__SOURCES_MANAGE__DELETE", websocketEvents.Settings__SourcesManage_Delete)

	// edit source
	clientHub.SetEventHandler("SETTINGS__SOURCES_EDIT__GET", websocketEvents.Settings__SourcesEdit_Get)
	clientHub.SetEventHandler("SETTINGS__SOURCES_EDIT__SAVE", websocketEvents.Settings__SourcesEdit_Save)

	// -- fileservers -----------------------------

	// add fileserver
	clientHub.SetEventHandler("SETTINGS__FILESERVERS_ADD__SET_URL", websocketEvents.Settings__FileserversAdd_SetURL)
	clientHub.SetEventHandler("SETTINGS__FILESERVERS_ADD__SET_SETTINGS", websocketEvents.Settings__FileserversAdd_SetSettings)

	// manage fileservers
	clientHub.SetEventHandler("SETTINGS__FILESERVERS_MANAGE__GET_ALL", websocketEvents.Settings__FileserversManage_GetAll)
	clientHub.SetEventHandler("SETTINGS__FILESERVERS_MANAGE__DELETE", websocketEvents.Settings__FileserversManage_Delete)
	clientHub.SetEventHandler("SETTINGS__FILESERVERS_MANAGE__TOGGLE", websocketEvents.Settings__FileserversManage_Toggle)

	// edit fileserver
	clientHub.SetEventHandler("SETTINGS__FILESERVERS_EDIT__GET", websocketEvents.Settings__FileserversEdit_Get)
	clientHub.SetEventHandler("SETTINGS__FILESERVERS_EDIT__SAVE", websocketEvents.Settings__FileserversEdit_Save)

	// fileserver settings
	clientHub.SetEventHandler("SETTINGS__FILESERVERS_SETTINGS__GET", websocketEvents.Settings__FileserversSettings_Get)
	clientHub.SetEventHandler("SETTINGS__FILESERVERS_SETTINGS__SAVE", websocketEvents.Settings__FileserversSettings_Save)

	// -- filters ---------------------------------
	clientHub.SetEventHandler("SETTINGS__FILTERS_MISC__GET_ALL", websocketEvents.Settings__FiltersMisc_GetAll)
	clientHub.SetEventHandler("SETTINGS__FILTERS_MISC__SAVE", websocketEvents.Settings__FiltersMisc_Save)
	clientHub.SetEventHandler("SETTINGS__FILTERS_CATEGORIES__GET_ALL", websocketEvents.Settings__FiltersCategories_GetAll)
	clientHub.SetEventHandler("SETTINGS__FILTERS_CATEGORIES__SAVE", websocketEvents.Settings__FiltersCategories_Save)

	// -- samples ---------------------------------
	clientHub.SetEventHandler("SETTINGS__SAMPLES_MANAGE__GET_ALL", websocketEvents.Settings__SamplesManage_GetAll)
	clientHub.SetEventHandler("SETTINGS__SAMPLES_MANAGE__SAVE", websocketEvents.Settings__SamplesManage_Save)

	// -- upload ----------------------------------
	clientHub.SetEventHandler("SETTINGS__UPLOAD__GET_ALL", websocketEvents.Settings__Upload_GetAll)
	clientHub.SetEventHandler("SETTINGS__UPLOAD__SAVE", websocketEvents.Settings__Upload_Save)
	clientHub.SetEventHandler("SETTINGS__UPLOAD__UPLOAD_TEST_TORRENT", websocketEvents.Settings__Upload_UploadTestTorrent)

	// -- users -----------------------------------
	clientHub.SetEventHandler("SETTINGS__USERS__GET", websocketEvents.Settings__Users_Get)
	clientHub.SetEventHandler("SETTINGS__USERS__GET_ALL", websocketEvents.Settings__Users_GetAll)
	clientHub.SetEventHandler("SETTINGS__USERS__UPDATE", websocketEvents.Settings__Users_Update)
	clientHub.SetEventHandler("SETTINGS__USERS__DELETE", websocketEvents.Settings__Users_Delete)
	clientHub.SetEventHandler("SETTINGS__USERS__ADD", websocketEvents.Settings__Users_Add)

	// -- api -------------------------------------
	clientHub.SetEventHandler("SETTINGS__API_MANAGE__RESET_API_TOKEN", websocketEvents.Settings__APIManage_ResetAPIToken)
	clientHub.SetEventHandler("SETTINGS__API_MANAGE__GET_TOKEN", websocketEvents.Settings__APIManage_GetToken)

	// -- Releases --------------------------------

	// browse
	clientHub.SetEventHandler("RELEASES__BROWSE__GET", websocketEvents.Releases__Browse_Get)

	// common
	clientHub.SetEventHandler("RELEASE__DETAILS__GET", websocketEvents.Release__GetDetails)
	clientHub.SetEventHandler("RELEASE__DETAILS__GET_FILES", websocketEvents.Release__GetFiles)
	clientHub.SetEventHandler("RELEASE__DELETE", websocketEvents.Release__Delete)
	clientHub.SetEventHandler("RELEASE__UPLOAD", websocketEvents.Release__Upload)

	// --- log ------------------------------------

	clientHub.SetEventHandler("LOG__GET", websocketEvents.Log__Get)
	clientHub.SetEventHandler("LOG__CLEAR", websocketEvents.Log__Clear)

	// -- debug -----------------------------------

	clientHub.SetEventHandler("DEBUG__GET_CACHE", websocketEvents.Debug__GetCache)

	// -- start the hub -----------------------------------------------------------------------------

	go clientHub.Run()

	// -- HTTP routes -------------------------------------------------------------------------------
	r := mux.NewRouter()

	// frontend
	frontendAPISR := r.PathPrefix("/frontend/api").Subrouter()
	frontendAPISR.Use(routes.MiddlewareHeaders)
	frontendAPISR.Use(routes.MiddlewareAuth)

	frontendAPISR.HandleFunc("/user/auth", routes.UserAuth).Methods("GET")
	frontendAPISR.HandleFunc("/user/login", routes.UserLogin).Methods("POST")
	frontendAPISR.HandleFunc("/user/register", routes.UserRegister).Methods("POST")
	frontendAPISR.HandleFunc("/ws", routes.SocketUserHandler(clientHub, atusInstance)).Methods("GET")

	// api
	apiSR := r.PathPrefix("/api").Subrouter()
	apiSR.Use(routes.MiddlewareHeaders)
	apiSR.Use(routes.MiddlewareAPIAuth)
	apiSR.HandleFunc("/releases", routes.API__Releases(atusInstance)).Methods("GET")
	apiSR.PathPrefix("/data/").HandlerFunc(routes.API__ServeFile)

	// catch all
	r.PathPrefix("/").HandlerFunc(routes.CatchAll)

	// ----------------------------------------------------------------------------------------------

	// add cors headers for development
	allowCredentials := handlers.AllowCredentials()
	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:3000", "http://127.0.0.1:3000"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Authorization"})

	listenAddr := config.GetEnv("LISTEN_ADDR", ":8000")
	s := &http.Server{
		Handler:           handlers.CORS(allowCredentials, allowedOrigins, allowedMethods, allowedHeaders)(r),
		Addr:              listenAddr,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	logger.ForceConsole().Debugf("ATUS is listening on %s", listenAddr)
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}

}
