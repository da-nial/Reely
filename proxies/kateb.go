package proxies

import (
	"IMDK/config"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

type KatebProxy struct {
	apiKey   string
	endpoint string
}

type KatebResponse struct {
	ResultIndex int `json:"result_index"`
	Results     []struct {
		Final        bool `json:"final"`
		Alternatives []struct {
			Transcript string  `json:"transcript"`
			Confidence float32 `json:"confidence"`
		} `json:"alternatives"`
	} `json:"results"`
}

func (k KatebProxy) getRequest(reviewVoiceFile *multipart.File) (*http.Request, error) {
	convertEndpoint := k.endpoint + "/v1/recognize"
	req, err := http.NewRequest("POST", convertEndpoint, *reviewVoiceFile)
	if err != nil {
		log.Println("KatebProxy: Error creating request: ", err)
		return nil, err
	}
	req.SetBasicAuth("apikey", k.apiKey)
	req.Header.Set("Content-Type", "audio/flac")

	return req, nil
}

func (k KatebProxy) parseResponse(resp *http.Response) (*KatebResponse, error) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("KatebProxy: Error reading response: ", err)
		return nil, err
	}
	bodyString := string(bodyBytes)

	if resp.StatusCode != 200 {
		log.Println("KatebProxy: Response not OK, status: ", resp.StatusCode)
		log.Println("KatebProxy: Response body: ", bodyString)
		return nil, err
	}

	katebResp := &KatebResponse{}
	err = json.Unmarshal(bodyBytes, katebResp)
	if err != nil {
		log.Println("KatebProxy, Error unmarshalling response to struct: ", bodyString)
		log.Println("KatebProxy, Error unmarshalling response to struct: ", err)
		return nil, err
	}

	return katebResp, nil
}

func (k KatebProxy) Transcribe(reviewVoiceFile *multipart.File) (string, error) {
	req, err := k.getRequest(reviewVoiceFile)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("KatebProxy: Error receiving response: ", err)
		return "", err
	}
	defer resp.Body.Close()

	katebResp, err := k.parseResponse(resp)
	if err != nil {

	}

	reviewText := katebResp.Results[0].Alternatives[0].Transcript
	log.Println("KatebProxy, Transcript text: ", reviewText)

	return reviewText, nil
}

var kateb KatebProxy

func InitKateb() {
	c := config.GetConfig()

	kateb = KatebProxy{
		apiKey:   c.GetString("kateb.apiKey"),
		endpoint: c.GetString("kateb.endpoint"),
	}

	log.Printf("Initialized KatebProxy")
}

func GetKateb() KatebProxy {
	return kateb
}
