package survivor

import (
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"robo-apocalypse/pkg/survivordb"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var TemplateFuncs = template.FuncMap{"rangeStruct": rangeStructer}

// Tracker structure of a Tracker object
type Apocalypse struct {
	DB               *survivordb.SurvivorDB
	HTMLTemplate     *template.Template
	HTMLTemplateName string
}

// DefaultPath endpoint to the default path
func (a *Apocalypse) DefaultPath(w http.ResponseWriter, r *http.Request) {
	logrus.WithFields(logrus.Fields{
		"EndPoint:": r.URL.Path,
	}).Info("Apocalypse.DefaultPath")
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"URL":   r.URL.Path,
		}).Info("Apocalypse.DefaultPath, ioutil.ReadAll")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	logrus.WithFields(logrus.Fields{
		"body": string(body),
	}).Info("Apocalypse.DefaultPath")

	w.WriteHeader(http.StatusNotFound)
}

// swagger:route GET /survivors/stats survivors getStats
// Return the statistics of infected survivors from the database
// responses:
//	200: statsResponse

// SurvivorStats handles GET requests and returns infected survivors statistics
func (a *Apocalypse) SurvivorStats(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Apocalypse.SurvivorStats")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	healthyCount := a.DB.CountSurvivors(false)
	if healthyCount == -1 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	infectedCount := a.DB.CountSurvivors(true)
	if infectedCount == -1 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var total float64 = float64(healthyCount + infectedCount)

	var healthyPercentage, infectedPercentage float64 = 0, 0

	if total > 0 {
		healthyPercentage = float64(healthyCount) / total * 100.0
		infectedPercentage = float64(infectedCount) / total * 100
	}
	stats := &struct {
		HealthyPercentage  float64 `json:"healthyPercentage"`
		InfectedPercentage float64 `json:"infectedPercentage"`
	}{HealthyPercentage: healthyPercentage, InfectedPercentage: infectedPercentage}

	statsBuffer, err := json.Marshal(stats)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"body":  stats,
			"Error": err,
		}).Error("Marshal")
		return
	}

	logrus.WithFields(logrus.Fields{
		"body": string(statsBuffer),
	}).Info("Data")

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Write(statsBuffer)
}

func (a *Apocalypse) listSurvivors(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Apocalypse.SurvivorStats")

	survivors := a.DB.GetAllSurvivors()

	survivorsBuffer, err := json.Marshal(survivors)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"body":  survivors,
			"Error": err,
		}).Error("Marshal")
		return
	}

	logrus.WithFields(logrus.Fields{
		"body": string(survivorsBuffer),
	}).Info("Data")

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Write(survivorsBuffer)
}

// newSurvivor endpoint to Apocalypse
func (a *Apocalypse) newSurvivor(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Apocalypse.newSurvivor")
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
		}).Info("Error reading response")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	survivor := &survivordb.Survivor{}
	if err := json.Unmarshal(body, survivor); err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"body":  string(body),
		}).Info("Error unmarshalling")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logrus.WithFields(logrus.Fields{
		"body": survivor,
	}).Info("Incoming")
	err = a.DB.Save(survivor)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"body":  survivor,
		}).Info("Error saving")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /survivors survivors getSurvivors
// Return a list of survivors from the database
// responses:
//	200: surivivorsResponse

// Survivor handles GET requests and returns all survivors

// swagger:route POST /survivors survivors createSurvivor
// Create a new Survivor
//
// responses:
//	200:
//  	500:

// Survivor handles POST requests to add new survivor
func (a *Apocalypse) Survivor(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Apocalypse.Survivor")
	switch r.Method {
	case http.MethodGet:
		a.listSurvivors(w, r)
	case http.MethodPost:
		a.newSurvivor(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// swagger:route PUT /survivors/location survivors updateLocation
// Return the HTTP response code: 200, 404, 500
// responses:
//	200:
//	404:
//	500:

// UpdateLocation handles PUT requests and returns an HTTP response code
func (a *Apocalypse) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Apocalypse.UpdateLocation")
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
		}).Info("Error reading response")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	locationPayload := &struct {
		IdNumber string `json:"id"`
		survivordb.LastLocation
	}{}
	if err := json.Unmarshal(body, locationPayload); err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"body":  string(body),
		}).Info("Error unmarshalling")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logrus.WithFields(logrus.Fields{
		"body": locationPayload,
	}).Info("Incoming")
	err = a.DB.UpdateLocation(locationPayload.IdNumber, locationPayload.Longitude, locationPayload.Latitude)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":  err,
			"action": locationPayload,
		}).Info("Error saving")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// updateInfected endpoint to update a survivor location
