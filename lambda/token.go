package main

import "time"

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

// isNotExpired - Check if token is expired or not...
func (t *Token) isValid() bool {
	expirationTime := time.Unix(t.ExpiresAt, 0).UTC()
	return time.Now().UTC().Before(expirationTime)
}
