package insightsGenerateModel

import (
	"encoding/json"
	"fmt"
	"strings"
	globalconstant "voice-hack-backend/globalConstant"
	"voice-hack-backend/utilities/genaiService"
	"voice-hack-backend/utilities/globalFunctions"
	llm "voice-hack-backend/utilities/llmService"
	"voice-hack-backend/utilities/urlMedia"

	"github.com/gin-gonic/gin"
	"google.golang.org/genai"
)

type ApiInputParams struct {
	Glid               int        `json:"glid" binding:"required"`               // GL user ID
	ExecutiveID        string     `json:"executive_id" binding:"required"`       // Executive ID
	CustomerType       string     `json:"customer_type" binding:"required"`      // Customer type: New or Existing
	CustomerCityName   string     `json:"customer_city_name" binding:"required"` // Customer city name
	CallData           []CallData `json:"call_data" binding:"required"`          // List of call details
	TrascriptionURLTxt []string   `json:"transcription_urlTxt"`                  // Transcription URL from call recording
}

type CallData struct {
	CallRecordingURL string `json:"call_recording_url"` // Call recording URL
	CallType         string `json:"call_type"`          // PNS, C2C, or other
	CallDate         string `json:"call_date"`          // Call date
}

type Ensights struct {
	EnsightType string `json:"EnsightType"`
	Concerns    string `json:"Concerns"`
	Resolution  string `json:"Resolution"`
	NextSteps   string `json:"NextSteps"`
	Alert       string `json:"Alert"`
	Sentiment   string `json:"Sentiment"`
	KeyPoints   string `json:"KeyPoints"`
}

type ContentGenerationResponse struct {
	Locations []Ensights `json:"ensights"`
}

type ApiResponse struct {
	Code     int                       `json:"code"`
	Status   string                    `json:"status"`
	Error    string                    `json:"error"`
	Response ContentGenerationResponse `json:"response"`
}

func GetSystemQuery(input ApiInputParams) string {
	var b strings.Builder

	// Header
	b.WriteString("You are an AI-powered Voice Analytics Insight Engine for IndiaMART.\n")
	b.WriteString("Your goal is to extract actionable insights from call transcripts, customer metadata, and historical patterns.\n\n")

	// RULES BASED ON CALL COUNT
	b.WriteString("### RULES FOR INSIGHT GENERATION ###\n")

	if len(input.CallData) > 1 {
		// MULTI-CALL RULES
		b.WriteString("- Multiple calls detected.\n")
		b.WriteString("- Generate one insight block for EACH call.\n")
		b.WriteString("- Use EnsightType = call_1, call_2, etc., based on call index.\n")
		b.WriteString("- After all call-level insights, generate ONE aggregated insight block with EnsightType = final.\n")
		b.WriteString("- The final block must summarize recurring issues, patterns, and business recommendations across all calls.\n")
	} else {
		// SINGLE-CALL RULES
		b.WriteString("- Only ONE call detected.\n")
		b.WriteString("- DO NOT generate any call_1 block.\n")
		b.WriteString("- Generate ONLY ONE insight block using EnsightType = final.\n")
		b.WriteString("- The final block represents the complete insight for this call.\n")
	}

	// COMMON RULES
	b.WriteString("- All insights MUST be actionable with clear next steps.\n")
	b.WriteString("- DO NOT invent any content; only use the provided transcript.\n")
	b.WriteString("- Output MUST be STRICT JSON following the exact structure.\n")
	b.WriteString("- Extract pain points, opportunities, risks, dissatisfaction, and escalation cues.\n")
	b.WriteString("- KeyPoints must be short bullet-style observations.\n\n")

	// FIELD LOGIC
	b.WriteString("### INSIGHT FIELD LOGIC ###\n")
	b.WriteString("EnsightType → call_x or final\n")
	b.WriteString("Concerns → Main customer issue/pain point\n")
	b.WriteString("Resolution → Recommended solution or corrective action\n")
	b.WriteString("NextSteps → Immediate steps to be executed\n")
	b.WriteString("Alert → Any urgency, risk, escalation, or churn signal\n")
	b.WriteString("Sentiment → Positive / Neutral / Negative\n")
	b.WriteString("KeyPoints → Bullet-style summary points\n\n")

	// JSON STRUCTURE EXAMPLE
	example := &ContentGenerationResponse{
		Locations: []Ensights{
			{
				EnsightType: "call_1",
				Concerns:    "Customer confused about subscription renewal",
				Resolution:  "Explain billing structure clearly",
				NextSteps:   "Share plan details and guide user through renewal",
				Alert:       "Possible churn if confusion persists",
				Sentiment:   "Neutral",
				KeyPoints:   "Billing confusion; Need clarity; Renewal assistance",
			},
		},
	}

	exampleBytes, _ := json.MarshalIndent(example, "", "  ")
	b.WriteString("- Output JSON MUST strictly follow this structure:\n```json\n" + string(exampleBytes) + "\n```\n\n")

	// HISTORICAL DATA (if any)
	rows := urlMedia.SafeFetchTranscriptAndSummary(fmt.Sprint(input.Glid), input.CustomerType, input.CustomerCityName)
	if len(rows) > 0 {
		b.WriteString("- Use historical patterns as reference:\n")
		for _, r := range rows {
			b.WriteString("  • Transcript: " + r.Transcript + "\n")
			b.WriteString("  • Summary: " + r.Summary + "\n")
		}
		b.WriteString("\n")
	}

	// BUSINESS PRIORITIES
	b.WriteString("### BUSINESS PRIORITIES ###\n")
	b.WriteString("- Identify recurring pain points and hidden opportunities.\n")
	b.WriteString("- Detect dissatisfaction, churn signals, and escalation risks.\n")
	b.WriteString("- Offer clear recommendations to improve business response.\n")
	b.WriteString("- Keep insights strictly actionable.\n\n")

	b.WriteString("- FINAL OUTPUT MUST be valid JSON following the required fields.\n")

	return b.String()
}


