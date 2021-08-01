package ftx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const URL = "https://ftx.com/api/"

func (client *FtxClient) signRequest(method string, path string, body []byte) *http.Request {
	ts := strconv.FormatInt(time.Now().UTC().Unix()*1000, 10)
	signaturePayload := ts + method + "/api/" + path + string(body)
	signature := client.sign(signaturePayload)
	req, _ := http.NewRequest(method, URL+path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("FTX-KEY", client.Api)
	req.Header.Set("FTX-SIGN", signature)
	req.Header.Set("FTX-TS", ts)
	if client.Subaccount != "" {
		req.Header.Set("FTX-SUBACCOUNT", client.Subaccount)
	}
	return req
}

// TODO Add retry logic HERE
func retry(retryCount int, f func() error) {
	var err error
	i := 0
	for i == 0 || (err != nil && i < retryCount) {
		i++
		err = f()
		time.Sleep(time.Second)
	}
}

func (client *FtxClient) _getRetry(path string, body []byte) (*http.Response, error) {
	var resp *http.Response
	var err error
	retry(3, func() error {
		var err error
		preparedRequest := client.signRequest("GET", path, body)
		resp, err = client.Client.Do(preparedRequest)
		if resp.StatusCode == 429 {
			time.Sleep(time.Second * 61)
			return errors.New("429")
		}
		return err
	})

	return resp, err
}

func (client *FtxClient) _get(path string, body []byte) (*http.Response, error) {
	preparedRequest := client.signRequest("GET", path, body)
	resp, err := client.Client.Do(preparedRequest)
	return resp, err
}

func (client *FtxClient) _post(path string, body []byte) (*http.Response, error) {
	preparedRequest := client.signRequest("POST", path, body)
	resp, err := client.Client.Do(preparedRequest)
	if err != nil {
		fmt.Println("Error _post", err)
	}

	return resp, err
}

func (client *FtxClient) _delete(path string, body []byte) (*http.Response, error) {
	preparedRequest := client.signRequest("DELETE", path, body)
	resp, err := client.Client.Do(preparedRequest)
	return resp, err
}

func _processResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error processing response:", err)
		return err
	}
	err = json.Unmarshal(body, result)
	if err != nil {
		log.Printf("Error processing response:", err)
		return err
	}
	return nil
}
