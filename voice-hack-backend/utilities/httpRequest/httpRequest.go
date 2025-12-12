package httpRequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
	"voice-hack-backend/utilities/globalFunctions"
)

// HttpRequest defines the input for the generic HTTP call
type HttpRequest struct {
	Method      string
	URL         string
	Headers     map[string]any
	Body        map[string]any
	QueryParams map[string]any
	Timeout     time.Duration //10 * time.Second
	// NEW — only used when sending multipart
	MultipartBody   *bytes.Buffer     // holds data
	MultipartWriter *multipart.Writer // for Content-Type
}

// HttpResponse defines the output for the generic HTTP call
type HttpResponse struct {
	StatusCode int
	Body       map[string]any
	Header     http.Header
	Err        error
}

// MakeHttpCall executes an HTTP request and returns the response
func MakeHttpCall(reqData HttpRequest) (respData *HttpResponse) {

	respData = &HttpResponse{
		StatusCode: http.StatusInternalServerError,
	}

	// Panic recovery
	defer func() {
		if r := recover(); r != nil {
			respData.StatusCode = http.StatusInternalServerError
			respData.Err = fmt.Errorf("panic in MakeHttpCall: %v", r)
		}
	}()

	// Validate the URL
	finalURL, err := ValidateAndPrepareRequest(&reqData)
	if err != nil {
		respData.Err = fmt.Errorf("ValidateAndPrepareRequest error: %w", err)
		return respData
	}

	// Apply query parameters
	if len(reqData.QueryParams) > 0 {
		q := finalURL.Query()
		for key, value := range reqData.QueryParams {
			q.Set(key, globalFunctions.ConvertValueToString(value))
		}
		finalURL.RawQuery = q.Encode()
	}

	// -------------------------
	// PREPARE REQUEST BODY
	// -------------------------
	var inputReqBody io.Reader

	if reqData.MultipartBody != nil {
		// ---- MULTIPART MODE ----
		inputReqBody = reqData.MultipartBody

	} else if reqData.Body != nil && reqData.Method != http.MethodGet {
		// ---- JSON MODE ----
		jsonBody, err := globalFunctions.ConvertValueToJson(reqData.Body)
		if err != "" {
			respData.Err = fmt.Errorf("error marshaling JSON body: %s", err)
			return respData
		}
		inputReqBody = bytes.NewBuffer(jsonBody)
	}

	// Create HTTP request
	request, err := http.NewRequest(reqData.Method, finalURL.String(), inputReqBody)
	if err != nil {
		respData.Err = fmt.Errorf("error creating HTTP request: %w", err)
		return respData
	}

	// -------------------------
	// HEADERS
	//  -------------------------
	for key, value := range reqData.Headers {
		request.Header.Set(key, globalFunctions.ConvertValueToString(value))
	}

	// Set default Content-Type
	if reqData.MultipartWriter != nil {
		request.Header.Set("Content-Type", reqData.MultipartWriter.FormDataContentType())
	} else if reqData.Body != nil && request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "application/json")
	}

	// Execute HTTP call
	client := &http.Client{
		Timeout: reqData.Timeout,
	}

	resp, err := client.Do(request)
	if err != nil {
		respData.Err = fmt.Errorf("http request failed: %w", err)
		return respData
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		respData.Err = fmt.Errorf("error reading response body: %w", err)
		return respData
	}

	// Unmarshal JSON
	var bodyMap map[string]any
	if len(respBody) > 0 {
		if err := json.Unmarshal(respBody, &bodyMap); err != nil {
			respData.Err = fmt.Errorf("error unmarshaling response JSON: %w", err)
			return respData
		}
	} else {
		bodyMap = make(map[string]any)
	}

	respData.StatusCode = resp.StatusCode
	respData.Body = bodyMap
	respData.Header = resp.Header

	return respData
}

