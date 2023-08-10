package proxies

import (
	"IMDK/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type HumanProxy struct {
	apiKey   string
	endpoint string
}

type HumanResponse struct {
	Emotion struct {
		Document struct {
			Emotion struct {
				Anger float32 `json:"anger"`
			} `json:"emotion"`
		} `json:"document"`
	} `json:"emotion"`
}

func (h HumanProxy) getRequest(text string) (*http.Request, error) {
	checkEndpoint := h.endpoint + "/v1/analyze?version=2021-08-01"
	checkBody, err := h.getRequestBody(text)
	if err != nil {
		log.Println("HumanProxy: Error creating request body: ", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", checkEndpoint, checkBody)
	if err != nil {
		log.Println("HumanProxy: Error creating request: ", err)
		return nil, err
	}
	req.SetBasicAuth("apikey", h.apiKey)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (h HumanProxy) getRequestBody(text string) (io.Reader, error) {
	schemaFile, err := os.Open("./proxies/schemas/human.json")
	if err != nil {
		log.Println("HumanProxy: Error reading schema: ", err)
		return nil, err
	}
	defer schemaFile.Close()

	schemaBytes, _ := ioutil.ReadAll(schemaFile)
	var data map[string]interface{}
	err = json.Unmarshal(schemaBytes, &data)
	if err != nil {
		return nil, err
	}

	data["text"] = text

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(data)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (h HumanProxy) parseResponse(resp http.Response) (*HumanResponse, error) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("HumanProxy: Error reading response: ", err)
		return nil, err
	}
	bodyString := string(bodyBytes)

	fmt.Println("HumanProxy resp: ", bodyString)

	if resp.StatusCode != 200 {
		log.Println("HumanProxy: Response not OK, status: ", resp.StatusCode)
		log.Println("HumanProxy: Response body: ", bodyString)
		return nil, err
	}

	var humanResp HumanResponse

	err = json.Unmarshal(bodyBytes, &humanResp)
	if err != nil {
		log.Println("HumanProxy, Error unmarshalling response to struct: ", bodyString)
		log.Println("HumanProxy, Error unmarshalling response to struct: ", err)
		return nil, err
	}

	return &humanResp, nil
}

func (h HumanProxy) Check(text string) (bool, error) {
	req, err := h.getRequest(text)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("HumanProxy: Error receiving response: ", err)
		return false, err
	}
	defer resp.Body.Close()

	humanResponse, err := h.parseResponse(*resp)
	if err != nil {
		return false, err
	}
	anger := humanResponse.Emotion.Document.Emotion.Anger
	if anger > 0.7 {
		return false, nil
	}

	return true, nil
}

var human HumanProxy

func InitHuman() {
	c := config.GetConfig()

	human = HumanProxy{
		apiKey:   c.GetString("human.apiKey"),
		endpoint: c.GetString("human.endpoint"),
	}

	log.Printf("Initialized HumanProxy")
}

func GetHuman() HumanProxy {
	return human
}
