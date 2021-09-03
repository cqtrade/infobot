package ftx

import (
	"bytes"
	"encoding/json"
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

func retry(retryCount int, f func(count int) (*http.Response, bool, error)) (*http.Response, error) {
	var err error
	var success bool
	var resp *http.Response
	i := 0
	for i == 0 || (i < retryCount && (err != nil || !success)) {
		i++
		resp, success, err = f(i)
	}
	return resp, err
}

func (client *FtxClient) _getRetry(retryCount int, path string, body []byte) (*http.Response, error) {
	var resp *http.Response
	var err error
	resp, err = retry(
		retryCount,
		func(count int) (*http.Response, bool, error) {

			preparedRequest := client.signRequest("GET", path, body)

			resp, err = client.Client.Do(preparedRequest)

			fmt.Println(fmt.Sprintf("Request %d", count))

			if err != nil {
				return resp, false, err
			}

			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return resp, true, err
			}

			if resp.StatusCode == 429 { // retry, rate limit
				fmt.Println(fmt.Sprintf("%+v", resp.Header))
				time.Sleep(time.Second * 61) // TODO needs value from header
				return resp, false, err
			}

			if resp.StatusCode >= 500 && resp.StatusCode < 600 {
				time.Sleep(time.Second * 5 * time.Duration(int64(count))) // retry, server side error, exponential backoff
				return resp, false, err
			}

			// client error or any other non retry
			fmt.Println(fmt.Sprintf("%+v", resp.Header))

			return resp, true, err
		},
	)

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
