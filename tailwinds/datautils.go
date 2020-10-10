package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// NOTE: For Testing: Download the First `numPages`
// of a user's activities
const numPages = 2

// S3Region - Region of Bucket to Post Activity files to Storage
var S3Region = os.Getenv("AWS_DEFAULT_REGION")

// S3Bucket - Storage Bucket to Post Activity files
var S3Bucket = os.Getenv("AWS_DEFAULT_S3_BUCKET")

// S3Profile - Profile to Use to Access Storage Bucket
var S3Profile = "tailwinds"

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
			"statusCode": resp.StatusCode,
			"Req":        req.URL,
		}).Error(
			"Could not execute HTTP request",
		)

		// Send Generic HTTP Error on channel
		err = fmt.Errorf(
			"Could not execute HTTP request to https://www.strava.com/api/v3/activities, %v", err,
		)
		errorChan <- err

	}

	defer resp.Body.Close()
	activities := []Activity{}

	// Using [] need to read full resp.Body before statusCode is known
	content, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		log.WithFields(log.Fields{
			"statusCode": resp.StatusCode,
		}).Error(
			"Returned NON-200 status code on request to https://www.strava.com/api/v3/activities",
		)

		// Send Error on channel
		err = fmt.Errorf(
			"Returned NON-200 status code on request to https://www.strava.com/api/v3/activities, %v", err,
		)
		errorChan <- err
	}

	// Unmarshal API response to []activity{}
	err = json.Unmarshal(content, &activities)
	if err != nil {
		log.Error("Unmarshal Error: Failed to unmarshal response from `./v3/activities`: %w", err)
		err = fmt.Errorf(
			"Unmarshal Error: Failed to unmarshal response from `./v3/activities`: %v", err,
		)

		errorChan <- err
	}

	// Read from []activity{x, y, z} and send over channel
	for _, activity := range activities {
		b, err := json.Marshal(activity)

		if err != nil {
			log.WithFields(log.Fields{
				"ActivityID": activity.ID,
				"AthleteID":  activity.Athlete.ID,
			}).Error("Marshal Error: Failed to encode response from `./v3/activities`: %w", err)
		}

		// Send encoded activity Over channel w. ID and content
		activityChan <- encodedactivity{
			ActivityID: activity.ID,
			Body:       b,
		}
	}
	return
}

// makeAWSSession - Create a AWS Session configured to Read + Write to S3
func makeAWSSession() (*session.Session, error) {

	// Create AWS Session
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String(S3Region),
		Credentials: credentials.NewSharedCredentials("", S3Profile),
	})

	if err != nil {
		// failed to create a new AWS Session; check profiles available at
		// ~/.aws/credentials and AWS environment variables
		log.Error(
			"Could Not Initialize AWS Session (Region: %w, Profile: %w) %w",
			S3Region, S3Profile, err,
		)
		return nil, fmt.Errorf(
			"Could Not Initialize AWS Session (Region: %v, Profile: %v) %v",
			S3Region, S3Profile, err,
		)
	}

	// Check that session has valid creds after initialized
	_, err = s.Config.Credentials.Get()

	if err != nil {
		// failed to create a new AWS Session; check profiles available at
		// ~/.aws/credentials and AWS environment variables
		log.Error(
			"Credentials do not exist for %w: %w", S3Profile, err,
		)

		return nil, fmt.Errorf(
			"Credentials do not exist for %v: %v", S3Profile, err,
		)
	}
	return s, nil
}

// AsyncGetActivities xxx
func AsyncGetActivities(auth Token) {

	var wg sync.WaitGroup   // Create a synchronizing WorkGroup for API Calls
	var s3wg sync.WaitGroup // Create a synchronizing WorkGroup for S3

	sess, err := makeAWSSession()
	if err != nil { // NOTE - Middlewares...
		log.Error("Could not connect to S3: %w", err)
		return
	}
	activityChan := make(chan encodedactivity)
	errorChan := make(chan error, 1)

	// NOTE: FIX THIS TO CALL THE TOTAL STATS FOR A USER
	for pgNum := 1; pgNum < numPages; pgNum++ {
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
		log.Error("Error Downloading Activity Data: %w", err)

	case msg := <-activityChan:
		s3wg.Add(1)
		go AddFileToS3(
			sess,
			&s3wg,
			fmt.Sprintf("%v.json", msg.ActivityID),
			msg.Body,
		)
	}

	s3wg.Wait()

}

// AddFileToS3 will upload a single file to S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func AddFileToS3(s *session.Session, wg *sync.WaitGroup, activityID string, b []byte) error {

	defer wg.Done()

	// Config settings: select the bucket, filename, content-type etc.
	// of the file to be uploaded.
	_, err := s3.New(s).PutObject(
		&s3.PutObjectInput{
			Bucket:               aws.String(S3Bucket),
			Key:                  aws.String(fmt.Sprintf("%v.json", activityID)),
			ACL:                  aws.String("private"),
			Body:                 bytes.NewReader(b),
			ContentLength:        aws.Int64(int64(len(b))),
			ContentType:          aws.String(http.DetectContentType(b)),
			ContentDisposition:   aws.String("attachment"),
			ServerSideEncryption: aws.String("AES256"),
		},
	)

	if err != nil {
		log.Error("Could Not Write %w to %w: %w", activityID, S3Bucket, err)
		return fmt.Errorf("Could Not Write %v to %v: %v", activityID, S3Bucket, err)
	}

	return nil
}
