package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	proto "github.com/golang/protobuf/proto"
)

// Activity ...
type Activity struct {
	// Strava Activity ID; In-browser at https://www.strava.com/activities/{id]
	ID uint64 `json:"id"`
	// AthleteID - Strava User ID activity owner
	Athlete struct {
		ID int `json:"id"`
	} `json:"athlete"`
	// Strava generated Polylines; base64encoded string
	Map struct {
		// Polyline *only* recorded in response from /v3/activities/{id}
		// PolylineShort recorded in response from /v3/activities/{id}
		// AND /v3/activities/
		Polyline      string `json:"polyline"`
		PolylineShort string `json:"summary_polyline"`
	} `json:"map"`
	// Device Data - /v3/activities/
	StartTime string `json:"start_date_local"`
	// Only recorded in response from /v3/activities/{id}
	Type   string `json:"type"`
	Device string `json:"device_name"`
}

func main() {

	var activity Activity

	f, _ := os.Open("./sample/sample.json")
	defer f.Close()

	// Original
	byteValue, _ := ioutil.ReadAll(f)
	json.Unmarshal(byteValue, &activity)

	// As PB
	activitypb := &ActivityPB{
		Id: activity.ID,
		Athlete: &Athlete{
			Id: int64(activity.Athlete.ID),
		},
		Map: &Map{
			Polyline:        activity.Map.Polyline,
			Summarypolyline: activity.Map.PolylineShort,
		},
		Startdttm: activity.StartTime,
		Type:      activity.Type,
		Device:    activity.Device,
	}

	_, _ = proto.Marshal(activitypb)

}