func (a *Apocalypse) updateInfected(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Apocalypse.updateInfected")

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
		}).Info("Error reading response")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	infectedPayload := &struct {
		IdNumber string `json:"id"`
	}{}
	if err := json.Unmarshal(body, infectedPayload); err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"body":  string(body),
		}).Info("Error unmarshalling")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logrus.WithFields(logrus.Fields{
		"body": infectedPayload,
	}).Info("Incoming")
	err = a.DB.UpdateInfected(infectedPayload.IdNumber)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":  err,
			"action": infectedPayload,
		}).Info("Error saving")
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// listSurvivors endpoint to Apocalypse
func (a *Apocalypse) listInfected(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Apocalypse.SurvivorStats")

	query := r.URL.Query()
	statusParameter := query.Get("status")
	status := true
	switch statusParameter {
	case "true":
		status = true
	case "false":
		status = false
	}
	infected := a.DB.GetSurvivors(status)

	infectedBuffer, err := json.Marshal(infected)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"body":  infectedBuffer,
			"Error": err,
		}).Error("Marshal")
		return
	}

	logrus.WithFields(logrus.Fields{
		"body":   string(infectedBuffer),
		"status": status,
	}).Info("Data")

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Write(infectedBuffer)
}

// swagger:parameters getInfected
type InfectedStatusParam struct {
	//
	// min items: 1
	// max items: 1
	// unique: true
	// in: query
	// example: status=true
	Status string `json:"status"`
}

// swagger:route GET /survivors/infected survivors getInfected
// Return a list of infected survivors from the database
// responses:
//	200: surivivorsResponse

// Infected handles GET requests and returns infected survivors

// swagger:route PUT /survivors/infected survivors setInfected
// Return the HTTP response code: 200, 404, 500
// responses:
//	200:
//	404:
//	500:

// Infected handles PUT requests and returns an HTTP response code
func (a *Apocalypse) Infected(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Apocalypse.Infected")
	switch r.Method {
	case http.MethodGet:
		a.listInfected(w, r)
	case http.MethodPut:
		a.updateInfected(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// swagger:route PUT /survivors/resource survivors updateResource
// Return the HTTP response code: 200, 404, 500
// responses:
//	200:

// UpdateResources handles PUT requests and returns an HTTP response code
func (a *Apocalypse) UpdateResources(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Apocalypse.UpdateResources")
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
		}).Info("Error reading response")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resourcePayload := &struct {
		IdNumber string `json:"id"`
		survivordb.Resources
	}{}
	if err := json.Unmarshal(body, resourcePayload); err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"body":  string(body),
		}).Info("Error unmarshalling")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logrus.WithFields(logrus.Fields{
		"body": resourcePayload,
	}).Info("Incoming")
	err = a.DB.UpdateResource(resourcePayload.IdNumber,
		resourcePayload.Water,
		resourcePayload.Food,
		resourcePayload.Medication,
		resourcePayload.Ammunition,
	)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"body":  resourcePayload,
		}).Info("Error saving")
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// Report endpoint to Apocalypse
func (a *Apocalypse) Report(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Apocalypse.Report")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	survivors := a.DB.GetAllSurvivors()

	err := a.HTMLTemplate.ExecuteTemplate(w, a.HTMLTemplateName, survivors)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
		}).Info("Error writing response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type RobotCpu struct {
	Model            string    `json:"model"`
	SerialNumber     string    `json:"serialNumber"`
	ManufacturedDate time.Time `json:"manufacturedDate"`
	Category         string    `json:"category"`
}

type RobotCpuSorter struct {
	robots []RobotCpu
	by     func(p1, p2 *RobotCpu) bool // Closure used in the Less method.
}

func (c *RobotCpuSorter) Len() int {
	return len(c.robots)
}

