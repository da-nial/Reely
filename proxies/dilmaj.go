package proxies

import (
	"IMDK/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type DilmajProxy struct {
	apiKey   string
	endpoint string
}

type DilmajResponse struct {
	Translations []struct {
		Translation string `json:"translation"`
	} `json:"translations"`
	WordCount      int `json:"word_count"`
	CharacterCount int `json:"character_count"`
}

func (d DilmajProxy) getRequest(text, lang string) (*http.Request, error) {
	translateEndpoint := d.endpoint + "/v3/translate?version=2018-05-01"
	translateBody, err := d.getRequestBody(text, lang)
	fmt.Println("host: ", d.endpoint)
	fmt.Println("endpoint: ", translateEndpoint)
	req, err := http.NewRequest("POST", translateEndpoint, translateBody)
	if err != nil {
		log.Println("DilmajProxy: Error creating request: ", err)
		return nil, err
	}
	req.SetBasicAuth("apikey", d.apiKey)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (d DilmajProxy) getRequestBody(text, lang string) (io.Reader, error) {
	var data = make(map[string]interface{})

	data["text"] = text
	data["model_id"] = "en-" + lang

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(data)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (d DilmajProxy) parseResponse(resp *http.Response) (*DilmajResponse, error) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("DilmajProxy: Error reading response: ", err)
		return nil, err
	}
	bodyString := string(bodyBytes)

	if resp.StatusCode != 200 {
		log.Println("DilmajProxy: Response not OK, status: ", resp.StatusCode)
		log.Println("DilmajProxy: Response body: ", bodyString)
		return nil, err
	}

	dilmajResp := &DilmajResponse{}
	err = json.Unmarshal(bodyBytes, dilmajResp)
	if err != nil {
		log.Println("DilmajProxy, Error unmarshalling response to struct: ", bodyString)
		log.Println("DilmajProxy, Error unmarshalling response to struct: ", err)
		return nil, err
	}

	return dilmajResp, nil
}

func (d DilmajProxy) Translate(text, lang string) (string, error) {
	req, err := d.getRequest(text, lang)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("DilmajProxy: Error receiving response: ", err)
		return "", err
	}
	defer resp.Body.Close()

	dilmajResp, err := d.parseResponse(resp)
	if err != nil {
		log.Println("DilmajProxy: Error parsing response: ", err)
		return "", err
	}

	translationText := dilmajResp.Translations[0].Translation
	log.Println("DilmajProxy: Transcript text: ", translationText)

	return translationText, nil
}

var dilmaj DilmajProxy

func InitDilmaj() {
	c := config.GetConfig()

	dilmaj = DilmajProxy{
		apiKey:   c.GetString("dilmaj.apiKey"),
		endpoint: c.GetString("dilmaj.endpoint"),
	}

	log.Printf("Initialized DilmajProxy")
}

func GetDilmaj() DilmajProxy {
	return dilmaj
}