func GenerateUserQuery(apiInputParams ApiInputParams) string {
	var b strings.Builder

	b.WriteString("Analyze the customer's call transcripts and generate structured, actionable insights based on the system prompt.\n\n")

	// CUSTOMER METADATA
	b.WriteString("### Customer Metadata:\n")
	b.WriteString(fmt.Sprintf("- GLID: %d\n", apiInputParams.Glid))
	b.WriteString(fmt.Sprintf("- Executive ID: %s\n", apiInputParams.ExecutiveID))
	b.WriteString(fmt.Sprintf("- Customer Type: %s\n", apiInputParams.CustomerType))
	b.WriteString(fmt.Sprintf("- Customer City: %s\n", apiInputParams.CustomerCityName))
	b.WriteString(fmt.Sprintf("- Total Calls Provided: %d\n\n", len(apiInputParams.CallData)))

	// CALL TRANSCRIPTS
	b.WriteString("### Call Transcripts:\n")
	for i, call := range apiInputParams.CallData {
		b.WriteString(fmt.Sprintf("CALL %d:\n", i+1))
		b.WriteString(fmt.Sprintf("- Call Type: %s\n", call.CallType))
		b.WriteString(fmt.Sprintf("- Call Date: %s\n", call.CallDate))

        // Transcription text
		if len(apiInputParams.TrascriptionURLTxt) > i {
			b.WriteString(fmt.Sprintf("Transcript %d:\n%s\n", i+1, apiInputParams.TrascriptionURLTxt[i]))
		} else {
			b.WriteString(fmt.Sprintf("Transcript %d: [No transcription text available]\n", i+1))
		}
		b.WriteString("\n")
	}

	// // SAMPLE PACKET (optional)
	// if len(apiInputParams.CallData) > 0 {
	// 	sample := FindSampleDatePacket(apiInputParams.CallData[0].CallType)
	// 	if sample != "" {
	// 		b.WriteString("### Sample Data Packet:\n")
	// 		b.WriteString(sample + "\n\n")
	// 	}
	// }

	// INSTRUCTIONS (DYNAMIC – MATCH SYSTEM RULES)
	b.WriteString("### INSTRUCTIONS FOR ANALYSIS ###\n")

	if len(apiInputParams.CallData) > 1 {
		b.WriteString("- Generate actionable insights for EACH CALL using EnsightType = call_1, call_2, etc.\n")
		b.WriteString("- After all call-level insights, generate ONE aggregated insight using EnsightType = final.\n")
	} else {
		b.WriteString("- Only ONE call provided.\n")
		b.WriteString("- DO NOT generate call_1 block.\n")
		b.WriteString("- Generate ONLY ONE insight block using EnsightType = final.\n")
	}

	b.WriteString("- Include Concerns, Resolution, NextSteps, Alert, Sentiment, and KeyPoints.\n")
	b.WriteString("- Do NOT summarize the transcript; only create ACTIONABLE insights.\n")
	b.WriteString("- Use ONLY the transcript and metadata; do NOT invent any content.\n")
	b.WriteString("- Output MUST be valid JSON as per the defined schema.\n")

	return b.String()
}

