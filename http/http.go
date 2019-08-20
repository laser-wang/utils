package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	clog "github.com/cihub/seelog"
)

const (
	RSP_CODE_OK                    = 200
	RSP_CODE_BAD_REQUEST           = 400
	RSP_CODE_UNAUTHORIZED          = 401
	RSP_CODE_FORBIDDEN             = 403
	RSP_CODE_NOT_FOUND             = 404
	RSP_CODE_INTERNAL_SERVER_ERROR = 500
)

var zeroByte = new([]byte) //pointer to empty []byte

// PostJSON sends an Http Request with using the "POST" method and with
// a "Content-Type" header with application/json and X-Auth-Token" header
// set to the specified token value. The inputValue is encoded to json
// and sent in the body of the request. The response json body is
// decoded into the outputValue. If the response does sends an invalid
// or error status code then an error will be returned. If the Content-Type
// value of the response is not "application/json" an error is returned.
func PostJSON(url string, token string, client http.Client, inputValue interface{}, outputValue interface{}) (err error) {
	body, err := json.Marshal(inputValue)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Auth-Token", token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 201 && resp.StatusCode != 202 {
		err = errors.New("Error: status code != 201 or 202, actual status code '" + resp.Status + "'")
		return
	}

	contentTypeValue := resp.Header.Get("Content-Type")
	if contentTypeValue != "application/json" {
		err = errors.New("Error: Expected a json payload but instead recieved '" + contentTypeValue + "'")
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&outputValue)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

// Delete sends an Http Request with using the "DELETE" method and with
// an "X-Auth-Token" header set to the specified token value. The request
// is made by the specified client.
func Delete(url string, token string, client http.Client) (err error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Auth-Token", token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// Expecting a successful delete
	if !(resp.StatusCode == 200 || resp.StatusCode == 202 || resp.StatusCode == 204) {
		err = fmt.Errorf("Unexpected server response status code on Delete '%s'", resp.StatusCode)
		return
	}

	return nil
}

//GetJSON sends an Http Request with using the "GET" method and with
//an "Accept" header set to "application/json" and the authenication token
//set to the specified token value. The request is made by the
//specified client. The val interface should be a pointer to the
//structure that the json response should be decoded into.
func GetJSON(url string, token string, client http.Client, val interface{}) (err error) {
	req, err := createJSONGetRequest(url, token)
	if err != nil {
		return err
	}

	err = executeRequestCheckStatusDecodeJSONResponse(client, req, val)
	if err != nil {
		return err
	}

	return nil
}

//CallAPI sends an HTTP request using "method" to "url".
//For uploading / sending file, caller needs to set the "content".  Otherwise,
//set it to zero length []byte. If Header fields need to be set, then set it in
// "h".  "h" needs to be even numbered, i.e. pairs of field name and the field
//content.
//
//fileContent, err := ioutil.ReadFile("fileName.ext");
//
//resp, err := CallAPI("PUT", "http://domain/hello/", &fileContent,
//"Name", "world")
//
//is similar to: curl -X PUT -H "Name: world" -T fileName.ext
//http://domain/hello/
func CallAPI(method, url string, content *[]byte, h ...string) (*http.Response, error) {
	if len(h)%2 == 1 { //odd #
		return nil, errors.New("syntax err: # header != # of values")
	}
	//I think the above err check is unnecessary and wastes cpu cycle, since
	//len(h) is not determined at run time. If the coder puts in odd # of args,
	//the integration testing should catch it.
	//But hey, things happen, so I decided to add it anyway, although you can
	//comment it out, if you are confident in your test suites.
	var req *http.Request
	var err error
	req, err = http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(h)-1; i = i + 2 {
		req.Header.Set(h[i], h[i+1])
	}
	req.ContentLength = int64(len(*content))
	if req.ContentLength > 0 {
		req.Body = readCloser{bytes.NewReader(*content)}
		//req.Body = *(new(io.ReadCloser)) //these 3 lines do not work but I am
		//req.Body.Read(content)           //keeping them here in case I wonder why
		//req.Body.Close()                 //I did not implement it this way :)
	}
	client := new(http.Client)
	client.Timeout = time.Second * 10
	return client.Do(req)
}

