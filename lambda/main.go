package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	log "github.com/sirupsen/logrus"
)

// NOTE: For Testing: Download the First `numPages`
// of a user's activities
const numPages = 2

// Activity - Type to Serialize Activity Data from JSON
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

// encodedactivity - Type to pass Strava objects w. an activityID
type encodedactivity struct {
	ActivityID uint64 // Strava ActivityID
	Body       []byte // Serialized object
}

// getActivityHistory - Makes single API call to ../v3/activities/ endpoint
// Requires a valid access token with scope activity:read to download
// user activities
func getActivityHistory(page int, authtoken Token, activityChan chan<- encodedactivity, errorChan chan<- error, wg *sync.WaitGroup) {

	defer wg.Done() // defer decrementing the sync.WG counter...

	// Initialize HTTP client
	client := &http.Client{Timeout: 10 * time.Second}

	// Prepare request and add headers; string format value for
	// page (page number) and per_page (results per page)
	req, _ := http.NewRequest("GET",
		fmt.Sprintf(
			"https://www.strava.com/api/v3/activities?page=%d&per_page=%d", page, 1,
		), nil,
	)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %v", authtoken.AccessToken))

	// Execute request
	resp, err := client.Do(req)

	if err != nil { // On HTTP Error...
		// Log Generic HTTP Error
		log.WithFields(log.Fields{
			"Req": req.URL,
		}).Fatalf(
			"Could not execute HTTP request to https://www.strava.com/api/v3/activities, %e", err)
	}

	defer resp.Body.Close()
	activities := []Activity{}

	// Using [] need to read full resp.Body before statusCode is known
	content, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {

		log.WithFields(log.Fields{
			"statusCode": resp.StatusCode,
		}).Fatalf(
			"Returned NON-200 status code on request to https://www.strava.com/api/v3/activities, %e", err)
	}

	// JSON -> Go Object: Unmarshal API response to []activity{}
	err = json.Unmarshal(content, &activities)
	if err != nil {
		log.Fatalf(
			"Unmarshal Error: Failed to unmarshal response from `./v3/activities` for %d: %e",
			authtoken.Athlete, err,
		)
	}

	// Read from []activity{x, y, z} and send over channel
	for _, activity := range activities {
		b, err := json.Marshal(activity)

		if err != nil {
			log.WithFields(log.Fields{
				"ActivityID": activity.ID,
				"AthleteID":  activity.Athlete.ID,
			}).Errorf("Marshal Error: Failed to marshal response from `./v3/activities`: %e", err)

			err = fmt.Errorf("Marshal Error: Failed to marshal response from `./v3/activities` for %d: %e", activity.ID, err)
			errorChan <- err
		}

		// Send encoded activity Over channel w. ID and content
		activityChan <- encodedactivity{
			ActivityID: activity.ID,
			Body:       b,
		}
	}
	return
}

func explicitRouter(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	switch request.HTTPMethod {
	case "POST":
		return AsyncGetActivities(ctx, request)
	default:
		// Return 405...
		return events.APIGatewayProxyResponse{
			StatusCode: 405,
			Body:       http.StatusText(405),
		}, nil
	}

}

// AsyncGetActivities - func(context.Context, TIn) error
func AsyncGetActivities(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var auth Token
	var wg sync.WaitGroup // Create a synchronizing WorkGroup for API Calls

	_ = json.Unmarshal([]byte(request.Body), &auth)
	log.Info(auth)

	activityChan := make(chan encodedactivity)
	errorChan := make(chan error, 1)

	// NOTE: FIX THIS TO CALL THE TOTAL STATS FOR A USER
	for pgNum := 1; pgNum < numPages; pgNum++ {
		// NOTE: This is Scuffed, but Very rarely does this result in more than 20 Goroutines...
		wg.Add(1)
		go getActivityHistory(pgNum, auth, activityChan, errorChan, &wg)
	}

	go func() {
		wg.Wait()           // this blocks the goroutine until WaitGroup counter is zero
		close(activityChan) // Channels need to be closed, otherwise the below loop will go on forever
	}()

	// Receive From Activities and Error Channel -
	select {
	case err := <-errorChan:
		log.Error("Error Downloading Activity Data:", err)

	case msg := <-activityChan:
		// TODO: DB Writer Workers Go Here...
		fmt.Println(msg.ActivityID)

	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "OK",
	}, nil

}

func main() {
	lambda.Start(explicitRouter)
}