// ValidateAndPrepareRequest validates input and returns parsed URL
func ValidateAndPrepareRequest(reqData *HttpRequest) (*url.URL, error) {
	// Default timeout
	if reqData.Timeout == 0 {
		reqData.Timeout = 10 * time.Second
	}

	// Valid HTTP method
	validMethods := map[string]bool{
		http.MethodGet:     true,
		http.MethodPost:    true,
		http.MethodPut:     true,
		http.MethodDelete:  true,
		http.MethodPatch:   true,
		http.MethodHead:    true,
		http.MethodOptions: true,
	}

	if !validMethods[reqData.Method] {
		return nil, fmt.Errorf("invalid HTTP method: %s", reqData.Method)
	}

	// Validate URL
	parsedURL, err := url.ParseRequestURI(reqData.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Initialize headers & query params if nil
	if reqData.Headers == nil {
		reqData.Headers = make(map[string]any)
	}
	if reqData.QueryParams == nil {
		reqData.QueryParams = make(map[string]any)
	}

	return parsedURL, nil
}

// [
//   {
//     "Input": {
//       "Method": "GET",
//       "URL": "https://jsonplaceholder.typicode.com/posts/1",
//       "Headers": {},
//       "QueryParams": {},
//       "Body": null,
//       "Timeout": 0
//     },
//     "ExpectedOutput": {
//       "StatusCode": 200,
//       "Body": {
//         "userId": 1,
//         "id": 1,
//         "title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
//         "body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum..."
//       },
//       "Header": {
//         "Content-Type": ["application/json; charset=utf-8"],
//         "Content-Length": ["292"]
//       },
//       "Err": null
//     }
//   },
//   {
//     "Input": {
//       "Method": "GET",
//       "URL": "https://jsonplaceholder.typicode.com/comments",
//       "Headers": {},
//       "QueryParams": { "postId": 1 },
//       "Body": null,
//       "Timeout": 0
//     },
//     "ExpectedOutput": {
//       "StatusCode": 200,
//       "Body": [
//         {
//           "postId": 1,
//           "id": 1,
//           "name": "id labore ex et quam laborum",
//           "email": "Eliseo@gardner.biz",
//           "body": "laudantium enim quasi est quidem magnam voluptate ipsam eos..."
//         },
//         {
//           "postId": 1,
//           "id": 2,
//           "name": "quo vero reiciendis velit similique earum",
//           "email": "Jayne_Kuhic@sydney.com",
//           "body": "est natus enim nihil est dolore omnis voluptatem numquam"
//         }
//       ],
//       "Header": { "Content-Type": ["application/json; charset=utf-8"] },
//       "Err": null
//     }
//   },
//   {
//     "Input": {
//       "Method": "POST",
//       "URL": "https://jsonplaceholder.typicode.com/posts",
//       "Headers": { "Authorization": "Bearer mytoken" },
//       "QueryParams": {},
//       "Body": { "title": "foo", "body": "bar", "userId": 1 },
//       "Timeout": 0
//     },
//     "ExpectedOutput": {
//       "StatusCode": 201,
//       "Body": { "title": "foo", "body": "bar", "userId": 1, "id": 101 },
//       "Header": { "Content-Type": ["application/json; charset=utf-8"] },
//       "Err": null
//     }
//   },
//   {
//     "Input": {
//       "Method": "PUT",
//       "URL": "https://jsonplaceholder.typicode.com/posts/1",
//       "Headers": {},
//       "QueryParams": {},
//       "Body": { "id": 1, "title": "updated", "body": "updated body", "userId": 1 },
//       "Timeout": 0
//     },
//     "ExpectedOutput": {
//       "StatusCode": 200,
//       "Body": { "id": 1, "title": "updated", "body": "updated body", "userId": 1 },
//       "Header": { "Content-Type": ["application/json; charset=utf-8"] },
//       "Err": null
//     }
//   },
//   {
//     "Input": {
//       "Method": "DELETE",
//       "URL": "https://jsonplaceholder.typicode.com/posts/1",
//       "Headers": {},
//       "QueryParams": {},
//       "Body": null,
//       "Timeout": 0
//     },
//     "ExpectedOutput": {
//       "StatusCode": 200,
//       "Body": {},
//       "Header": { "Content-Type": ["application/json; charset=utf-8"] },
//       "Err": null
//     }
//   },
//   {
//     "Input": {
//       "Method": "GET",
//       "URL": "invalid-url",
//       "Headers": {},
//       "QueryParams": {},
//       "Body": null,
//       "Timeout": 0
//     },
//     "ExpectedOutput": {
//       "StatusCode": 0,
//       "Body": {},
//       "Header": {},
//       "Err": "invalid URL"
//     }
//   },
//   {
//     "Input": {
//       "Method": "GET",
//       "URL": "https://jsonplaceholder.typicode.com/posts/1",
//       "Headers": {},
//       "QueryParams": {},
//       "Body": null,
//       "Timeout": 1
//     },
//     "ExpectedOutput": {
//       "StatusCode": 200,
//       "Body": { "userId": 1, "id": 1, "title": "...", "body": "..." },
//       "Header": { "Content-Type": ["application/json; charset=utf-8"] },
//       "Err": null
//     }
//   },
//   {
//     "Input": {
//       "Method": "POST",
//       "URL": "https://jsonplaceholder.typicode.com/posts",
//       "Headers": {},
//       "QueryParams": {},
//       "Body": { "title": 123, "body": true, "userId": "abc" },
//       "Timeout": 0
//     },
//     "ExpectedOutput": {
//       "StatusCode": 201,
//       "Body": { "title": "123", "body": "true", "userId": "abc", "id": 101 },
//       "Header": { "Content-Type": ["application/json; charset=utf-8"] },
//       "Err": null
//     }
//   },
//   {
//     "Input": {
//       "Method": "PATCH",
//       "URL": "https://jsonplaceholder.typicode.com/posts/1",
//       "Headers": {},
//       "QueryParams": {},
//       "Body": { "title": "partial update" },
//       "Timeout": 0
//     },
//     "ExpectedOutput": {
//       "StatusCode": 200,
//       "Body": { "id": 1, "title": "partial update", "body": "...", "userId": 1 },
//       "Header": { "Content-Type": ["application/json; charset=utf-8"] },
//       "Err": null
//     }
//   },
//   {
//     "Input": {
//       "Method": "GET",
//       "URL": "https://jsonplaceholder.typicode.com/posts/999999",
//       "Headers": {},
//       "QueryParams": {},
//       "Body": null,
//       "Timeout": 0
//     },
//     "ExpectedOutput": {
//       "StatusCode": 404,
//       "Body": {},
//       "Header": { "Content-Type": ["application/json; charset=utf-8"] },
//       "Err": null
//     }
//   },
//   {
//     "Input": {
//       "Method": "POST",
//       "URL": "https://jsonplaceholder.typicode.com/posts",
//       "Headers": {},
//       "QueryParams": {},
//       "Body": null,
//       "Timeout": 0
//     },
//     "ExpectedOutput": {
//       "StatusCode": 201,
//       "Body": { "id": 101 },
//       "Header": { "Content-Type": ["application/json; charset=utf-8"] },
//       "Err": null
//     }
//   },
//   {
//     "Input": {
//       "Method": "GET",
//       "URL": "https://jsonplaceholder.typicode.com/posts/1",
//       "Headers": { "Accept": "application/xml" },
//       "QueryParams": {},
//       "Body": null,
//       "Timeout": 0
//     },
//     "ExpectedOutput": {
//       "StatusCode": 200,
//       "Body": { "userId": 1, "id": 1, "title": "...", "body": "..." },
//       "Header": { "Content-Type": ["application/json; charset=utf-8"] },
//       "Err": null
//     }
//   }
// ]
