package survivor

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"robo-apocalypse/pkg/survivordb"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

var survivorRequest string = `{
	"name":		"Jane Doe",
	"age":		1,
	"gender":	"Female",
	"id":		"HD138VOP34219",
	"longitude":	0,
	"latitude":	0,
	"water":	2000,
	"food":		"Beef, Fish, Milk, Pasta, Ceriel",
	"medication":	"Antibiotics, venteze cfc free, cough syrup",
	"infected":	false,
	"timestamp":	"2022-03-11T08:19:35Z"
}`

var updaterLocationRequest string = `{
	"id":		"HD138VOP34219",
	"longitude":	1,
	"latitude":	2
}`

var updaterResourceRequest string = `{
	"id":		"HD138VOP34219",
	"water":	0,
	"food":		"Fish, Bread",
	"medication":	"Antibiotics"
}`

var updaterInfectedRequest string = `{
	"id":		"HD138VOP34219"
}`

// TestApocalypseApi_NewSurvivor checks if the api endpoint
// returns a success http status
func TestApocalypseApi_NewSurvivor(t *testing.T) {
	robo := &Apocalypse{}
	os.Remove("./test.db")
	robo.DB = survivordb.Open("./test.db")
	if robo.DB == nil {
		return
	}
	err := robo.DB.Setup()
	if err != nil {
		t.Errorf("Error setting up database: %v", err)
		return
	}

	reader := strings.NewReader(survivorRequest)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/survivors", reader)
	robo.Survivor(w, r)
	resp := w.Result()
	if resp.Status != fmt.Sprintf("%d OK", http.StatusOK) {
		t.Errorf("Apocalypse.NewSurvivor(w http.ResponseWriter, r *http.Request): want: %v, got: %v", http.StatusOK, resp.Status)
	}
	survivor := &survivordb.Survivor{}
	if err := json.Unmarshal([]byte(survivorRequest), survivor); err != nil {
		t.Errorf("Apocalypse.NewSurvivor(w http.ResponseWriter, r *http.Request): could not json.Unmarshal: %v", survivorRequest)
	}
	logrus.WithFields(logrus.Fields{
		"survivor": survivor,
	}).Info("TestApocalypseApi_NewSurvivor info")
	newSurvivor := robo.DB.GetSurvivor(survivor.IdNumber)
	if newSurvivor == nil {
		t.Errorf("SurvivorDB.GetSurvivor() - %q: want: not nil, got: %v", survivor.Name, newSurvivor)
	}
}

// TestApocalypseApi_UpdateLocation checks if the api endpoint
// returns a success http status
func TestApocalypseApi_UpdateLocation(t *testing.T) {
	robo := &Apocalypse{}
	os.Remove("./test.db")
	robo.DB = survivordb.Open("./test.db")
	if robo.DB == nil {
		return
	}
	err := robo.DB.Setup()
	if err != nil {
		t.Errorf("Error setting up database: %v", err)
		return
	}

	addreader := strings.NewReader(survivorRequest)
	addw := httptest.NewRecorder()
	addr := httptest.NewRequest(http.MethodPost, "/survivors/add", addreader)
	robo.Survivor(addw, addr)
	addresp := addw.Result()
	if addresp.Status != fmt.Sprintf("%d OK", http.StatusOK) {
		t.Errorf("Apocalypse.UpdateLocation(w http.ResponseWriter, r *http.Request): want: %v, got: %v", http.StatusOK, addresp.Status)
	}

	reader := strings.NewReader(updaterLocationRequest)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/survivors/location", reader)
	robo.UpdateLocation(w, r)
	resp := w.Result()
	if resp.Status != fmt.Sprintf("%d OK", http.StatusOK) {
		t.Errorf("Apocalypse.UpdateLocation(w http.ResponseWriter, r *http.Request): want: %v, got: %v", http.StatusOK, resp.Status)
	}
	locationPayload := &struct {
		IdNumber string `json:"id"`
		survivordb.LastLocation
	}{}
	if err := json.Unmarshal([]byte(updaterLocationRequest), locationPayload); err != nil {
		t.Errorf("Apocalypse.NewSurvivor(w http.ResponseWriter, r *http.Request): could not json.Unmarshal: %v", survivorRequest)
	}
	logrus.WithFields(logrus.Fields{
		"survivorRequest": updaterLocationRequest,
		"locationPayload": locationPayload,
	}).Info("TestApocalypseApi_UpdateLocation info")
	newSurvivor := robo.DB.GetSurvivor(locationPayload.IdNumber)
	if newSurvivor == nil || (newSurvivor.LastLocation.Longitude != 1 || newSurvivor.LastLocation.Latitude != 2) {
		t.Errorf("SurvivorDB.GetSurvivor() - %v: want: not nil, got: %v", locationPayload.IdNumber, newSurvivor)
	}
}