// func FindSampleDatePacket(callType string) string {
// 	sampleData := map[string]string{
// 		"PN":  "Customer called to inquire about new product features and pricing plans.",
// 		"C2C": "Customer expressed dissatisfaction with recent service outages and requested compensation.",
// 	}
// 	if data, exists := sampleData[callType]; exists {
// 		return data
// 	}
// 	return "General inquiry about services and support."
// }
func GenerateInsightsFromLLM(ginCtx *gin.Context, userQuery string, input ApiInputParams) (result ContentGenerationResponse, err error) {
	// Call LLM service
	systemQuery := GetSystemQuery(input)

	response, err := llm.GenerateInsightsViaLLM(userQuery, systemQuery)
	if err != nil {
		err = fmt.Errorf("LLM service call failed: %v", err)
		return
	}
	fmt.Println("LLM Response:", response)
	// Unmarshal response into result
	rawResponse, _ := globalFunctions.ExtractJson(response)
	if unmarshalErr := json.Unmarshal([]byte(rawResponse), &result); unmarshalErr != nil {
		err = fmt.Errorf("failed to unmarshal LLM response: %v", unmarshalErr)
		return
	}

	return
}
func GenerateInsights(ginCtx *gin.Context, userQuery string, input ApiInputParams) (result ContentGenerationResponse, err error) {
	// Get GenAI client
	client, clientErr := genaiService.GetClient()
	if clientErr != nil {
		err = clientErr
		return
	}

	// Check if model is available
	model := globalconstant.GEMINI_MODEL
	if isModelAvailableErr := genaiService.IsModelAvailable(ginCtx, client, model); isModelAvailableErr != nil {
		err = isModelAvailableErr
		return
	}

	// Prepare conversation contents
	contents := []*genai.Content{
		{
			Role: "model",
			Parts: []*genai.Part{
				{
					Text: GetSystemQuery(input), // System prompt for insights generation
				},
			},
		},
		{
			Role: "user",
			Parts: []*genai.Part{
				{
					Text: userQuery, // User input including transcripts
				},
			},
		},
	}

	// Configure the response schema to match Ensights structure
	contentGenerateConfig := genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"ensights": { // Matches ContentGenerationResponse.Locations
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"EnsightType": {Type: genai.TypeString},
							"Concerns":    {Type: genai.TypeString},
							"Resolution":  {Type: genai.TypeString},
							"Next Steps":  {Type: genai.TypeString},
							"Alert":       {Type: genai.TypeString},
							"Sentiment":   {Type: genai.TypeString},
							"KeyPoints":   {Type: genai.TypeString},
						},
						Required: []string{"EnsightType", "Concerns", "Resolution", "Next Steps", "Alert", "Sentiment", "KeyPoints"},
					},
				},
			},
			Required: []string{"ensights"},
		},
		Tools: nil, //[]*genai.Tool{{GoogleSearch: &genai.GoogleSearch{}}},
		// Optional: can add search or other tools if needed
	}

	// Generate content
	modelResponse, respErr := client.Models.GenerateContent(ginCtx, model, contents, &contentGenerateConfig)
	if respErr != nil {
		err = respErr
		return
	}
	fmt.Println("gemini Response:", modelResponse.Text())
	// Extract and unmarshal JSON response
	rawResponse, _ := globalFunctions.ExtractJson(modelResponse.Text())
	if unmarshalErr := json.Unmarshal([]byte(rawResponse), &result); unmarshalErr != nil {
		err = unmarshalErr
		return
	}

	return
}

func CreateApplicationLogs(ginCtx *gin.Context, apiInputParams ApiInputParams, apiResponse ApiResponse) {
	fileName := "insights_Generate"

	// Build structured log data
	logData := map[string]any{
		// Customer & Executive Info
		"glid":          apiInputParams.Glid,
		"executive_id":  apiInputParams.ExecutiveID,
		"customer_type": apiInputParams.CustomerType,
		"customer_city": apiInputParams.CustomerCityName,
		"total_calls":   len(apiInputParams.CallData),
	}

	// Optional: Include details for each call
	callDetails := make([]map[string]any, len(apiInputParams.CallData))
	for i, call := range apiInputParams.CallData {
		callDetails[i] = map[string]any{
			"call_index":         i + 1,
			"call_type":          call.CallType,
			"call_date":          call.CallDate,
			"transcription_urls": apiInputParams.TrascriptionURLTxt[i],
		}
	}
	logData["call_data"] = callDetails

	// API Response Info
	logData["code"] = apiResponse.Code
	logData["status"] = apiResponse.Status
	logData["response"] = apiResponse.Response // assuming this is already JSON compatible

	// Write JSON log
	globalFunctions.WriteJsonLogs(ginCtx, fileName, logData)
}
