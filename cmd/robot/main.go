package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"robo-apocalypse/pkg/survivor"
	"robo-apocalypse/pkg/survivordb"
	"syscall"

	goflags "flag"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-openapi/runtime/middleware"
)

// rootCmd application command object
var rootCmd = &cobra.Command{
	Use:   "apocalypse",
	Short: "Track survivors in a robot apocalypse",
	Long:  "Track survivors resources, location and status of infection in a robot apocalypse",
	Run:   run,
}

// init set configuration defaults
func init() {
	rootCmd.Flags().AddGoFlagSet(goflags.CommandLine)

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("port", "8080", "Port to listen on")
	rootCmd.PersistentFlags().String("host", "", "Host IP to listen on. If the host is empty it will listen on all IPs")
	rootCmd.PersistentFlags().String("dbName",
		"./apocalypse.db", "Apocalypse statistics database")
	rootCmd.PersistentFlags().String("webTemplate",
		"index.tmpl", "HTML web template")
	rootCmd.PersistentFlags().String("styleSheet",
		"./style.css", "Web cascading style sheet")
	rootCmd.PersistentFlags().String("destEndpoint",
		"https://robotstakeover20210903110417.azurewebsites.net/robotcpu", "endpoint for the robot CPU system")
}

func initConfig() {
	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		logrus.Error(err, "viper.BindPFlags")
	}

	viper.AutomaticEnv()
	viper.AddConfigPath(".")

	viper.SetConfigName("apocalypse")

	if err := viper.ReadInConfig(); err == nil {
		logrus.WithFields(logrus.Fields{
			"file": viper.ConfigFileUsed(),
		}).Info("viper.ReadInConfig.")
	} else {
		logrus.Error(err, "viper.ReadInConfig failed")
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	loglevel := logrus.Level(viper.GetInt("loglevel"))

	logrus.WithFields(logrus.Fields{"loglevel": loglevel}).Info("Logging config.")

	logrus.SetLevel(loglevel)
	logrus.SetReportCaller(true)

	for _, v := range viper.AllKeys() {
		logrus.WithFields(logrus.Fields{
			v: viper.Get(v),
		}).Info("Configs loaded")
	}
}

// fileServe object to serve stylesheets
type fileServe struct {
	http.FileSystem
	filename string
}

// Open open a file to serve as a stylesheet
func (fs fileServe) Open(name string) (http.File, error) {
	return fs.FileSystem.Open(fs.filename)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err, "Error starting rootCmd.Execute()")
		os.Exit(1)
	}
}

// CatchCtrlC function performs a graceful shutdown
func catchCtrlC(srv *http.Server) {
	sigint := make(chan os.Signal, 1)

	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint
	logrus.Info("We received an interrupt signal, gracefully shutting down")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.WithFields(logrus.Fields{"Error": err}).Info("Server shutdown error")
	}
}

func run(cmd *cobra.Command, args []string) {
	robo := &survivor.Apocalypse{}
	robo.DB = survivordb.Open(viper.GetString("dbName"))
	defer robo.DB.DB.Close()
	if robo.DB == nil {
		return
	}
	err := robo.DB.Setup()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
		}).Info("Error setting up database")
		return
	}
	templ := template.New("").Funcs(survivor.TemplateFuncs)
	robo.HTMLTemplateName = viper.GetString("webTemplate")
	robo.HTMLTemplate, err = templ.ParseFiles(robo.HTMLTemplateName)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
		}).Info("Error parsing the web template")
		return
	}

	mux := http.NewServeMux()
	mux.Handle(viper.GetString("styleSheet"),
		http.FileServer(fileServe{
			FileSystem: http.Dir("."),
			filename:   "style.css",
		}))
	mux.HandleFunc("/", robo.DefaultPath)
	mux.HandleFunc("/survivors", robo.Survivor)
	mux.HandleFunc("/survivors/stats", robo.SurvivorStats)
	mux.HandleFunc("/survivors/location", robo.UpdateLocation)
	mux.HandleFunc("/survivors/infected", robo.Infected)
	mux.HandleFunc("/survivors/resources", robo.UpdateResources)
	mux.HandleFunc("/robotcpu", robo.RobotCPU)
	mux.HandleFunc("/reportweb", robo.Report)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	mux.Handle("/docs", sh)
	mux.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	svr := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", viper.GetString("host"), viper.GetString("port")),
		Handler: mux,
	}

	go catchCtrlC(svr)

	if err := svr.ListenAndServe(); err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
		}).Info("HTTP Server shutdown response")
		return
	}
}
