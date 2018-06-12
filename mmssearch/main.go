package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"bitbucket.org/ckvist/twilio/twiml"
	_ "github.com/go-sql-driver/mysql"
)

// MMSSecret holds all Watson and Twilio and other credentials
type MMSSecret struct {
	JPetstoreURL string `json:"jpetstoreurl"`
	WatsonSecret struct {
		URL     string `json:"url"`
		Note    string `json:"note"`
		APIKey  string `json:"api_key"`
		APIKey2 string `json:"apikey"`
	} `json:"watson"`
	TwilioSecret struct {
		AccountSID   string `json:"sid"`
		AccountToken string `json:"token"`
		Number       string `json:"number"`
	} `json:"twilio"`
}

type ErrorInfo struct {
	ErrorID     string `json:"error_id"`
	Description string `json:"description"`
}

type WarningInfo struct {
	WarningID   string `json:"warning_id"`
	Description string `json:"description"`
}

type ClassResult struct {
	Class         string  `json:"class"`
	Score         float64 `json:"score"`
	TypeHierarchy string  `json:"type_hierarchy,omitempty"`
}

type ClassifierResult struct {
	Classes      []ClassResult `json:"classes"`
	ClassifierID string        `json:"classifier_id"`
	Name         string        `json:"name"`
}

type ClassifiedImage struct {
	Classifiers []ClassifierResult `json:"classifiers"`
	Image       string             `json:"image,omitempty"`
	SourceURL   string             `json:"source_url"`
	ResolvedURL string             `json:"resolved_url"`
	Error       ErrorInfo          `json:"error,omitempty"`
}

type ClassifiedImages struct {
	CustomClasses   int64             `json:"custom_classes"`
	ImagesProcessed int64             `json:"images_processed"`
	Images          []ClassifiedImage `json:"images"`
	Warnings        []WarningInfo     `json:"warnings,omitempty"`
}

type SimulatorResponse struct {
	Response string `json:"response"`
	IMAGEURL string `json:"imageURL"`
}

var mmsSecretsFileName = "mms-secrets.json"
var watsonVRURLStr string
var watsonVersion = "2016-05-20"

// var twilioSecrets TwilioSecret
// var watsonSecrets WatsonSecret
var mmsSecrets MMSSecret
var db *sql.DB
var (
	productid string
	category  string
	name      string
	descn     string
	listprice float64
	attr1     string
)

func main() {
	var mmsSecretFile []byte

	// Get Watson secrets
	mmsSecretFile, err := ioutil.ReadFile("/etc/secrets/" + mmsSecretsFileName)
	if err != nil {
		// If you can't find it in /etc/secrets, check the current dir
		mmsSecretFile, err = ioutil.ReadFile(mmsSecretsFileName)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("Couldn't find file '" + mmsSecretsFileName + "' in /etc/secrets or current dir.")
			os.Exit(1)
		}
	}
	err = json.Unmarshal(mmsSecretFile, &mmsSecrets)
	if err != nil {
		fmt.Println("Something is wrong with your JSON : ", err)
	}
	fmt.Println("(Updated2) Watson URL: ", mmsSecrets.WatsonSecret.URL)

	// if api_key is empty, check apikey
	if mmsSecrets.WatsonSecret.APIKey == "" {
		mmsSecrets.WatsonSecret.APIKey = mmsSecrets.WatsonSecret.APIKey2
	}
	watsonVRURLStr = mmsSecrets.WatsonSecret.URL + "/v3/classify"

	// fmt.Println("Twilio SID: ", mmsSecrets.TwilioSecret.AccountSID)
	// fmt.Println("Twilio Token: ", mmsSecrets.TwilioSecret.AccountToken)
	// fmt.Println("Twilio Number: ", mmsSecrets.TwilioSecret.Number)
	fmt.Println("JPetStore URL: ", mmsSecrets.JPetstoreURL)

	http.HandleFunc("/sms/receive", receiveSMSHandler)
	http.HandleFunc("/simulator/receive", receiveSimulatorHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	fmt.Println("http://localhost:8080/")
	http.ListenAndServe(":8080", nil)

}

// Handles responses when using the web chat interface
func receiveSimulatorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New message received from simulator")

	var Buf bytes.Buffer
	//parse the data POSTed
	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile("picture") // file has the image
	if err != nil {
		fmt.Fprintf(w, "Error reading uploaded image")
		fmt.Println("whoops:", err)
		return
	}
	defer file.Close()

	// grab the image
	bodyImage := &bytes.Buffer{}             // address of a buffer
	writer := multipart.NewWriter(bodyImage) //writer that will write to bodyImage
	part, err := writer.CreateFormFile("images_file", header.Filename)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	_, err = io.Copy(part, file) //copy contents from file to part

	err = writer.Close()
	if err != nil {
		fmt.Println("whoops:", err)
	}
	// call Watson visual recognition API
	req, _ := http.NewRequest("POST", watsonVRURLStr, bodyImage)
	req.Header.Add("Accept", "application/json")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.SetBasicAuth("apikey", mmsSecrets.WatsonSecret.APIKey)
	msgQ := req.URL.Query()
	msgQ.Add("api_key", mmsSecrets.WatsonSecret.APIKey)
	msgQ.Add("version", watsonVersion)
	req.URL.RawQuery = msgQ.Encode()
	client := &http.Client{}
	watsonResp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error calling Watson:", err)
		fmt.Println("Trying again with http instead (Istio)")
		// TODO: Need a better way of determining if running with ISTIO. env var?
		// If running with Istio, communication with Envoy sidecar is http. Envoy will use HTTPS.
		req.URL.Scheme = "http"
		req.URL.Host = req.URL.Host + ":443"
		watsonResp, err = client.Do(req)
		if err != nil {
			fmt.Println("Error calling Watson (Istio):", err)
		}
	}
	// parse data from Watson
	textResponse, dbMediaURL := parseWatsonResponse(watsonResp)
	Buf.Reset()

	// format response
	sResponse := SimulatorResponse{
		Response: textResponse,
		IMAGEURL: dbMediaURL,
	}

	b, _ := json.Marshal(sResponse)
	fmt.Fprintf(w, string(b))
}