// TestApocalypseApi_UpdateLocation checks if the api endpoint
// returns a success http status
func TestApocalypseApi_UpdateInfected(t *testing.T) {
	robo := &Apocalypse{}
	os.Remove("./test.db")
	robo.DB = survivordb.Open("./test.db")
	if robo.DB == nil {
		return
	}
	err := robo.DB.Setup()
	if err != nil {
		t.Errorf("Error setting up database: %v", err)
		return
	}

	addreader := strings.NewReader(survivorRequest)
	addw := httptest.NewRecorder()
	addr := httptest.NewRequest(http.MethodPost, "/survivors/add", addreader)
	robo.Survivor(addw, addr)
	addresp := addw.Result()
	if addresp.Status != fmt.Sprintf("%d OK", http.StatusOK) {
		t.Errorf("Apocalypse.UpdateLocation(w http.ResponseWriter, r *http.Request): want: %v, got: %v", http.StatusOK, addresp.Status)
	}

	reader := strings.NewReader(updaterInfectedRequest)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/survivors/infected", reader)
	robo.Infected(w, r)
	resp := w.Result()
	if resp.Status != fmt.Sprintf("%d OK", http.StatusOK) {
		t.Errorf("Apocalypse.UpdateInfected(w http.ResponseWriter, r *http.Request): want: %v, got: %v", http.StatusOK, resp.Status)
	}
	locationPayload := &struct {
		IdNumber string `json:"id"`
		survivordb.LastLocation
	}{}
	if err := json.Unmarshal([]byte(updaterLocationRequest), locationPayload); err != nil {
		t.Errorf("Apocalypse.NewSurvivor(w http.ResponseWriter, r *http.Request): could not json.Unmarshal: %v", survivorRequest)
	}
	logrus.WithFields(logrus.Fields{
		"survivorRequest": updaterLocationRequest,
		"locationPayload": locationPayload,
	}).Info("TestApocalypseApi_UpdateLocation info")
	newSurvivor := robo.DB.GetSurvivor(locationPayload.IdNumber)
	if newSurvivor == nil || (newSurvivor.Infected == false) {
		t.Errorf("SurvivorDB.GetSurvivor() - %v: want: not nil, got: %v", locationPayload.IdNumber, newSurvivor)
	}
}

