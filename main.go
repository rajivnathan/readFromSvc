package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const namespace = "rajivtest"

func main() {

	httpClient := newHTTPClient()
	for {
		fmt.Println(doRequest(httpClient))
		time.Sleep(10 * time.Second)
	}
}

func doRequest(httpClient *http.Client) string {

	// url := fmt.Sprintf("http://%s.che-db-cleaner/%s", namespace, userID)
	url := "http://devfile-sample-python-basic-git-rajivtest.apps.rajivdev-07-11-0925.devcluster.openshift.com/rajiv"

	// create request
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, url, nil)
	if err != nil {
		return err.Error()
	}

	// if queryParams != nil {
	// 	req.URL.RawQuery = queryParams.Encode()
	// }

	// do the request
	res, err := httpClient.Do(req)
	if err != nil {
		return err.Error()
	}

	defer closeResponse(res)
	resBody, readError := readBody(res.Body)
	if readError != nil {
		fmt.Println("error while reading body of the response")
		return err.Error()
	}
	return fmt.Sprintf("Response status: '%s' Body: '%s'", res.Status, resBody)
}

// newHTTPClient returns a new HTTP client with some timeout and TLS values configured
func newHTTPClient() *http.Client {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true}, // nolint:gosec
	}
	var httpClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	return httpClient
}

// readBody reads body from a ReadCloser and returns it as a string
func readBody(body io.ReadCloser) (string, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(body)
	return buf.String(), err
}

// closeResponse reads the body and close the response. To be used to prevent file descriptor leaks.
func closeResponse(res *http.Response) {
	if res != nil {
		io.Copy(ioutil.Discard, res.Body) //nolint: errcheck
		defer res.Body.Close()
	}
}
