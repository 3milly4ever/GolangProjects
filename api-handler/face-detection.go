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

// 1.
// Detect initiates the face detection process for a given image URL.
func (f *faceDetector) Detect(url string) chan bool {

	//channel is created to receive the detection result
	resCh := make(chan bool) //will be used to receive the detection result.
	//the request is used to send and receive the detection result

	//Create a request struct with URL and result channel
	req := fdRequest{
		url:   url,   //url of the image to be processed
		resCh: resCh, //channel for receiving the detection result
	}
	//we download the image and send it to the request channel

	//Launch a goroutine to handle the detection process
	go func(req fdRequest) {

		// download the image and assign it to req.imgData
		req.imgData = f.downloadImage(req.url) //we download the image to the request

		// send the modified request to the reqCh channel for Start()
		f.reqCh <- req //this request is sent to the Start function
	}(req) //this req in the function literal ensures that a snapshot of the req variable is created in this Goroutine.
	//even after it is sent to f.reqCh, a copy is kept here.
	//this happens to prevent any  unintended changes that could happen to the req varialbe after this Goroutine is launched.
	//to summarize, the (req) is a way to maintain data isolation between the Goroutine's internal operations and changes
	//that might happen to the original req variable after the Goroutine has started.

	//then the result is returned
	return resCh
}

// Start starts the face detection process
// processes requests sent to the reqCh channel.
// returns the necessary result in the result channel.
func (f *faceDetector) Start() {
	//listens for requests on the reqCh channel using infinite loop. When a request is received, the face detection process is simulated,
	//then a random result is generated.
	for {
		//2.
		req := <-f.reqCh //receives request from faceDetector / Detect() through the request channel
		fmt.Printf("fd processing request: %s\n", req.url)
		time.Sleep(time.Second * 5)                                                 //Simulating detection process
		result := f.detectFace(req.imgData)                                         //finishes detecting, returns if the byte is even or not.
		fmt.Printf("fd processing done - face detected: %v: %s\n", result, req.url) //proccess is complete
		req.resCh <- result                                                         //sents result to the fdRequest result channel
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
			//this goroutine initiates the detection of the urls
			go func(url string) {
				defer wg.Done()

				//1. call the Detect method on the FaceDetector instance
				result := detector.Detect(url)

				//2. the goroutine stores the result and sends it to the result channel we defined above.
				resultCh <- Result{ImageURL: url, DetectionResult: <-result}
			}(imageURL)
		}

		//3. Waits for all goroutines/detections to finish
		//after all goroutines and detections are finished, the results are added to the result channel
		wg.Wait()
		close(resultCh)

		//Process and respond with the collected results
		//4.each result is collected from the resultCh, and appended to the results slice
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
