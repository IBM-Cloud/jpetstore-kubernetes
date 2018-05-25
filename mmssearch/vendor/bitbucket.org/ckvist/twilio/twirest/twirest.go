// Copyright (C) 2014 Cristoffer Kvist. All rights reserved.
// This project is licensed under the terms of the MIT license in LICENSE.

// Package twirest provides a interface to Twilio REST API allowing the user to
// query meta-data from their account and, to initiate calls and send SMS.
package twirest

import (
	"crypto/tls"
	//"crypto/x509"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const ApiVer string = "2010-04-01"

const (
	tag   = 0
	value = 1
)

// TwilioClient struct for holding a http client and user credentials
type TwilioClient struct {
	httpclient            *http.Client
	accountSid, authToken string
}

// Create a new client
func NewClient(accountSid, authToken string) *TwilioClient {
	// certPool := x509.NewCertPool()
	// pemFile, err := os.Open("cacert.pem")
	// if err != nil {
	// 	err = fmt.Errorf("Using host's root CA\n\t%s", err)
	// 	certPool = nil
	// } else {
	// 	defer pemFile.Close()
	// 	bytes, _ := ioutil.ReadAll(pemFile)
	// 	certPool.AppendCertsFromPEM(bytes)
	// }
	tr := &http.Transport{TLSClientConfig: &tls.Config{RootCAs: nil},
		DisableCompression: true}
	client := &http.Client{Transport: tr}

	return &TwilioClient{client, accountSid, authToken}
}

// Request makes a REST resource or action request from twilio servers and
// returns the response. The type of request is determined by the request
// struct supplied.
func (twiClient *TwilioClient) Request(reqStruct interface{}) (
	TwilioResponse, error) {

	twiResp := TwilioResponse{}

	// setup a POST/GET/DELETE http request from request struct
	httpReq, err := httpRequest(reqStruct, twiClient.accountSid)
	if err != nil {
		return twiResp, err
	}
	// add authentication and headers to the http request
	httpReq.SetBasicAuth(twiClient.accountSid, twiClient.authToken)
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("Accept", "*/*")

	response, err := twiClient.httpclient.Do(httpReq)
	if err != nil {
		return twiResp, err
	}

	// Save http status code to response struct
	twiResp.Status.Http = response.StatusCode

	body, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	// parse xml response into twilioResponse struct
	xml.Unmarshal(body, &twiResp)

	twiResp.Status.Twilio, err = exceptionToErr(twiResp)
	return twiResp, err
}

// exceptiontToErr converts a Twilio response exception (if any) to a go error
func exceptionToErr(twir TwilioResponse) (code int, err error) {
	if twir.Exception != nil {
		return twir.Exception.Code, fmt.Errorf("%s",
			twir.Exception.Message)
	}
	return
}

// httpRequest creates a http REST request from the supplied request struct
// and the account Sid
func httpRequest(reqStruct interface{}, accountSid string) (
	httpReq *http.Request, err error) {

	url, err := urlString(reqStruct, accountSid)
	if err != nil {
		return httpReq, err
	}

	queryStr := queryString(reqStruct)

	switch reqStruct.(type) {
	// GET query method
	default:
		if queryStr != "" {
			url = url + "?" + queryStr
		}
		httpReq, err = http.NewRequest("GET", url, nil)
	// DELETE query method
	case DeleteNotification, DeleteOutgoingCallerId,
		DeleteRecording, DeleteParticipant, DeleteQueue:
		if queryStr != "" {
			url = url + "?" + queryStr
		}
		httpReq, err = http.NewRequest("DELETE", url, nil)
	// POST query method
	case SendMessage, MakeCall, ModifyCall, CreateQueue, ChangeQueue,
		DeQueue, UpdateParticipant, UpdateOutgoingCallerId,
		AddOutgoingCallerId:
		requestBody := strings.NewReader(queryStr)
		httpReq, err = http.NewRequest("POST", url, requestBody)

	}

	return httpReq, err
}

// queryString constructs the request string by combining struct tags and
// elements from the request struct. Each element string is being url
// encoded/escaped before included.
func queryString(reqSt interface{}) (qryStr string) {
	switch reqSt := reqSt.(type) {
	default:
	case SendMessage, Messages, MakeCall, Calls, ModifyCall, Accounts,
		Notifications, OutgoingCallerIds, Recordings, UsageRecords,
		CreateQueue, ChangeQueue, DeQueue, Conferences, Participants:
		for i := 0; i < reflect.ValueOf(reqSt).NumField(); i++ {
			fld := reflect.ValueOf(reqSt).Type().Field(i)
			val := reflect.ValueOf(reqSt).Field(i).String()

			if fld.Type.Kind() == reflect.String &&
				string(fld.Tag) != "" && val != "" {
				qryStr += string(fld.Tag) +
					url.QueryEscape(val) + "&"
			}
		}
		// remove the last '&' if we created a query string
		if len(qryStr) > 0 {
			qryStr = qryStr[:len(qryStr)-1]
		}
	}
	return qryStr
}

// urlString constructs the REST resource url
func urlString(reqStruct interface{}, accSid string) (url string, err error) {

	url = "https://api.twilio.com/" + ApiVer + "/Accounts"

	m := make(map[string][2]string)
	// Map the name of the fields in the struct with the values and tags
	switch reqSt := reqStruct.(type) {
	default:
		for i := 0; i < reflect.ValueOf(reqSt).NumField(); i++ {
			fld := reflect.ValueOf(reqSt).Type().Field(i)
			val := reflect.ValueOf(reqSt).Field(i).String()

			m[fld.Name] = [2]string{string(fld.Tag), val}
		}
	}

	// Make base resource URL by adding fields if they exists
	// ... /Accounts/{accSid}/{resource}/{Sid}/{subresource}/{CallSid}
	if fld, ok := m["resource"]; ok {
		url = url + "/" + accSid + fld[tag]
	}
	if fld, ok := m["Sid"]; ok {
		err = required(fld[value])
		url = url + "/" + fld[value]
	}
	if fld, ok := m["subresource"]; ok {
		url = url + fld[tag]
	}
	if fld, ok := m["CallSid"]; ok && fld[tag] == "" {
		url = url + "/" + fld[value]
	}

	// Request cases with additional/optional resources added
	switch reqSt := reqStruct.(type) {
	default:
	case Message:
		if reqSt.Media == true {
			url = url + "/Media"
			if reqSt.MediaSid != "" {
				url = url + "/" + reqSt.MediaSid
			}
		}
	case Call:
		if reqSt.Recordings == true {
			url = url + "/Recordings"
		} else if reqSt.Notifications == true {
			url = url + "/Notifications"
		}
	case UsageRecords:
		url = url + "/" + reqSt.SubResource
	case QueueMember:
		if reqSt.Front && reqSt.CallSid == "" {
			url = url + "/Front"
		}
	case DeQueue:
		if reqSt.Front && reqSt.CallSid == "" {
			url = url + "/Front"
		}
	}

	return url, err
}

// check that string(s) is(are) not empty, return error otherwise
func required(rs ...string) (err error) {
	for _, s := range rs {
		if s == "" {
			return fmt.Errorf("required field missing")
		}
	}
	return
}