// Handles responses when using the Twilio SMS interface
func receiveSMSHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New SMS message received from Twilio")

	err := r.ParseForm()
	if err != nil {
		fmt.Println("whoops:", err)
	}
	sender := r.FormValue("From")
	// body := r.FormValue("Body")
	numMedia, _ := strconv.Atoi(r.FormValue("NumMedia"))
	mediaURL := r.FormValue("MediaUrl0")
	// mediaType := r.FormValue("MediaContentType0")

	twilioResponse := "No response generated"
	dbMediaURL := ""

	if numMedia == 0 {
		twilioResponse = "No images provided"
	} else {
		req, _ := http.NewRequest("POST", watsonVRURLStr, nil)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetBasicAuth("apikey", mmsSecrets.WatsonSecret.APIKey)
		msgQ := req.URL.Query()
		msgQ.Add("api_key", mmsSecrets.WatsonSecret.APIKey)
		msgQ.Add("version", watsonVersion)
		msgQ.Add("url", mediaURL)
		req.URL.RawQuery = msgQ.Encode()
		client := &http.Client{}
		watsonResp, err := client.Do(req)
		if err != nil {
			fmt.Println("whoops:", err)
		}
		twilioResponse, dbMediaURL = parseWatsonResponse(watsonResp)
	}

	resp := twiml.NewResponse()
	if dbMediaURL != "" {
		resp.Action(twiml.Message{
			Body:  fmt.Sprintf(twilioResponse),
			From:  mmsSecrets.TwilioSecret.Number,
			To:    sender,
			Media: dbMediaURL,
		})
	} else {
		resp.Action(twiml.Message{
			Body: fmt.Sprintf(twilioResponse),
			From: mmsSecrets.TwilioSecret.Number,
			To:   sender,
		})
	}
	resp.Send(w)
}

func parseWatsonResponse(watsonResp *http.Response) (string, string) {
	watsonJSON := &ClassifiedImages{}
	watsonBody, err := ioutil.ReadAll(watsonResp.Body)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	err = json.Unmarshal(watsonBody, &watsonJSON)
	if err != nil {
		fmt.Println("whoops:", err)
	}

	var response = ""
	var dbMediaURL = ""
	bestResponse := watsonJSON.Images[0].Classifiers[0].Classes[0].Class
	bestScore := watsonJSON.Images[0].Classifiers[0].Classes[0].Score

	if bestScore > 0.5 {
		// strip response to first word and lowercase
		if strings.LastIndex(bestResponse, " ") != -1 {
			bestResponse = bestResponse[:strings.LastIndex(bestResponse, " ")-1]
		}
		bestResponse = strings.ToLower(bestResponse)

		dbLoc := "/jpetstore"
		val, ok := os.LookupEnv("DB_LOCATION")
		if ok {
			dbLoc = val
		} else {
			response = bestResponse + " : Unable to connect to DB"
			return response, dbMediaURL
		}
		db, err := sql.Open("mysql", "jpetstore:foobar@"+dbLoc)
		if err != nil {
			fmt.Println(err.Error()) // You should use proper error handling
		}

		// Open doesn't open a connection. Validate DSN data:
		err = db.Ping()
		if err != nil {
			fmt.Println(err.Error()) // You should use proper error handling
		}
		defer db.Close()

		rows, err := db.Query("SELECT * FROM product WHERE descn LIKE '%" + bestResponse + "%' OR name LIKE + '%" + bestResponse + "%'")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		rows.Next()
		err = rows.Scan(&productid, &category, &name, &descn)
		if err != nil {
			response = "Our catalog does not have a matching item for " + bestResponse
		} else {
			longName := name + ": " + descn[strings.LastIndex(descn, ">")+1:]
			dbMediaURL = mmsSecrets.JPetstoreURL + descn[strings.LastIndex(descn, "../")+2:strings.LastIndex(descn, ">")-1]
			if !strings.HasPrefix(dbMediaURL, "http") {
				dbMediaURL = "http://" + dbMediaURL
			}
			var ()
			rows, err = db.Query("SELECT listprice, attr1 FROM item WHERE productid LIKE '" + productid + "'")
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			rows.Next()
			err = rows.Scan(&listprice, &attr1)
			if err != nil {
				log.Fatal(err)
			}
			longNamepostfix := " (" + attr1 + ") for $" + strconv.FormatFloat(listprice, 'f', 2, 64)
			response = "Our catalog has a " + longName + longNamepostfix
		}
	} else {
		response = "Photo cannot be adequately be classified. Please try another photo."
	}
	fmt.Println("response: ", response)
	fmt.Println("dbMediaURL: ", dbMediaURL)
	return response, dbMediaURL
}
