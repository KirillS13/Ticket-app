package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
)

func main() {
	const workers = 50
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo2LCJlbWFpbCI6ImZ1cnR1bmFjZXJpbEBnbWFpbC5jb20iLCJleHAiOjE3OTU0NDgxNDEsImlhdCI6MTc5NTM2MTc0MX0.SU5jBSKMAO5CUrcnQSdOZggaQgK4zA87fjC6Si7sTp4"
	payload := []byte(`{
    "event_id": 1,
    "seats" : [{
        "row": 1,
        "column": 2
        }
    ]
}`)
	content := []byte(`{
    "seat": [{
        "row": 1,
        "column": 2
    }]
}`)
	start := make(chan struct{})
	wg := &sync.WaitGroup{}
	client := &http.Client{}
	results := map[string]int{}
	var mux sync.Mutex
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			<-start
			req, err := http.NewRequest("POST", "http://localhost:8080/bookings", bytes.NewBuffer(payload))
			if err != nil {
				fmt.Println(err)
				return
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			res, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer res.Body.Close()
			fmt.Println("response Status:", res.StatusCode)
			if i == 0 {
				body, _ := io.ReadAll(res.Body)
				fmt.Println("Body:", string(body))
			} else {
				_, _ = io.Copy(io.Discard, res.Body)
			}
			mux.Lock()
			results[res.Status]++
			mux.Unlock()
		}()
	}
	req, _ := http.NewRequest("POST", "http://localhost:8080/events/1/hold", bytes.NewBuffer(content))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	holdRes, err := client.Do(req)
	if err != nil {
		fmt.Println("Hold error:", err)
		return
	}
	_, _ = io.Copy(io.Discard, holdRes.Body)
	holdRes.Body.Close()
	fmt.Println("Hold status:", holdRes.Status)

	close(start)
	wg.Wait()
	fmt.Println("Success")
	fmt.Println("Results:", results)
}
