package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	log "github.com/sirupsen/logrus"
)

// Memcache - Object used to abstract memcache.Client
type Memcache struct {
	// Connection Params...
	db *memcache.Client // mc := memcache.New("127.0.0.1:11211", "10.0.0.2:11211", "10.0.0.3:11212")
}

// Token - Contains data from response of https://www.strava.com/oauth/token
// Passed w. all requests to Strava API, note that different activities
// require different scopes
// Example Response:
// {
// "token_type": "Bearer",
// "expires_at": <UNIXTIME>,
// "expires_in": <UINT64>,
// "refresh_token": <base64 string>,
// "access_token": <base64 string>,
// "athlete": {
// "id": 25591100,
// <NON_EXPORTED FIELDS>
// }
// }
type Token struct {
	Athlete struct {
		ID int `json:"id"`
	} `json:"athlete"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

// writeTokenCache - Send a Token to the cache
func (cache *Memcache) writeTokenCache(t Token) (err error) {

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(t)

	// Set w. no expiration to preserve the `refresh-token`
	// for when we need to exchange this for a new token...
	athleteID := strconv.Itoa(t.Athlete.ID)

	err = cache.db.Set(
		&memcache.Item{
			Key:   athleteID,
			Value: b.Bytes(),
		})

	if err != nil {
		// Unknown Cache Set Error
		log.WithFields(log.Fields{
			"AthleteID": athleteID,
		}).Error(
			"Could not set value in Cache",
		)
		return fmt.Errorf(" Could not set value in Cache for %v : %v", athleteID, err)
	}

	log.WithFields(log.Fields{
		"AthleteID": athleteID,
	}).Info(
		"Set value in Cache",
	)

	return nil

}

// readTokenCache - Read a Token from the cache
func (cache *Memcache) readTokenCache(athleteID string) (Token, error) {

	var token Token
	item, err := cache.db.Get(athleteID)

	// Note: cache.db throws Cache Miss Error if and only
	// if item == nil. I.E. these cases are sufficient to
	// check for errors, as MemCache should never return
	// item != nil w.o a Cache Miss Error
	if err != nil && item == nil {
		// Cache Miss
		log.WithFields(log.Fields{
			"AthleteID": athleteID,
		}).Info("Cache Miss", err)

		return token,
			fmt.Errorf("Could not get cached value for %v : %v", athleteID, err)

	} else if err != nil {
		// Some other Error
		log.WithFields(log.Fields{
			"AthleteID": athleteID,
		}).Info(
			"Cache Error", err,
		)
		return token,
			fmt.Errorf("Could not get cached value for %v : %v", athleteID, err)
	}

	// If recive a key, attempt to unmarshal into token
	err = json.Unmarshal(item.Value, &token)
	if err != nil {
		// JSON Unmarshal Error
		log.WithFields(log.Fields{
			"AthleteID": athleteID,
			"Content":   string(item.Value),
		}).Info(
			"Read Token from Cache but failed to Unmarsal", err,
		)
		return Token{},
			fmt.Errorf("Read Token from Cache but failed to Unmarsal responsse for %v : %v",
				athleteID, err,
			)
	}

	return token, nil
}

// refreshTokenCache - Refresh Token from the cache
// NOTE: *MC OR MC??
func refreshTokenCache(origToken Token, cache Memcache) (token Token, err error) {

	// GetUserAccessToken(); Use refresh_token
	token, err = GetUserAccessToken(
		origToken.RefreshToken,
		"refresh_token",
	)

	if err != nil {
		// Some Error from GetUserAccessToken, most often HTTP
		log.WithFields(log.Fields{
			"AthleteID": origToken.Athlete.ID,
		}).Info(err)

		return token,
			fmt.Errorf("Could not get cached value for %v : %v", origToken.Athlete.ID, err)
	}

	_ = cache.writeTokenCache(token)
	return token, nil
}

// isNotExpired - Check if token is expired or not...
func (t *Token) isValid() bool {
	expirationTime := time.Unix(t.ExpiresAt, 0).UTC()
	return time.Now().UTC().Before(expirationTime)
}

// GetUserAccessToken - Make a POST request to the Strava API
// at https://www.strava.com/api/v3/oauth/token to get a short-lived
// user access_token
// For Initial:
//
// For Refresh:
// curl -X POST https://www.strava.com/api/v3/oauth/token \
//   -d client_id=ReplaceWithClientID \
//   -d client_secret=ReplaceWithClientSecret \
//   -d grant_type=refresh_token \
//   -d refresh_token=ReplaceWithRefreshToken
func GetUserAccessToken(code string, grantType string) (token Token, err error) {

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	values := map[string]string{
		"client_id":     os.Getenv("STRAVA_CLIENT_ID"),
		"client_secret": os.Getenv("STRAVA_CLIENT_SECRET"),
		"grant_type":    grantType, // "authorization_code" or "refresh_token"
	}

	// Set grantType to authorization_code or refresh_token
	if grantType == "authorization_code" {
		values["code"] = code
	} else if grantType == "refresh_token" {
		values["refresh_token"] = code
	} else { // If grantType is other; then return other...
		return token, fmt.Errorf("L")
	}

	// Encode map values into form data
	form := url.Values{}
	for k, v := range values {
		form.Add(k, v)
	}

	// Prepare http.Request() w. Encoded form
	req, err := http.NewRequest(
		"POST", "https://www.strava.com/oauth/token", strings.NewReader(form.Encode()),
	)

	req.Header.Add(
		"Content-Type", "application/x-www-form-urlencoded",
	)

	resp, err := client.Do(req)

	// Check API Token Response
	if resp.StatusCode != 200 {
		log.WithFields(log.Fields{
			"statusCode": resp.StatusCode,
			"requestURL": req.URL,
		}).Error(
			"Bad request to Strava Token API",
		)

		return token, fmt.Errorf(
			"Strava Token API returned Non-200 status code %v", resp.StatusCode,
		)
	}

	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(content, &token)

	if err != nil {
		log.Warn("Unmarshalling Error: Failed to unmarshal Strava API response", err)
		return token, fmt.Errorf(
			"Unmarshalling Error: Failed to unmarshal Strava API response",
		)
	}

	return token, nil
}

// getHTTPParam - Check URL for param
func getHTTPParam(req *http.Request, key string) (string, bool) {

	val, ok := req.URL.Query()[key]

	if !ok {
		log.WithFields(log.Fields{
			"requestURL": req.URL,
			"Param":      key,
		}).Error(
			"Redirect does not contain required parameter",
		)
		return "", false
	} else if len(val) == 0 {
		log.WithFields(log.Fields{
			"requestURL": req.URL,
			"Param":      key,
		}).Error(
			"Redirect does not contain required parameter",
		)
		return "", false
	}

	return val[0], true
}