// TestApocalypseApi_UpdateResources checks if the api endpoint
// returns a success http status
func TestApocalypseApi_UpdateResources(t *testing.T) {
	robo := &Apocalypse{}
	os.Remove("./test.db")
	robo.DB = survivordb.Open("./test.db")
	if robo.DB == nil {
		return
	}
	err := robo.DB.Setup()
	if err != nil {
		t.Errorf("Error setting up database: %v", err)
		return
	}

	addreader := strings.NewReader(survivorRequest)
	addw := httptest.NewRecorder()
	addr := httptest.NewRequest(http.MethodPost, "/survivors/add", addreader)
	robo.Survivor(addw, addr)
	addresp := addw.Result()
	if addresp.Status != fmt.Sprintf("%d OK", http.StatusOK) {
		t.Errorf("Apocalypse.UpdateResources(w http.ResponseWriter, r *http.Request): want: %v, got: %v", http.StatusOK, addresp.Status)
	}

	reader := strings.NewReader(updaterResourceRequest)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/survivors/resources", reader)
	robo.UpdateResources(w, r)
	resp := w.Result()
	if resp.Status != fmt.Sprintf("%d OK", http.StatusOK) {
		t.Errorf("Apocalypse.UpdateResources(w http.ResponseWriter, r *http.Request): want: %v, got: %v", http.StatusOK, resp.Status)
	}
	resourcePayload := &struct {
		IdNumber string `json:"id"`
		survivordb.Resources
	}{}
	if err := json.Unmarshal([]byte(updaterLocationRequest), resourcePayload); err != nil {
		t.Errorf("Apocalypse.NewSurvivor(w http.ResponseWriter, r *http.Request): could not json.Unmarshal: %v", survivorRequest)
	}
	logrus.WithFields(logrus.Fields{
		"survivorRequest": updaterResourceRequest,
		"locationPayload": resourcePayload,
	}).Info("TestApocalypseApi_UpdateLocation info")
	newSurvivor := robo.DB.GetSurvivor(resourcePayload.IdNumber)
	if newSurvivor == nil || (newSurvivor.Resources.Water != 0) {
		t.Errorf("SurvivorDB.GetSurvivor() - %v: want: not nil, got: %v", resourcePayload.IdNumber, newSurvivor)
	}
}

var tmplStr string = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN"                            
"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">                                
<html xmlns="http://www.w3.org/1999/xhtml">                                         
  <head>                                                                            
    <title>Stats</title>                                                            
    <meta http-equiv="Content-Type"                                                 
      content="text/html; charset=utf-8"/>                                          
    <link href="style.css" rel="stylesheet" type="text/css"/>                       
  </head>                                                                           
  <body>                                                                            
    <table summary="Test Statistics">                                               
      <caption>Test Statistics</caption>                                            
      <tr>                                                                          
        <th>Event</th>                                                              
        <th>VentureConfigId</th>                                                    
        <th>VentureReference</th>                                                   
        <th>CreatedAt</th>                                                          
        <th>Culture</th>                                                            
        <th>ActionType</th>                                                         
        <th>ActionReference</th>                                                    
        <th>Version</th>                                                            
        <th>Route</th>                                                              
        <th>Payload</th>                                                            
      </tr>                                                                         
      {{range .}}<tr>                                                               
      {{range rangeStruct .}}<td>{{.}}</td>                                         
      {{end}}</tr>                                                                  
      {{end}}                                                                       
    </table>                                                                        
  </body>                                                                           
</html> `

// TestApocalypseApi_ReportWeb checks if the web endpoint
// returns a success http status
func TestApocalypseApi_ReportWeb(t *testing.T) {
	robo := &Apocalypse{}
	os.Remove("./test.db")
	robo.DB = survivordb.Open("./test.db")
	if robo.DB == nil {
		return
	}
	err := robo.DB.Setup()
	if err != nil {
		t.Errorf("Error setting up database: %v", err)
		return
	}
	templ := template.New("index.tmpl").Funcs(TemplateFuncs)
	robo.HTMLTemplateName = "index.tmpl"
	robo.HTMLTemplate, err = templ.Parse(tmplStr)
	if err != nil {
		t.Error(err, "Error parsing the web template")
		return
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/report", nil)
	robo.Report(w, r)
	resp := w.Result()
	if resp.Status != fmt.Sprintf("%d OK", http.StatusOK) {
		t.Errorf("Apocalypse.ReportWeb(w http.ResponseWriter, r *http.Request): want: %v, got: %v", http.StatusOK, resp.Status)
	}
}
