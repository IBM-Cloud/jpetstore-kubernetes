package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"bitbucket.org/ckvist/twilio/twiml"
	_ "github.com/go-sql-driver/mysql"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

// MMSSecret holds all Twilio and other credentials
type MMSSecret struct {
	JPetstoreURL string `json:"jpetstoreurl"`
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

type ClassifyResult struct {
	Filename string        `json:"filename"`
	Labels   []LabelResult `json:"labels"`
}

type LabelResult struct {
	Label       string  `json:"label"`
	Probability float32 `json:"probability"`
}

var (
	graphModel   *tf.Graph
	sessionModel *tf.Session
	labels       []string
)

var mmsSecretsFileName = "mms-secrets.json"

// var twilioSecrets TwilioSecret
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

	// Get mmssearch secrets
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

	// fmt.Println("Twilio SID: ", mmsSecrets.TwilioSecret.AccountSID)
	// fmt.Println("Twilio Token: ", mmsSecrets.TwilioSecret.AccountToken)
	// fmt.Println("Twilio Number: ", mmsSecrets.TwilioSecret.Number)
	fmt.Println("JPetStore URL: ", mmsSecrets.JPetstoreURL)

	if err := loadModel(); err != nil {
		log.Fatal(err)
		return
	}
	http.HandleFunc("/sms/receive", receiveSMSHandler)
	http.HandleFunc("/simulator/receive", receiveSimulatorHandler)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	fmt.Println("http://localhost:8080/")
	http.ListenAndServe(":8080", nil)

}

func loadModel() error {
	fmt.Println("Loading model...")
	// Load inception model
	model, err := ioutil.ReadFile("/model/tensorflow_inception_graph.pb")
	if err != nil {
		return err
	}
	graphModel = tf.NewGraph()
	if err := graphModel.Import(model, ""); err != nil {
		return err
	}

	sessionModel, err = tf.NewSession(graphModel, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Load labels
	labelsFile, err := os.Open("/model/imagenet_comp_graph_label_strings.txt")
	if err != nil {
		return err
	}
	defer labelsFile.Close()
	scanner := bufio.NewScanner(labelsFile)
	// Labels are separated by newlines
	for scanner.Scan() {
		labels = append(labels, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func contains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}

// Handles responses when using the web chat interface
func receiveSimulatorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New message received from simulator")
	imageFile, header, err := r.FormFile("picture")
	// Will contain filename and extension
	imageName := strings.Split(header.Filename, ".")
	if err != nil {
		responseError(w, "Could not read image", http.StatusBadRequest)
		return
	}
	defer imageFile.Close()
	var imageBuffer bytes.Buffer
	// Copy image data to a buffer
	io.Copy(&imageBuffer, imageFile)

	//parse the data POSTed
	/*r.ParseMultipartForm(32 << 20)
	fmt.Println("I am here")
	file, header, err := r.FormFile("picture") // file has the image
	imageName := strings.Split(header.Filename, ".")
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
	var imageBuffer bytes.Buffer
	io.Copy(&imageBuffer, file)
	err = writer.Close()
	if err != nil {
		fmt.Println("whoops:", err)
	}*/
	// Make tensor
	tensor, err := makeTensorFromImage(&imageBuffer, imageName[:1][0])
	if err != nil {
		responseError(w, "Invalid image", http.StatusBadRequest)
		return
	}

	// Run inference
	output, err := sessionModel.Run(
		map[tf.Output]*tf.Tensor{
			graphModel.Operation("input").Output(0): tensor,
		},
		[]tf.Output{
			graphModel.Operation("output").Output(0),
		},
		nil)
	if err != nil {
		responseError(w, "Could not run inference", http.StatusInternalServerError)
		return
	}
	textResponse, dbMediaURL := parseResponse(output[0].Value().([][]float32)[0])
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
		response, e := http.Get(mediaURL)
		if e != nil {
			log.Fatal(e)
		}
		defer response.Body.Close()

		//open a file for writing
		file, err := os.Create("/tmp/asdf.jpg")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Use io.Copy to just dump the response body to the file. This supports huge files
		_, err = io.Copy(file, response.Body)
		var imageBuffer bytes.Buffer
		io.Copy(&imageBuffer, response.Body)
		if err != nil {
			log.Fatal(err)
		}

		tensor, err := makeTensorFromImage(&imageBuffer, "jpg")
		if err != nil {
			responseError(w, "Invalid image", http.StatusBadRequest)
			return
		}

		// Run inference
		output, err := sessionModel.Run(
			map[tf.Output]*tf.Tensor{
				graphModel.Operation("input").Output(0): tensor,
			},
			[]tf.Output{
				graphModel.Operation("output").Output(0),
			},
			nil)
		if err != nil {
			responseError(w, "Could not run inference", http.StatusInternalServerError)
			return
		}

		twilioResponse, dbMediaURL = parseResponse(output[0].Value().([][]float32)[0])
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

type ByProbability []LabelResult

func (a ByProbability) Len() int           { return len(a) }
func (a ByProbability) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByProbability) Less(i, j int) bool { return a[i].Probability > a[j].Probability }

func parseResponse(probabilities []float32) (string, string) {
	var resultLabels []LabelResult
	for i, p := range probabilities {
		if i >= len(labels) {
			break
		}
		resultLabels = append(resultLabels, LabelResult{Label: labels[i], Probability: p})
	}
	// Sort by probability
	sort.Sort(ByProbability(resultLabels))
	// Return top 5 labels

	var response = ""
	var dbMediaURL = ""
	bestResponse := resultLabels[0].Label
	bestScore := resultLabels[0].Probability

	bestResponse2 := resultLabels[1].Label

	if bestScore > 0.2 {
		// strip response to first word and lowercase
		input := []string{"cat", "dog", "fish", "reptile", "bulldog"}
		sort.Strings(input)

		if strings.LastIndex(bestResponse, " ") != -1 {
			bestResponse = bestResponse[strings.LastIndex(bestResponse, " ")+1:]
		}
		if !contains(input, bestResponse) {
			bestResponse = bestResponse2[strings.LastIndex(bestResponse2, " ")+1:]
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
