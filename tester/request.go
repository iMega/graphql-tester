package tester

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

const defaultStyle = "\x1b[0m"
const cyanColor = "\x1b[36m"
const yellowColor = "\x1b[33m"

var (
	APIURL  string
	Dump    bool
	Headers map[string]string

	dumpReq = func(req *http.Request) {
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			MessageError("Error %s", err)
		}
		Message("%s\nREQUEST:\n%s\n%s\n", cyanColor, string(dump), defaultStyle)
	}

	dumpRes = func(res *http.Response) {
		dump, err := httputil.DumpResponse(res, true)
		if err != nil {
			MessageError("Error %s", err)
		}
		Message("%s\nRESPONSE:\n%s\n%s\n", yellowColor, string(dump), defaultStyle)
	}
)

func sendRequest(query []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, APIURL, bytes.NewBuffer(query))
	if err != nil {
		return nil, err
	}
	for k, v := range Headers {
		req.Header.Set(k, v)
	}

	if Dump {
		dumpReq(req)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if Dump {
		dumpRes(res)
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, nil
}
