package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

//this would be a separate file

// FaceDetector is an interface that defines the behavior of a face detection system.
type FaceDetector interface {
	Start()
	Detect(url string) chan bool
}

// faceDetector implements the FaceDetector interface.
type faceDetector struct {
	//this channel is used to send face detection requests
	reqCh chan fdRequest
}

// fdRequest represents a face detection request which is in the faceDetector struct.
type fdRequest struct {
	url     string
	imgData []byte
	resCh   chan bool
}

// Result represents the detection result for an image
type Result struct {
	ImageURL        string `json:"imageURL"`
	DetectionResult bool   `json:"detectionResult"`
}

const bufSize = 1024

// NewFaceDetector creates a new instance of the faceDetector type.
func NewFaceDetector() FaceDetector {
	return &faceDetector{
		reqCh: make(chan fdRequest, bufSize),
	}
}

// Detect initiates the face detection process for a given image URL.
func (f *faceDetector) Detect(url string) chan bool {
	resCh := make(chan bool) //will be used to receive the detection result.
	//the request is used to send and receive the detection result
	req := fdRequest{
		url:   url,
		resCh: resCh,
	}
	//we download the image and send it to the request channel
	go func(req fdRequest) {
		req.imgData = f.downloadImage(req.url) //we download the image to the request
		f.reqCh <- req                         //this request is sent to the Start function
	}(req)
	//then the result is returned from
	return resCh
}

// Start starts the face detection process
// processes requests sent to the reqCh channel.
// returns the necessary result in the result channel.
func (f *faceDetector) Start() {
	//listens for requests on the reqCh channel using infinite loop. When a request is received, the face detection process is simulated,
	//then a random result is generated.
	for {
		req := <-f.reqCh //receives request from faceDetector / Start()
		fmt.Printf("fd processing request: %s\n", req.url)
		time.Sleep(time.Second * 5)         // Simulating detection process
		result := f.detectFace(req.imgData) //finishes detecting, returns if the byte is even or not.
		fmt.Printf("fd processing done - face detected: %v: %s\n", result, req.url)
		req.resCh <- result //sents result to the fdRequest result channel
	}
}

// downloadImage simulates downloading image data.
func (f *faceDetector) downloadImage(url string) []byte {
	// Simulate downloading image data
	imgData := []byte{1, 2, 3, 4} // Placeholder image data
	return imgData
}

// detectFace simulates the face detection process.
func (f *faceDetector) detectFace(imgData []byte) bool {
	// Simulate face detection by checking if the first byte of the image data is even
	return imgData[0]%2 == 0
}

// this would be the main
func main() {
	// Create a new FaceDetector instance
	detector := NewFaceDetector()

	// Define the HTTP server and request handler
	http.HandleFunc("/detect", func(w http.ResponseWriter, r *http.Request) {
		// Get image URLs from the query parameters
		imageURLs := r.URL.Query()["url"]

		// Create a wait group to wait for all detections to finish
		var wg sync.WaitGroup

		// Channel to receive detection results
		resultCh := make(chan Result, len(imageURLs))

		// Initiate detection for each image URL
		for _, imageURL := range imageURLs {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()

				result := detector.Detect(url)
				//the goroutine stores the result and sends it to the result channel we defined above.
				resultCh <- Result{ImageURL: url, DetectionResult: <-result}
			}(imageURL)
		}

		// Wait for all goroutines/detections to finish
		//after all goroutines and detections are finished, the results are added to the result channel
		wg.Wait()
		close(resultCh)

		// Process and respond with the collected results
		//each result is collected from the resultCh, and appended to the results slice
		var results []Result
		for res := range resultCh {
			results = append(results, res)
		}

		// Send the results in the HTTP response
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%v", results)
	})

	// Start the HTTP server
	http.ListenAndServe(":8080", nil)
}
