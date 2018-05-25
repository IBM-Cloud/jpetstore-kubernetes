Twilio
======

Twilio is a go (golang) package for the twilio 
[www.twilio.com](https://www.twilio.com) cloud telephone service API. The package
supports initiating calls, sending SMS and generating TwiML.

Usage
=====

Example Code
============
##TwiML
```go
package main

import (
        "bitbucket.org/ckvist/twilio/twiml"
        "net/http"
)

func helloMonkey(w http.ResponseWriter, r *http.Request) {

        callers := map[string]string{"+15005550001": "Langur"}

        resp := twiml.NewResponse()

        r.ParseForm()
        from := r.Form.Get("From")
        caller, ok := callers[from]

        msg := "Hello monkey"
        if ok {
                msg = "Hello " + caller
        }

        resp.Action(twiml.Say{Text: msg},
                twiml.Play{Url: "http://demo.twilio.com/hellomonkey/monkey.mp3")
        resp.Send(w)
}

func main() {
	http.HandleFunc("/", helloMonkey)
	http.ListenAndServe(":8080", nil)
}
```

##REST
```go
package main

import (
        "bitbucket.org/ckvist/twilio/twirest"
        "fmt"
)

func main() {
        // Test account Sid/Token
        accountSid := "ACdf045ee0ab0e2212ae091a3217660db6"
        authToken := "f74298ebab3a31e099f7161235764b0a"

        client := twirest.NewClient(accountSid, authToken)

        msg := twirest.SendMessage{
                Text: "Hello monkey",
                To:   "+15005550001",
                From: "+15005550005"}

        resp, err := client.Request(msg)
        if err != nil {
                fmt.Println(err)
                return
        }

        fmt.Println(resp.Message.Status)
```

Status
======
Not all functionality is supported nor tested. For example, I have not 
implemented Short Codes and Transcriptions at this time among a few others.

Other features are implemented but they are not fully tested.

License
=======
Released under MIT license. See LICENSE file.