// Swap is part of sort.Interface.
func (c *RobotCpuSorter) Swap(i, j int) {
	c.robots[i], c.robots[j] = c.robots[j], c.robots[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (c *RobotCpuSorter) Less(i, j int) bool {
	return c.by(&c.robots[i], &c.robots[j])
}

// A list of robotcpus
// swagger:response robotcpuResponse
type robotcpuResponseWrapper struct {
	// All current robotcpus
	// in: body
	Body []RobotCpu
}

// swagger:parameters getRobotCPU
type RobortSortCPUParam struct {
	// a BarSlice has bars which are strings
	//
	// min items: 1
	// max items: 1
	// unique: true
	// in: query
	// example: sortby=category
	Sortby string `json:"sortby"`
}

// swagger:parameters getRobotCPU
type RobortCategoryCPUParam struct {
	// a BarSlice has bars which are strings
	//
	// min items: 1
	// max items: 1
	// unique: true
	// in: query
	// example: category=Flying
	Category string `json:"category"`
}

// swagger:route GET /survivors/infected survivors getRobotCPU
// Returns a list of infected survivors from the database
// responses:
//	200: robotcpuResponse

// RobotCPU handles GET requests and returns robotCPUs
func (a *Apocalypse) RobotCPU(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Apocalypse.RobotCPU")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	enpoint := viper.GetString("destEndpoint")
	resp, err := http.Get(enpoint)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":    err,
			"Endpoint": enpoint,
		}).Info("Error GET")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
		}).Info("Error reading response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	robotcpus := &RobotCpuSorter{}

	if err := json.Unmarshal(body, &robotcpus.robots); err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"body":  string(body),
		}).Info("Error unmarshalling")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := r.URL.Query()
	category := query.Get("category")
	switch category {
	case "Flying":
		var tmp []RobotCpu
		for _, robot := range robotcpus.robots {
			if robot.Category == "Flying" {
				tmp = append(tmp, robot)
			}
		}
		robotcpus.robots = tmp
	case "Land":
		var tmp []RobotCpu
		for _, robot := range robotcpus.robots {
			if robot.Category == "Land" {
				tmp = append(tmp, robot)
			}
		}
		robotcpus.robots = tmp
	}
	sortColumn := query.Get("sortby")
	switch sortColumn {
	case "model":
		cmp := func(r1, r2 *RobotCpu) bool {
			return r1.Model < r2.Model
		}
		robotcpus.by = cmp
		sort.Sort(robotcpus)
	case "serialNumber":
		cmp := func(r1, r2 *RobotCpu) bool {
			return r1.SerialNumber < r2.SerialNumber
		}
		robotcpus.by = cmp
		sort.Sort(robotcpus)
	case "manufacturedDate":
		cmp := func(r1, r2 *RobotCpu) bool {
			return r1.ManufacturedDate.Before(r2.ManufacturedDate)
		}
		robotcpus.by = cmp
		sort.Sort(robotcpus)
	case "category":
		cmp := func(r1, r2 *RobotCpu) bool {
			return r1.Category < r2.Category
		}
		robotcpus.by = cmp
		sort.Sort(robotcpus)
	}

	robotsBuffer, err := json.Marshal(robotcpus.robots)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"body":  robotcpus.robots,
			"Error": err,
		}).Error("Marshal")
		return
	}

	logrus.WithFields(logrus.Fields{
		"body":        string(robotsBuffer),
		"sort column": sortColumn,
		"category":    category,
	}).Info("Data")

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Write(robotsBuffer)
}

// rangeStructer takes the first argument, which must be a struct, and
// returns the value of each field in a slice. It will return nil
// if there are no arguments or first argument is not a struct
func rangeStructer(args ...interface{}) []interface{} {
	if len(args) == 0 {
		return nil
	}

	v := reflect.ValueOf(args[0])
	if v.Kind() != reflect.Struct {
		return nil
	}

	out := make([]interface{}, 0)
	for i := 0; i < v.NumField(); i++ {
		switch u := v.Field(i).Interface().(type) {
		case bool, string, float64, int:
			out = append(out, u)
			continue
		case time.Time:
			out = append(out, u.String())
			continue
		}
		if v.Field(i).Kind() == reflect.Struct {
			v2 := v.Field(i)
			for j := 0; j < v2.NumField(); j++ {
				switch w := v2.Field(j).Interface().(type) {
				case bool, string, float64, int:
					out = append(out, w)
				case time.Time:
					out = append(out, w.String())
				}
			}
			continue
		}
	}

	return out
}
