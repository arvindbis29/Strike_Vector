package urlMedia

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"time"
	globalconstant "voice-hack-backend/globalConstant"
	"voice-hack-backend/utilities/httpRequest"
)

type MetaKeys struct {
	ReceiverId string `json:"receiverId"`
	CallerId   string `json:"callerId"`
	ModId      string `json:"modid"`
}
type TranscribeInput struct {
	CallRecordingLink string   `json:"callRecordingLink"`
	CallType          string   `json:"callType"` // PN, C2C or other
	MetaKeys          MetaKeys `json:"metaKeys"`
}
type TranscribeAPIResponse struct {
	Code   int                `json:"Code"`
	Status string             `json:"Status"`
	Data   *TranscribeAPIData `json:"Data"`
}

type TranscribeAPIData struct {
	MediaId          string `json:"MediaId"`
	Status           string `json:"Status"`
	TranscriptionURL string `json:"TranscriptionURL"`
}

func SetInputParamForTranscribeAPI(callRecordingLink string, receiverId string, callerId string) TranscribeInput {
	return TranscribeInput{
		CallRecordingLink: callRecordingLink,
		CallType:          "PNS",
		MetaKeys: MetaKeys{
			ReceiverId: receiverId,
			CallerId:   callerId,
			ModId:      "LMS",
		},
	}
}
func CallTranscribeAPI(input TranscribeInput) (resposne TranscribeAPIData, err string) {

	// ---- Build multipart body ----
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	writer.WriteField("callRecordingLink", input.CallRecordingLink)
	writer.WriteField("callType", input.CallType)

	metaJSON, _ := json.Marshal(input.MetaKeys)
	writer.WriteField("metaKeys", string(metaJSON))

	writer.Close()
	// fmt.Print("Multipart body prepared", writer.FormDataContentType())

	// ---- Prepare request ----
	req := httpRequest.HttpRequest{
		Method:          http.MethodPost,
		URL:             "http://34.47.186.170/transcribe",
		Headers:         map[string]any{},
		MultipartBody:   &buf,
		MultipartWriter: writer,
		Timeout:         40 * time.Second,
	}

	// ---- Make the API call ----
	resp := httpRequest.MakeHttpCall(req)

	// ---- Error handling ----
	if resp.Err != nil {
		return TranscribeAPIData{}, fmt.Sprintf("failed to call Transcribe API: %v", resp.Err)
	}
	if resp.StatusCode != http.StatusOK {
		return TranscribeAPIData{}, fmt.Sprintf("Transcribe API returned non-200 status: %d", resp.StatusCode)
	}

	// ---- Convert map[string]any → typed struct ----
	var parsed TranscribeAPIResponse
	jsonBody, _ := json.Marshal(resp.Body) // convert map → JSON
	if err := json.Unmarshal(jsonBody, &parsed); err != nil {
		return TranscribeAPIData{}, fmt.Sprintf("failed to parse Transcribe API response: %v", err)
	}

	if parsed.Code != http.StatusOK {
		return TranscribeAPIData{}, fmt.Sprintf("Transcribe API returned non-200 status: %d", parsed.Code)
	}

	// ---- Return the parsed response ----
	return *parsed.Data, ""

}

// GetTextFromURL downloads any .txt or plain-text file from a URL
// and returns the content as a string.
func GetTextFromURL(url string) (string, error) {

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// Make GET request
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to call URL: %w", err)
	}
	defer resp.Body.Close()

	// Non-200 response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Convert to string and return
	return string(body), nil
}

func SafeFetchTranscriptAndSummary(inputGLID, inputCustomerType, inputCity string) []DataRow {
	file, err := os.Open("sample/sample_data.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return []DataRow{}
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return []DataRow{}
	}

	var allData []struct {
		GLID         string
		CustomerType string
		CityName     string
		Transcript   string
		Summary      string
	}

	for i, row := range records {
		if i == 0 {
			continue // skip header
		}
		if len(row) < 5 {
			continue
		}
		allData = append(allData, struct {
			GLID         string
			CustomerType string
			CityName     string
			Transcript   string
			Summary      string
		}{
			GLID:         row[0],
			CustomerType: row[1],
			CityName:     row[2],
			Transcript:   row[3],
			Summary:      row[4],
		})
	}

	var result []DataRow
	used := make(map[int]bool)

	// Step 1: Filter by GLID
	for i, row := range allData {
		if row.GLID == inputGLID && len(result) < globalconstant.Max_Length_For_Sample_data {
			result = append(result, DataRow{Transcript: row.Transcript, Summary: row.Summary})
			used[i] = true
		}
		if globalconstant.Max_Length_For_Sample_data <= len(result) {
			return result
		}
	}

	// Step 2: Fill remaining by CustomerType
	for i, row := range allData {
		if !used[i] && row.CustomerType == inputCustomerType && len(result) < globalconstant.Max_Length_For_Sample_data {
			result = append(result, DataRow{Transcript: row.Transcript, Summary: row.Summary})
			used[i] = true
		}
		if globalconstant.Max_Length_For_Sample_data <= len(result) {
			return result
		}
	}

	// Step 3: Fill remaining by CityName
	for i, row := range allData {
		if !used[i] && row.CityName == inputCity && len(result) < globalconstant.Max_Length_For_Sample_data {
			result = append(result, DataRow{Transcript: row.Transcript, Summary: row.Summary})
			used[i] = true
		}
		if globalconstant.Max_Length_For_Sample_data <= len(result) {
			return result
		}
	}

	// Step 4: Fill remaining randomly
	rand.Seed(time.Now().UnixNano())
	for len(result) < globalconstant.Max_Length_For_Sample_data && len(result) < len(allData) {
		idx := rand.Intn(len(allData))
		if !used[idx] {
			result = append(result, DataRow{Transcript: allData[idx].Transcript, Summary: allData[idx].Summary})
			used[idx] = true
		}
		if globalconstant.Max_Length_For_Sample_data <= len(result) {
			return result
		}
	}

	return result
}

// DataRow holds the fields we care about
type DataRow struct {
	Transcript string
	Summary    string
}
