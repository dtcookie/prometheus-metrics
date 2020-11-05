package rest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Client has no documentation
type Client interface {
	SetHeader(name string, value string) Client                                                      // SetHeader defines a HTTP Request Header to be sent along with every subsequent HTTP request
	Get(url string, callback func(status int, message string, body []byte) error) error              // Get sends a GET request to the specified URL and decodes the response (assumed to be in JSON format) into the given interface{}
	Post(url string, body []byte, calback func(status int, message string, body []byte) error) error // Post sends a POST request to the specified URL using the given []byte as post body and decodes the response (assumed to be in JSON format) into the given interface{}
	Trace(bool)
}

// NewClient creates a new Client
func NewClient() Client {
	return &client{Headers: map[string]string{}}
}

type client struct {
	TraceRequests bool
	Headers       map[string]string
}

func (c *client) Trace(trace bool) {
	c.TraceRequests = trace
}

// SetHeader defines a HTTP Request Header to be sent along with every subsequent HTTP request
func (c *client) SetHeader(name string, value string) Client {
	c.Headers[name] = value
	return c
}

// Get sends a GET request to the specified URL and decodes the response (assumed to be in JSON format) into the given interface{}
func (c *client) Get(url string, callback func(status int, message string, body []byte) error) error {
	trace := c.TraceRequests
	c.TraceRequests = false
	if trace {
		log.Println("GET " + url)
	}

	var response *http.Response
	var err error
	var data []byte

	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	for name, value := range c.Headers {
		request.Header.Set(name, value)
	}

	if response, err = client.Do(request); err != nil {
		return err
	}
	defer response.Body.Close()

	if data, err = ioutil.ReadAll(response.Body); err != nil {
		return err
	}
	if trace {
		log.Println("  " + response.Status)
	}

	return callback(response.StatusCode, response.Status, data)
}

// Post sends a POST request to the specified URL using the given []byte as post body and decodes the response (assumed to be in JSON format) into the given interface{}
func (c *client) Post(url string, body []byte, callback func(status int, message string, body []byte) error) error {
	trace := c.TraceRequests
	c.TraceRequests = false
	if trace {
		log.Println(fmt.Sprintf("POST (%d bytes) %s", len(body), url))
	}

	var response *http.Response
	var err error
	var data []byte

	client := &http.Client{}
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	for name, value := range c.Headers {
		request.Header.Set(name, value)
	}

	if response, err = client.Do(request); err != nil {
		return err
	}
	defer response.Body.Close()

	if data, err = ioutil.ReadAll(response.Body); err != nil {
		return err
	}
	if trace {
		log.Println("  " + response.Status)
	}
	return callback(response.StatusCode, response.Status, data)
}
