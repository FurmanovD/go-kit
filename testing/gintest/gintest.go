package gintest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
)

// GeneratePostJSON generates a gin context for a POST request with
// with given URL and PostData(will be marshaled into JSON)
func GeneratePostJSONContext(
	url string,
	urlKeys map[string]interface{},
	params map[string]interface{},
	postData interface{},
) (*httptest.ResponseRecorder, *gin.Context, error) {

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c = addParams(c, params)

	if c.Request == nil {
		c.Request = httptest.NewRequest("POST", BuildURLWithParams(url, urlKeys), nil)
	}

	if c.Request.Header == nil {
		c.Request.Header = make(http.Header)
	}

	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(postData)
	if err != nil {
		return nil, nil, err
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(jsonbytes))

	return w, c, nil
}

// GenerateGetContext generates a gin context for a GET request
// with given URL and URL parameters
func GenerateGetContext(
	url string,
	urlKeys map[string]interface{},
	params map[string]interface{},
) (*httptest.ResponseRecorder, *gin.Context) {

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c = addParams(addParams(c, urlKeys), params)

	if c.Request == nil {
		c.Request = httptest.NewRequest("GET", BuildURLWithParams(url, urlKeys), nil)
	}

	return w, c
}

func BuildURLWithParams(url string, params map[string]interface{}) string {

	var urlWithParams strings.Builder
	urlWithParams.WriteString(url)
	if len(params) > 0 && !strings.Contains(url, "?") {
		urlWithParams.WriteString("?")
	}

	amp := ""
	for key, val := range params {
		valStr := fmt.Sprintf("%v", val)

		urlWithParams.WriteString(amp + key + "=" + valStr)
		amp = "&"
	}

	return urlWithParams.String()
}

func addParams(c *gin.Context, params map[string]interface{}) *gin.Context {
	for key, val := range params {
		c.Params = append(
			c.Params,
			gin.Param{
				Key:   key,
				Value: fmt.Sprintf("%v", val),
			},
		)
	}

	return c
}
