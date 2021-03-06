package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gorilla/securecookie"
	log "github.com/sirupsen/logrus"
)

var tpl *template.Template
var mc = Memcache{
	db: memcache.New("127.0.0.1:11211"),
}

var hashKey = []byte("32b7ed1b72e6ac06c47dc57d9aebf1d6") // Should be env var...
var s = securecookie.New(hashKey, nil)

const targetScope = "activity:read"

func init() {
	tpl = template.Must(template.ParseGlob("./templates/*"))
}

// index - Handler for Index Page...
func index(w http.ResponseWriter, req *http.Request) {
	// sliced cookie values to only send over images - WTF?
	tpl.ExecuteTemplate(w, "index.html", nil)
}

// manualAuth - Handler for `/exchange_token` route
//
// User manually grants access to application
// Redirects to `/exchange_token` w. code embeded in URL
// code is echanged for a short-lived user access_token`
func manualAuth(w http.ResponseWriter, req *http.Request) {

	// Parse url and extract `code` param
	scope, hasScope := getHTTPParam(req, "scope")
	code, hasCode := getHTTPParam(req, "code")

	if !hasScope || !hasCode { // Check for case when user clicks `Cancel`
		_, hasCancelled := getHTTPParam(req, "error")

		if hasCancelled {
			// If user cancels - redirect back to index.html
			tpl.ExecuteTemplate(w, "index.html", nil)
			return
		}
		// Should never Exit from here...
		return
	}

	// Check for case when user clicks `Auth`; Should match
	// desired app scope
	if strings.Contains(scope, targetScope) {
		// Request access token from Strava for Authenticated User
		// and then Post to Memcache.
		token, _ := GetUserAccessToken(code, "authorization_code")
		mc.writeTokenCache(token)

		// If you have a valid token, you can store Athlete ID in Cookies
		// for now 10-09-2020, unencrypted
		athleteID := strconv.Itoa(token.Athlete.ID)
		encodedID, _ := s.Encode("AthleteID", athleteID)

		http.SetCookie(
			w, &http.Cookie{
				Name:     "AthleteID",
				Value:    encodedID,
				HttpOnly: true,
			},
		)

	} else {
		// Note; replace this...
		tpl.ExecuteTemplate(w, "bad_exchange.html", nil)
		return
	}

	// Render Template
	http.Redirect(w, req, "/download", http.StatusSeeOther)
}

func sendLambdaRequest(authtoken Token) {

	client := &http.Client{Timeout: 10 * time.Second}

	// Prepare request and add headers; string format value for
	// page (page number) and per_page (results per page)
	fmt.Println(authtoken)

	b, _ := json.Marshal(authtoken)

	req, _ := http.NewRequest(
		"POST",
		os.Getenv("AWS_LAMBDA_PROC_ENDOINT"),
		bytes.NewBuffer(b),
	)

	req.Header.Add("Accept", "application/json")

	// Execute request
	_, err := client.Do(req)

	if err != nil { // On HTTP Error...
		fmt.Println(err)
		// Log Generic HTTP Error
		log.WithFields(
			log.Fields{"Req": req.URL},
		).Warn(
			"Could not execute HTTP request to %v, %e", os.Getenv("AWS_LAMBDA_PROC_ENDOINT"), err,
		)
	}

}

// queueDownloadRequest
func queueDownloadRequest(w http.ResponseWriter, req *http.Request) {

	var value string

	// Get Cookie and Decode
	if c, err := req.Cookie("AthleteID"); err == nil {
		s.Decode("AthleteID", c.Value, &value)
	} else {
		// NOTE: LOG!!
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	// Read from Cache && Download Activities
	t, err := mc.readTokenCache(value)
	if err != nil {
		// Returned No token; need to ask for auth again!!
		// NOTE: LOG!!
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !t.isValid() {
		// Returned an Old token; need to refresh...
		t, err = refreshTokenCache(t, mc)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Send to Lambda...
	tpl.ExecuteTemplate(w, "downloading.html", nil)
	sendLambdaRequest(t)
	return
}

func main() {

	// Configure Logging...
	file, err := os.OpenFile("./logs/dload.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // For read access.
	if err != nil {
		log.Fatal(err)
	}

	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
	log.SetOutput(file)

	// Add routes to serve home and download pages
	http.HandleFunc("/", index)
	http.HandleFunc("/exchange_token", manualAuth)
	http.HandleFunc("/download", queueDownloadRequest)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)

	// Test Lambda

}