// CallAPI2 把结果转换成interface再返回
func CallAPI2(method, url string, content *[]byte, h ...string) (interface{}, error) {

	if strings.Contains(url, "manager_private/api/v1/events") == false &&
		strings.Contains(url, "manager_private/api/v1/monitor/keepalive") == false &&
		strings.Contains(url, "nms/api/v1/heart_beat") == false &&
		strings.Contains(url, "transfer_private/api/v1/events") == false {
		clog.Trace("call api url:" + url + " post data:" + string(*content))
	}

	resp, err := CallAPI(method, url, content, h...)

	if err == nil {
		defer resp.Body.Close()

		if resp.StatusCode != RSP_CODE_OK {
			bodyNg, errNg1 := ioutil.ReadAll(resp.Body)
			if errNg1 != nil {
				clog.Error(errNg1.Error())
				return nil, errors.New(resp.Status)
			} else {
				if bodyNg != nil && string(bodyNg) != "" {
					var f interface{}
					errNg2 := json.Unmarshal(bodyNg, &f)
					if errNg2 != nil {
						clog.Error(errNg2.Error())
						clog.Error("bodyNg:" + string(bodyNg))
						return nil, errors.New(resp.Status)
					}
					return f, errors.New(resp.Status)
				} else {
					return nil, errors.New(resp.Status)
				}
			}

			//			return nil, errors.New(resp.Status)
		}

		body, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			clog.Error(err2.Error())
			return nil, err2
		} else {
			if body != nil && string(body) != "" {
				var f interface{}
				err3 := json.Unmarshal(body, &f)
				if err3 != nil {
					clog.Error(err3.Error())
					clog.Error("body:" + string(body))
					return nil, err3
				}

				return f, nil
			} else {
				return nil, nil
			}
		}
	} else {
		clog.Error(err.Error())
		return nil, err
	}
}

// CallAPI2 把结果转换成interface再返回
func CallAPI2Str(method, url string, content *[]byte, h ...string) (string, error) {

	if strings.Contains(url, "manager_private/api/v1/events") == false &&
		strings.Contains(url, "manager_private/api/v1/monitor/keepalive") == false &&
		strings.Contains(url, "nms/api/v1/heart_beat") == false &&
		strings.Contains(url, "transfer_private/api/v1/events") == false {
		clog.Trace("call api url:" + url + " post data:" + string(*content))
	}

	resp, err := CallAPI(method, url, content, h...)

	if err == nil {
		defer resp.Body.Close()

		if resp.StatusCode != RSP_CODE_OK {
			bodyNg, errNg1 := ioutil.ReadAll(resp.Body)
			if errNg1 != nil {
				clog.Error(errNg1.Error())
				return "", errors.New(resp.Status)
			} else {
				if bodyNg != nil && string(bodyNg) != "" {
					return string(bodyNg), errors.New(resp.Status)
				} else {
					return "", errors.New(resp.Status)
				}
			}

			//			return nil, errors.New(resp.Status)
		}

		body, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			clog.Error(err2.Error())
			return "", err2
		} else {
			if body != nil && string(body) != "" {
				return string(body), nil
			} else {
				return "", nil
			}
		}
	} else {
		clog.Error(err.Error())
		return "", err
	}
}

type readCloser struct {
	io.Reader
}

func (readCloser) Close() error {
	//cannot put this func inside CallAPI; golang disallow nested func
	return nil
}

// CheckHTTPResponseStatusCode compares http response header StatusCode against expected
// statuses. Primary function is to ensure StatusCode is in the 20x (return nil).
// Ok: 200. Created: 201. Accepted: 202. No Content: 204. Partial Content: 206.
// Otherwise return error message.
func CheckHTTPResponseStatusCode(resp *http.Response) error {
	switch resp.StatusCode {
	case 200, 201, 202, 204, 206:
		return nil
	case 400:
		return errors.New("Error: response == 400 bad request")
	case 401:
		return errors.New("Error: response == 401 unauthorised")
	case 403:
		return errors.New("Error: response == 403 forbidden")
	case 404:
		return errors.New("Error: response == 404 not found")
	case 405:
		return errors.New("Error: response == 405 method not allowed")
	case 409:
		return errors.New("Error: response == 409 conflict")
	case 413:
		return errors.New("Error: response == 413 over limit")
	case 415:
		return errors.New("Error: response == 415 bad media type")
	case 422:
		return errors.New("Error: response == 422 unprocessable")
	case 429:
		return errors.New("Error: response == 429 too many request")
	case 500:
		return errors.New("Error: response == 500 instance fault / server err")
	case 501:
		return errors.New("Error: response == 501 not implemented")
	case 503:
		return errors.New("Error: response == 503 service unavailable")
	}
	return errors.New("Error: unexpected response status code")
}

func createJSONGetRequest(url string, token string) (req *http.Request, err error) {
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Auth-Token", token)

	return req, nil
}

func executeRequestCheckStatusDecodeJSONResponse(client http.Client, req *http.Request, val interface{}) (err error) {
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	err = CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&val)
	defer resp.Body.Close()

	return err
}
