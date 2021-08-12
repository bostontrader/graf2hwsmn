package main

import (
	"./core"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// https://support.huaweicloud.com/en-us/api-smn/en-us_topic_0036017316.html

type Config struct {
	AccessKey string
	SecretKey string
	EndPoint  string
	ProjectID string
	TopicURN  string
}

func main() {

	configFileNameDefault := "graf2hwsmn.json"
	var configFileName string
	flag.StringVar(&configFileName, "c", configFileNameDefault,
		fmt.Sprintf("The name of a config file to use.  Default is %s", configFileNameDefault))

	flag.Parse()

	// Now read and parse the config file
	content, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	//http.HandleFunc("/", requestHandler)
	http.HandleFunc("/", requestHandlerWrapper(config))

	fmt.Printf("Starting graf2hwsmn server using configFile = %s\n", configFileName)
	if err := http.ListenAndServe(":9112", nil); err != nil {
		log.Fatal(err)
	}

}

// This is how we feed the configuration into the handler
func requestHandlerWrapper(config Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			// If error, do nothing, just return.
			return
		}
		fmt.Printf("Incoming request...\n")

		switch r.Method {

		case "POST":
			fmt.Printf("Incoming POST request.\n")
			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				// If error, do nothing, just return.
				return
			}

			// Unmarshal json of unknown structure and pick out a field named "message"
			// This is useful for Grafana, but you might need to dissect this for your application.
			var arr map[string]interface{}
			if err := json.Unmarshal(reqBody, &arr); err != nil {
				// If error, do nothing, just return.
				return
			}
			var msg = arr["message"]

			send2Huawei(msg, config)

		default:
			fmt.Fprintf(w, "Only the POST method is supported.")
		}
	}
}

func send2Huawei(message2Send interface{}, config Config) {

	//Set the AK/SK to sign and authenticate the request.
	s := core.Signer{
		Key:    config.AccessKey,
		Secret: config.SecretKey,
	}

	// Convert the message2Send that we have now, to a []byte json that huawei wants.
	type JString struct {
		Message string `json:"message"`
	}

	j := JString{
		Message: fmt.Sprintf("%v", message2Send),
	}
	j1, _ := json.Marshal(j)

	// https://support.huaweicloud.com/en-us/api-smn/PublishMessage.html
	// {URI-scheme}://{Endpoint}/{resource-path}?{query-string}
	resourcePath := fmt.Sprintf("v2/%s/notifications/topics/%s/publish", config.ProjectID, config.TopicURN)
	url := fmt.Sprintf("https://%s/%s", config.EndPoint, resourcePath)
	r, _ := http.NewRequest("POST", url, bytes.NewBuffer(j1))
	fmt.Println(url)

	r.Header.Set("Content-Type", "application/json; charset=UTF-8")
	s.Sign(r)
	fmt.Println(r.Header)
	client := http.DefaultClient
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(body))
}
