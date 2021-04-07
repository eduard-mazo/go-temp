package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"example.com/hello/dbo"
	"github.com/gorilla/mux"
)

func UpdateTemp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sensorID, _ := strconv.Atoi(vars["sensorID"])
	temp, _ := strconv.ParseFloat(vars["temp"], 64)
	fmt.Fprintf(w, "pushed sensor: %d with temp %f\n", sensorID, temp)
	dbo.Update(sensorID, temp)
}

//Server Hi!
func Grettings(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Sensor Data!")
}

func GetTemp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sensorID, _ := strconv.Atoi(vars["sensorID"])

	measures := dbo.Get(sensorID)
	/*sensor := sensor{
		SensorID:   1,
		NumSamples: 250,
		Day:        time.Now().Format("01-02-2006"),
		First:      0,
		Last:       0,
		Samples:    []sample{{0, 0}},
	}*/
	js, err := json.Marshal(measures)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
