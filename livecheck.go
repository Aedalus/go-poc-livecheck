package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type ResponseStatus string

var RESPONSE_OK ResponseStatus = "ok"
var RESPONSE_FAIL ResponseStatus = "fail"

type Manifest struct {
	HTTPLivechecks []*HTTPLivecheck `yaml:"http_livechecks"`
}

func (m *Manifest) Run(concurrency int) {
	if concurrency < 1 {
		concurrency = 1
	}

	fmt.Printf("Found %d livechecks\n", len(m.HTTPLivechecks))
	fmt.Printf("Running with concurrency: %d\n\n", concurrency)

	sem := make(chan int, concurrency)
	wg := sync.WaitGroup{}
	writeLock := sync.Mutex{}

	for _, livecheck := range m.HTTPLivechecks {
		sem <- 1 // Increment the semaphore
		wg.Add(1)

		go func(l *HTTPLivecheck) {
			defer wg.Done()

			// Add an extra second here to showcase concurrency
			time.Sleep(time.Second * 1)

			result := l.Check()

			writeLock.Lock()
			defer writeLock.Unlock()

			fmt.Printf("%+v\n", l.Name)
			if result.Status == RESPONSE_FAIL {
				fmt.Printf("\tfail: %v\n", result.Error)
			} else if result.Status == RESPONSE_OK {
				fmt.Printf("\tsuccess\n")
			}

			<-sem // Release semaphore
		}(livecheck)
	}

	wg.Wait()
}

type HTTPLivecheck struct {
	Name     string
	Endpoint string
	Expected struct {
		Code int
		Body string
	}
}

func (l *HTTPLivecheck) String() string {
	return fmt.Sprintf("%+v", *l)
}

type LivecheckResponse struct {
	Status ResponseStatus
	Error  error
}

func (l *HTTPLivecheck) Check() LivecheckResponse {
	resp, err := http.Get(l.Endpoint)
	if err != nil {
		return LivecheckResponse{
			Status: RESPONSE_FAIL,
			Error:  err,
		}
	}

	if resp.StatusCode != l.Expected.Code {
		return LivecheckResponse{
			Status: RESPONSE_FAIL,
			Error:  fmt.Errorf("status code failed: got %v, want %v", resp.StatusCode, l.Expected.Code),
		}
	}

	if l.Expected.Body != "" {
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return LivecheckResponse{
				Status: RESPONSE_FAIL,
				Error:  fmt.Errorf("unable to parse body: %v", err),
			}
		}
		body := string(bytes)

		if body != l.Expected.Body {
			return LivecheckResponse{
				Status: RESPONSE_FAIL,
				Error:  fmt.Errorf("body failed: got %v, want %v", body, l.Expected.Body),
			}
		}
	}

	return LivecheckResponse{
		Status: RESPONSE_OK,
		Error:  nil,
	}
}
