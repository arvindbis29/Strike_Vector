package insightsGenerateModel

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	globalconstant "voice-hack-backend/globalConstant"
	"voice-hack-backend/utilities/genaiService"
	getdatafromvectordb "voice-hack-backend/utilities/getDataFromVectorDB"
	"voice-hack-backend/utilities/globalFunctions"
	llm "voice-hack-backend/utilities/llmService"

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
	MaxCallLimit       int        `json:"max_call_limit"`                        // Maximum call limit
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

func GetCustomRules() string {
	return `
	Context:
IndiaMart definition:
IndiaMART is an online B2B marketplace for buyers and suppliers. It connects sellers with buyers. On the IndiaMart portal, sellers list their products & buyer once sees them online send enquiries & make a call to the sellers (This is real business to sellers).
Who are the sellers?
Sellers on IndiaMart can be MSMEs, Large Enterprises, and individual users.
IndiaMart provides premium services to the sellers on a subscription basis, where the seller can purchase Leads (Buyleads), get more direct enquiries (Paid sellers are shown on top of Listing Pages).
Various types of team handling Sellers?
IndiaMart has its own sales, servicing, onboarding & support team for these sellers to provide the above services.
Sales Team-> They sell packages to premium sellers or smartly nudge existing paid sellers for higher packages in case they want more Business (Buyleads & Enquiry).
Onboarding Team/ Servicing Team/ Support Team-> Once the payment is done, these teams give a brief walkthrough of the product & guide them to upload high-quality photos & information that would attract new buyers. 
Based on the seller’s concerns the team gets calls from sellers, I am listing common types of concerns along with their expected resolution that executives should provide over the call, and determine next steps with the stage of call. 
Along with every concern a ticket is being opened & Our objective is to assure in minimum steps with hassle free experience the tickets should get closed in minimum time by resolving the concern.
I am listing some common type of concerns, understand these concerns & how they are being tackled. Similarly whenever a new concern comes provide the resolution & next steps accordingly.
Concern 1: Irrelevant Buyleads
Sellers/customer is getting leads that are not relevant to them, maybe due to 
Category issue (Wrong Product leads), 
Resolution: 
Executive (IndiaMart Team) have the access to sellers category, he should ask questions related to finding relevant category & add that specific category product in Seller’s catalogue, help seller search & find the leads of that category.
Actionable (Executive): Follow-up with seller tomorrow whether he is getting relevant leads or not.
Actionable (Seller): Consume Leads from his desired & check whether they are now relevant or not. 


Location issue (Wrong preferred location-> can be hyperlocal, district, all India, any type of issue), 
Resolution: Executive (IndiaMart Team) have the access to sellers location, should ask questions related to finding correct location from where he want leads & add then add that location in preferred location of the seller, help seller search & find the leads of that category.
Actionable (Executive): Follow-up with seller tomorrow whether he is getting relevant leads or not.
Actionable (Seller): Consume Leads from his desired & check whether they are now relevant or not. 


Low quantity/ Low order value issue, 
Resolution: Executive (IndiaMart Team) have the access to leads which are being generated in the sellers category, should inform about those leads that are there, asks seller to purchase lead as fast as possible. So, that he can get high Quantity leads.


Less Buyleads/Enquiries
Resolution: Executive (IndiaMart Team) should guide seller to put high Quality content (Product should have high quality multiple images, complete specification, Product PDF & Product demo video, complete description of the product) These details can be identified by looking at Product Quality score of each product present on his catlaogue.

Along with this the executive should take a golden opportunity to upsell or upgrade his existing plan, to get more Buyleads & enquiries. As now he he will be listed above higher on listing pages. Generate more trust from Buyers & can increase his Business.

If the seller has already upgraded his package, still not getting leads, then escalate it to IndiaMart Category Teams for more lead generation in his category. [This will only hint the category Team to work on the SEO of that category but will not immediately solve the seller’s problem]
Concern 2: Domain renew/ Domain Change
Sellers/customer want to change the domain (Not recommended as seller would loose the  entire SEO done on that website) or renew the existing domain.
Resolution: Renew/Purchase Domain: Ask the client to login to his domain service provider(example- godaddy etc) , make the renewal of the domain & then reply on the mail provided thread provided by the executive with the login credentials.
Concern 3: Update his catalogue- Images, specification, description, photos, videos etc.
Sellers/customer want to update his catalogue with relevant details add high quality images, remove low quality or wrong images, add or remove product, add or remove title , description, price, specification.
Resolution: 
Executive should guide with steps on the Android App, Desktop to update the Product with the given details. If the seller is not able do it on their own, he should ask the details over the mail & then connect with seller further to get the details updated by himself.
Basis Image, Specification, Price, Title, Description, PDF & videos Product Score is calculated for every Product that executive can also see. According to this he should guide seller further to keep the score as 100 & guide seller to update the details.
Further executive should take follow-up with seller for completing these details & get the ticket closed.
Concern 4: Invoice/Payment concerns
Sellers/customer invoice is being not reflected in his IndiaMart portal, or he want to navigate to the last purchased invoice, or he want to make changes in the invoice or the seller is not able to see his invoice.
Resolution: 
Executive should guide with steps on Desktop/App to see his invoice & executive should verify the issues with himself by seeing on the IndiaMart portal & then escalate it to payments team for resolution. Accordingly guide the seller with further steps & get it solved.
Concern 5: Issue with the Indiamart Product
Seller/customer is not being able to use or open or finding difficult to use any specific sections of the Desktop/Android App
Resolution: Executive should navigate to the problem that seller is facing. Verify the issue with relevant screenshot of that specific section of the product. Also identify the steps that can be used to reproduce that specific issue. 
Then advise & check with general solutions for Desktop like restarting the browser, clearing the cache, opening in incognito mode, login- logout from account, try opening in other computer. Similarly for App the solutions could be like restarting the App, clearing the cache, login-logout, restarting the App.
If the above issue is resolved or not with generic solutions then also escalate the issue to relevant Product teams. 
Example:
Buylead Page related concern in Desktop-  Desktop Buylead Product Team.
Lead Manager Related concern on Android App- Lead Manager Android App team (Lead Manager Team)
Very important to identify the device, section & then escalate to that specific Device & section Teams.
Concern 6: Account details Updation
Seller/customer wanted to update/delete specific details in his account. The seller can update the details like Name, CEO Name, contact person name, Phone number of his own & phone number of his employee, address , changing password etc.
Resolution: Executive should navigate to the details that the seller want to update. Guide him on the device that he is using & get the details updated. 
If there’s any detail that seller cannot update, then ask him to share that specific detail along with relevant proof over the mail. Escalate it to relevant team & get those details updated if he doesn’t have the access to update.
Concern 7: Setting related changes
Seller/customer want some changes in communication of Enquiry/ Buylead like the notification through -email/sms/whatsapp of the modules like 

followup reminder
Replies to your responses
Answered or Missed PNS* Calls
BuyLead Alerts
WhatsApp Communication
Changes in PNS call settings (Preferred Number allotted by IndiaMart: Number shown to buyer which on called route to seller’s attached employees (5) which simultaneously ring on getting call from Buyer). Here are the seller may want to link/unlink the number or disable the PNS call but want his number attached to IndiaMart.
Buylead Preference Settings which contains Location Preferences of Preferred for seller & not preferred for seller, Category Preferences in similar manner Preferred & not Preferred.
Resolution: Executive should first navigate to the exact concern, basis this he should advise him the settings that the seller can change to resolve the issue.
Example: If there’s any issue that is specific to certain location where seller doesn’t want the buyelead from then the same can be updated in the Not Preferred Location.
Take the above concerns as reference, will keep on adding more examples. First Identify the concerns then provide with Next steps for both Customer & executive.
In Next steps make sure to include follow-up for the executives until the issue is resolved
Alert-> If any thing is High Priority, if customer is dissatisfied or will probably not continue his services because of issue, 
Also alert if the seller is potential for upsell, he can upgrade his plan as he wants to increase more Buylead & enquiries
Similarly, If executive is harse, not professional, used abusive languages during the call, telling wrong information about Indiamart, Saying false about indiamart, fight with the customer. Immediately flag.
Very important: Keep insights strictly actionable.
Sentiments has to Positive, Negative or Neutral.
Key topics should have various keywords like mentioned above classification of concers
Agent 1 (Indiviual call:
Role: Act as an Voice Insight expert who listens to the call happening between IndiaMart Sales/servicing/onboaring & support teams, Summarise & understand the current call, also take context from the previous call of the same seller.
Goal: Understand the entire call, extract insight in the below format
Concerns -
Resolution
Next Steps
Alert (If Any)
Sentiment
Key Topics
Agent 2: (Multiple calls of same seller)
Role: Act as an Insight expert who gets the input of various different types of call insights & concern like above. 
Then summarise the insights & actionables for the Persona of IndiaMart Higher authorities ( Higher level teams like Product Manager, Senior sales managers, Senior operations manager, Vice President Level, CEO Level etc)
These insights should be in actionable, quantitative format to take some actions.
Action can be like in 20% calls there was call clarity issues from IndiaMart side only, Multiple product not working escalation of any specific segment, Multiple Less Buylead enquiry issues where we are not convincing seller for upsell.
Goal: Extract the insights from the multiple calls in below quantitative format:
Concerns (Summarise all the concerns quantitively)
Resolution (Summarise the resolution given to all concerns quantaitvely)
Next Steps (Here define actionables with Persona like executive, manager, VP wherever it is required basis call
Alert (If Any)
Summarise the entire alerting of the previous calls with the concerned person to whom the alert has been raised
Sentiment
Summarise the sentiments quantitatively with Positive, Negative & Neutral
Key Topics
With every keywords mention the number of occurrences.

	`

}

// func GetSystemQuery(input ApiInputParams) string {
// 	var b strings.Builder

// 	// Header
// 	b.WriteString("You are an AI-powered Voice Analytics Insight Engine for IndiaMART.\n")
// 	b.WriteString("Your goal is to extract actionable insights from call transcripts, customer metadata, and historical patterns.\n\n")

// 	// RULES BASED ON CALL COUNT
// 	b.WriteString("### RULES FOR INSIGHT GENERATION ###\n")
// 	b.WriteString(GetCustomRules())

// 	if len(input.CallData) > 1 {
// 		// MULTI-CALL RULES
// 		b.WriteString("- Multiple calls detected.\n")
// 		b.WriteString("- Generate one insight block for EACH call.\n")
// 		b.WriteString("- Use EnsightType = call_1, call_2, etc., based on call index.\n")
// 		b.WriteString("- After all call-level insights, generate ONE aggregated insight block with EnsightType = final.\n")
// 		b.WriteString("- The final block must summarize recurring issues, patterns, and business recommendations across all calls.\n")
// 	} else {
// 		// SINGLE-CALL RULES
// 		b.WriteString("- Only ONE call detected.\n")
// 		b.WriteString("- DO NOT generate any call_1 block.\n")
// 		b.WriteString("- Generate ONLY ONE insight block using EnsightType = final.\n")
// 		b.WriteString("- The final block represents the complete insight for this call.\n")
// 	}

// 	// COMMON RULES
// 	// b.WriteString("- All insights MUST be actionable with clear next steps.\n")
// 	// b.WriteString("- DO NOT invent any content; only use the provided transcript.\n")
// 	b.WriteString("- Output MUST be STRICT JSON following the exact structure.\n")
// 	// b.WriteString("- Extract pain points, opportunities, risks, dissatisfaction, and escalation cues.\n")
// 	// b.WriteString("- KeyPoints must be short bullet-style observations.\n\n")

// 	// FIELD LOGIC
// 	// b.WriteString("### INSIGHT FIELD LOGIC ###\n")
// 	// b.WriteString("EnsightType → call_x or final\n")
// 	// b.WriteString("Concerns → Main customer issue/pain point\n")
// 	// b.WriteString("Resolution → Recommended solution or corrective action\n")
// 	// b.WriteString("NextSteps → Immediate steps to be executed\n")
// 	// b.WriteString("Alert → Any urgency, risk, escalation, or churn signal\n")
// 	// b.WriteString("Sentiment → Positive / Neutral / Negative\n")
// 	// b.WriteString("KeyPoints → Bullet-style summary points\n\n")

// 	// JSON STRUCTURE EXAMPLE
// 	example := &ContentGenerationResponse{
// 		Locations: []Ensights{
// 			{
// 				EnsightType: "call_1",
// 				Concerns:    "Customer confused about subscription renewal",
// 				Resolution:  "Explain billing structure clearly",
// 				NextSteps:   "Share plan details and guide user through renewal",
// 				Alert:       "Possible churn if confusion persists",
// 				Sentiment:   "Neutral",
// 				KeyPoints:   "Billing confusion; Need clarity; Renewal assistance",
// 			},
// 		},
// 	}

// 	exampleBytes, _ := json.MarshalIndent(example, "", "  ")
// 	b.WriteString("- Output JSON MUST strictly follow this structure:\n```json\n" + string(exampleBytes) + "\n```\n\n")

// 	// HISTORICAL DATA (if any)
// 	rows := urlMedia.SafeFetchTranscriptAndSummary(fmt.Sprint(input.Glid), input.CustomerType, input.CustomerCityName)
// 	if len(rows) > 0 {
// 		b.WriteString("- Use historical patterns as reference:\n")
// 		for _, r := range rows {
// 			b.WriteString("  • Transcript: " + r.Transcript + "\n")
// 			b.WriteString("  • Summary: " + r.Summary + "\n")
// 		}
// 		b.WriteString("\n")
// 	}

// 	// BUSINESS PRIORITIES
// 	// b.WriteString("### BUSINESS PRIORITIES ###\n")
// 	// b.WriteString("- Identify recurring pain points and hidden opportunities.\n")
// 	// b.WriteString("- Detect dissatisfaction, churn signals, and escalation risks.\n")
// 	// b.WriteString("- Offer clear recommendations to improve business response.\n")
// 	// b.WriteString("- Keep insights strictly actionable.\n\n")

// 	b.WriteString("- FINAL OUTPUT MUST be valid JSON following the required fields.\n")

// 	return b.String()
// }

func GetSystemQuery(input ApiInputParams) string {
	var b strings.Builder

	// --- 1. ROLE DEFINITION & CONTEXT ---
	b.WriteString("You are an expert Voice Analytics Insight Engine for IndiaMART.\n")
	b.WriteString("Your primary goal is to extract actionable insights from call transcripts, customer metadata, and historical patterns.\n")
	b.WriteString("Your analysis must support two objectives:\n")
	b.WriteString("1. Improve IM Executive performance (Sales, Servicing, Onboarding, Support).\n")
	b.WriteString("2. Flag systemic concerns to relevant management teams (Product, Category, Sales VP).\n\n")

	// TEAM ROLES SUMMARY (Focus on the Call Context)
	b.WriteString("### IM EXECUTIVE ROLES ###\n")
	b.WriteString("- **Sales Team:** Focus on upselling higher-value premium packages (Buyleads, Enquiry).\n")
	b.WriteString("- **Servicing/Onboarding/Support:** Focus on product walkthrough, catalogue quality, and resolving concerns to close tickets quickly and efficiently.\n\n")
	b.WriteString("### KEY PLATFORM CONCEPTS ###\n")
	b.WriteString("- IndiaMART: B2B marketplace connecting Sellers (Customers) with Buyers.\n")
	b.WriteString("- Key Goals: Maximize Buyleads, optimize Catalogue Quality Score, ensure hassle-free experience.\n\n")

	// --- 2. LOGIC BRANCHING (SINGLE VS MULTI CALL) ---
	b.WriteString("### PROCESSING MODE ###\n")
	if len(input.CallData) > 1 {
		// AGENT 2: MULTI-CALL LOGIC
		b.WriteString("MODE: MULTI-CALL ANALYSIS (Agent 2 Logic)\n")
		b.WriteString("- Analyze multiple calls for the same seller.\n")
		b.WriteString("- Generate individual insight blocks for each call (EnsightType = call_1, call_2...).\n")
		b.WriteString("- Generate ONE 'final' aggregated block.\n")
		b.WriteString("- The FINAL block must be QUANTITATIVE:\n")
		b.WriteString("  * Summarize concerns with counts (e.g., 'Irrelevant Leads (2 calls)').\n")
		b.WriteString("  * Identify recurring failures or improvements across the timeline.\n")
		b.WriteString("  * Define Actionables for specific personas (Executive, Manager, VP) based on severity.\n")
	} else {
		// AGENT 1: SINGLE-CALL LOGIC
		b.WriteString("MODE: SINGLE CALL ANALYSIS (Agent 1 Logic)\n")
		b.WriteString("- Analyze a specific individual call.\n")
		b.WriteString("- Output ONLY ONE insight block with EnsightType = final.\n")
	}

	// --- 3. DOMAIN SPECIFIC RULES (CONCERNS & RESOLUTIONS) ---
	b.WriteString("\n### DOMAIN KNOWLEDGE & RESOLUTION MATRIX ###\n")
	b.WriteString("Classify issues into these specific categories and verify if the Executive followed the correct Resolution/Next Steps:\n\n")

	b.WriteString("1. IRRELEVANT BUYLEADS (Category/Location/Value issue)\n")
	b.WriteString("   - Resolution: Executive must check category/location settings, add specific product categories, or suggest 'Filters'.\n")
	b.WriteString("   - Upsell Opportunity: If leads are less, nudge for higher package/TrustSeal/STAR/LEADER.\n")

	b.WriteString("2. DOMAIN RENEW/CHANGE\n")
	b.WriteString("   - Resolution: Guide seller to Godaddy/Provider login & upgrade from there. changing domain is NOT recommended (SEO loss).\n")

	b.WriteString("3. CATALOGUE UPDATE (Images, Specs, Score)\n")
	b.WriteString("   - Resolution: Guide on App/Desktop. Aim for Product Score 100 (Images, PDF, Video, Desc).\n")

	b.WriteString("4. INVOICE/PAYMENT ISSUES\n")
	b.WriteString("   - Resolution: Verify on portal, guide user to invoice section. Escalate to Payments Team if system error.\n")

	b.WriteString("5. PRODUCT/TECH ISSUES (App/Desktop)\n")
	b.WriteString("   - Resolution: Troubleshoot (Clear Cache, Incognito, Update App). If fails, escalate to Product Team (Device specific).\n")

	b.WriteString("6. ACCOUNT UPDATES (Contact, Address, Password)\n")
	b.WriteString("   - Resolution: Update on call if allowed. If not, ask for proof via mail and escalate.\n")

	b.WriteString("7. SETTINGS (PNS, Alerts, Whatsapp)\n")
	b.WriteString("   - Resolution: Configure 'Preferred/Not Preferred' locations or categories. Link/Unlink PNS numbers.\n")

	// --- 4. ALERTING LOGIC ---
	b.WriteString("\n### ALERT CATEGORIES (MANDATORY) ###\n")
	b.WriteString("Trigger an 'Alert' field ONLY for these specific scenarios:\n")
	b.WriteString("1. INTERNAL PROCESS FAILURE: Customer bounced between teams, conflicting info given.\n")
	b.WriteString("2. COMPETITOR/CHURN RISK: Mention of competitors, better external offers, or threat to leave.\n")
	b.WriteString("3. EXECUTIVE INEFFICIENCY: Hold time >120s, no clear next step, rude/unprofessional behavior, giving false info.\n")
	b.WriteString("4. UPSELL OPPORTUNITY: Need more leads.\n")

	// --- 5. FIELD REQUIREMENTS ---
	b.WriteString("\n### OUTPUT REQUIREMENTS ###\n")
	b.WriteString("- Concerns: Short, bullet-style.\n")
	b.WriteString("- Resolution: What was done/advised.\n")
	b.WriteString("- NextSteps: Specific follow-ups (e.g., 'Check lead quality tomorrow', 'Share catalog report').\n")
	b.WriteString("- Sentiment: Positive / Neutral / Negative / Angry -> Neutral.\n")
	b.WriteString("- KeyPoints: Concise keywords (e.g., 'Buy leads, Catalog, Filter').\n")
	b.WriteString("- Output MUST be valid JSON strictly matching the example below.\n")

	// --- 6. JSON STRUCTURE ---
	example := &ContentGenerationResponse{
		Locations: []Ensights{
			{
				EnsightType: "final",
				Concerns:    "Insufficient buy leads; Irrelevant location leads",
				Resolution:  "Advised bulk filters; Checked location settings",
				NextSteps:   "Executive to call back tomorrow to verify lead quality",
				Alert:       "Upsell Opportunity: Seller wants high quantity leads",
				Sentiment:   "Negative -> Neutral",
				KeyPoints:   "Buy leads, Location Filter, Upsell",
			},
		},
	}
	exampleBytes, _ := json.MarshalIndent(example, "", "  ")
	b.WriteString("\nJSON STRUCTURE:\n" + string(exampleBytes) + "\n\n")

	// --- 7. HISTORICAL CONTEXT INJECTION ---
	// If history exists, inject it to help the model detect recurring patterns
	// rows := urlMedia.SafeFetchTranscriptAndSummary(fmt.Sprint(input.Glid), input.CustomerType, input.CustomerCityName)
	// if len(rows) > 0 {
	// 	b.WriteString("### HISTORICAL CONTEXT (Previous Calls) ###\n")
	// 	for _, r := range rows {
	// 		b.WriteString("• Summary: " + r.Summary + "\n")
	// 	}
	// 	b.WriteString("Use this history to detect repeat tickets or unresolved recurring issues.\n\n")
	// }
	b.WriteString("- Use historical patterns as reference:\n")

	n := len(input.TrascriptionURLTxt)
	if n == 0 {
		b.WriteString("  No transcription data available.\n")

	}

	callsPerTranscript := input.MaxCallLimit / n
	if callsPerTranscript == 0 {
		callsPerTranscript = 1 // ensure at least 1 call per transcript
	}

	for i := 0; i < n; i++ {
		sampleCalls, _ := getdatafromvectordb.GetSampleCalls(input.TrascriptionURLTxt[i], callsPerTranscript)
		b.WriteString(sampleCalls)
		b.WriteString("\n") // optional: add spacing between each transcript
	}

	b.WriteString("Ensure the response is strictly valid JSON.")
	fmt.Println(b.String())
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

//	func FindSampleDatePacket(callType string) string {
//		sampleData := map[string]string{
//			"PN":  "Customer called to inquire about new product features and pricing plans.",
//			"C2C": "Customer expressed dissatisfaction with recent service outages and requested compensation.",
//		}
//		if data, exists := sampleData[callType]; exists {
//			return data
//		}
//		return "General inquiry about services and support."
//	}
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
		"glid":          apiInputParams.Glid,
		"executive_id":  apiInputParams.ExecutiveID,
		"customer_type": apiInputParams.CustomerType,
		"customer_city": apiInputParams.CustomerCityName,
		"total_calls":   len(apiInputParams.CallData),
	}

	// Optional: Include details for each call safely
	callDetails := make([]map[string]any, len(apiInputParams.CallData))
	for i, call := range apiInputParams.CallData {
		transcript := "[No transcription available]"
		if len(apiInputParams.TrascriptionURLTxt) > i {
			transcript = apiInputParams.TrascriptionURLTxt[i]
		}

		callDetails[i] = map[string]any{
			"call_index":         i + 1,
			"call_type":          call.CallType,
			"call_date":          call.CallDate,
			"transcription_urls": transcript,
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

// StoreInsightsToJSON appends new insights into insights_data.json
func StoreInsightsToJSON(newData []Ensights) error {
	fileName := "insights_data.json"

	var existingData []Ensights

	// -------------------------------
	// 1. Check if file exists
	// -------------------------------
	_, err := os.Stat(fileName)
	if err == nil {
		// File exists → read old data
		fileBytes, err := os.ReadFile(fileName)
		if err == nil {
			// Unmarshal old data
			json.Unmarshal(fileBytes, &existingData)
		}
	}

	// -------------------------------
	// 2. Append new data to old data
	// -------------------------------
	if len(newData) > 1 {
		newData = newData[:len(newData)-1]
	}
	existingData = append(existingData, newData...)
	// If newData contains more than one element, remove the last index of existingData
	// if len(newData) > 1 && len(existingData) > 0 {
	// 	existingData = existingData[:len(existingData)-1]
	// }

	// -------------------------------p
	// 3. Save the updated array
	// -------------------------------
	jsonData, err := json.MarshalIndent(existingData, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON file: %w", err)
	}

	return nil
}

func FinalSummary(apiInputParams ApiInputParams) (*Ensights, error) {
	// 1️⃣ Load stored insights from JSON
	allInsights, err := LoadAllInsightsFromJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to load insights from JSON: %v", err)
	}

	if len(allInsights) == 0 {
		return nil, errors.New("no insights found in JSON")
	}

	// 2️⃣ Build system/user query for final insight
	systemQuery := BuildFinalInsightsQuery(allInsights, apiInputParams)
	fmt.Println("Final Summary System Query:", systemQuery)

	// 3️⃣ Call the LLM service (Gemini)
	respText, err := llm.GenerateInsightsViaLLM("", systemQuery)
	if err != nil {
		return nil, fmt.Errorf("LLM service call failed: %v", err)
	}

	fmt.Println("LLM Raw Response:", respText)

	// 4️⃣ Extract JSON from LLM response
	rawJSON, err := globalFunctions.ExtractJson(respText)
	if err != nil {
		return nil, fmt.Errorf("failed to extract JSON from LLM response: %v", err)
	}

	// 5️⃣ Unmarshal JSON into Ensights struct
	var finalEnsight Ensights
	if err := json.Unmarshal([]byte(rawJSON), &finalEnsight); err != nil {
		return nil, fmt.Errorf("failed to unmarshal LLM JSON response: %v", err)
	}

	// ✅ Return the final aggregated insight
	return &finalEnsight, nil
}

func LoadAllInsightsFromJSON() ([]Ensights, error) {
	fileName := "insights_data.json"

	data, err := os.ReadFile(fileName)
	if err != nil {
		// If file not found → return empty
		if os.IsNotExist(err) {
			return []Ensights{}, nil
		}
		return nil, err
	}

	var insights []Ensights
	if len(data) > 0 {
		if err := json.Unmarshal(data, &insights); err != nil {
			return nil, err
		}
	}

	return insights, nil
}
// func BuildFinalInsightsQuery(allInsights []Ensights, input ApiInputParams) string {
// 	var b strings.Builder

// 	b.WriteString("You are an expert Voice Insights Engine for IndiaMART.\n")
// 	b.WriteString("You are given multiple individual insights from different calls.\n")
// 	b.WriteString("Your job is to generate ONE FINAL aggregated insight block.\n\n")

// 	b.WriteString("### INPUT INSIGHTS ###\n")

// 	for i, ins := range allInsights {
// 		if input.MaxCallLimit < i+1 {
// 			break
// 		}
// 		b.WriteString(fmt.Sprintf("CALL %d:\n", i+1))
// 		b.WriteString(fmt.Sprintf("- Concerns: %s\n", ins.Concerns))
// 		b.WriteString(fmt.Sprintf("- Resolution: %s\n", ins.Resolution))
// 		b.WriteString(fmt.Sprintf("- NextSteps: %s\n", ins.NextSteps))
// 		b.WriteString(fmt.Sprintf("- Alert: %s\n", ins.Alert))
// 		b.WriteString(fmt.Sprintf("- Sentiment: %s\n", ins.Sentiment))
// 		b.WriteString(fmt.Sprintf("- KeyPoints: %s\n\n", ins.KeyPoints))
// 	}

// 	b.WriteString("### INSTRUCTIONS ###\n")
// 	b.WriteString("- Combine and summarize all concerns.\n")
// 	b.WriteString("- Detect repeated patterns across calls.\n")
// 	b.WriteString("- Identify root causes.\n")
// 	b.WriteString("- Generate high-quality Resolution and NextSteps.\n")
// 	b.WriteString("- Output only ONE insight block.\n")
// 	b.WriteString("- EnsightType must be 'final'.\n")
// 	b.WriteString("- Output strictly valid JSON matching this structure:\n\n")

// 	ex := Ensights{
// 		EnsightType: "final",
// 		Concerns:    "Summary of all major repeated concerns",
// 		Resolution:  "Combined resolution from all calls",
// 		NextSteps:   "Clear next actions for Executive & Customer",
// 		Alert:       "Any critical risk or upsell opportunity",
// 		Sentiment:   "Overall sentiment trend",
// 		KeyPoints:   "keywords",
// 	}
// 	exampleBytes, _ := json.MarshalIndent(ex, "", "  ")

// 	b.WriteString(string(exampleBytes) + "\n\n")
// 	b.WriteString("Return ONLY JSON.")

// 	return b.String()
// }
func BuildFinalInsightsQuery(allInsights []Ensights, input ApiInputParams) string {
    var b strings.Builder

    b.WriteString("You are an expert, quantitative Voice Insights Engine who generates Insight for IndiaMART Higher Authorities (VP, Product Manager, Senior Sales/Ops Manager).\n")
    b.WriteString("You are given multiple individual, call-level insights for the same or different seller.\n")
    b.WriteString("Your job is to generate ONE FINAL aggregated insight block that is QUANTITATIVE and ACTIONABLE.\n\n")

    // --- 1. INPUT INSIGHTS ---
    b.WriteString("### RAW CALL INSIGHTS ###\n")

    // Input loop remains the same to feed all prior insights to the model
    for i, ins := range allInsights {
        if input.MaxCallLimit < i+1 {
            break
        }
        b.WriteString(fmt.Sprintf("CALL %d:\n", i+1))
        b.WriteString(fmt.Sprintf("- Concerns: %s\n", ins.Concerns))
        b.WriteString(fmt.Sprintf("- Resolution: %s\n", ins.Resolution))
        b.WriteString(fmt.Sprintf("- NextSteps: %s\n", ins.NextSteps))
        b.WriteString(fmt.Sprintf("- Alert: %s\n", ins.Alert))
        b.WriteString(fmt.Sprintf("- Sentiment: %s\n", ins.Sentiment))
        b.WriteString(fmt.Sprintf("- KeyPoints: %s\n\n", ins.KeyPoints))
    }
    
    // --- 2. INSTRUCTIONS FOR QUANTITATIVE AGGREGATION ---
    b.WriteString("### QUANTITATIVE AGGREGATION INSTRUCTIONS ###\n")
    b.WriteString("- Analyze all calls to identify recurring issues, failure rates, and opportunities.\n")
    b.WriteString("- The output MUST be QUANTITATIVE, using percentages or counts (e.g., 20% of calls, 4/10 cases).\n")
    b.WriteString("- Ensure NextSteps and Resolutions are targeted at specific personas (Executive, Manager, Sales Head, Product, Category).\n")
    b.WriteString("- Output only ONE insight block with EnsightType = 'final'.\n")
    b.WriteString("- All pointers in the final block MUST be concise and in a list/bullet format.\n\n")

    // --- 3. FIELD MAPPING TO QUANTITATIVE FORMAT ---
    b.WriteString("### FIELD LOGIC (Quantitative Format) ###\n")
    b.WriteString("Concerns: Summarise all recurring issues with their quantitative recurrence (e.g., 'Irrelevant Buyleads (40% of calls)').\n")
    b.WriteString("Resolution: Summarize resolutions given, identifying systemic failures or best practices quantitatively.\n")
    b.WriteString("NextSteps: Define clear actionables for relevant personas (Executive, Product Manager, Sales Manager, etc.).\n")
    b.WriteString("Alert: Aggregate all alerts and state the concerned person/team.\n")
    b.WriteString("Sentiment: Summarise the sentiment distribution quantitatively (e.g., '60% Negative, 30% Neutral, 10% Positive').\n")
    b.WriteString("KeyPoints: List all distinct keywords with their total occurrence count across all calls (e.g., 'Buyleads (10), Catalogue (4)').\n")


    // --- 4. JSON STRUCTURE (Quantitative Example) ---
    ex := Ensights{
        EnsightType: "final",
        Concerns:    "Irrelevant Buyleads (40% of calls); Product Catalogue issues (30% of calls)",
        Resolution:  "4/10 executives failed to update category settings; Need automated bulk lead filter guidance.",
        NextSteps:   "Executive: 4/10 need training on lead qualification; Category Manager: Advise categories for SEO work; Product Manager: Simplify catalogue upload journey.",
        Alert:       "20% Upsell Opportunity cases identified for Sales Manager follow-up; 10% Executive Inefficiency (rude behavior).",
        Sentiment:   "Negative (60%) -> Neutral (30%) -> Positive (10%)",
        KeyPoints:   "Buyleads (10), Catalogue (4), Upsell (2), Inefficiency (1)",
    }
    exampleBytes, _ := json.MarshalIndent(ex, "", "  ")

    b.WriteString("\n### OUTPUT JSON MUST STRICTLY FOLLOW THIS QUANTITATIVE STRUCTURE ###\n")
    b.WriteString(string(exampleBytes) + "\n\n")
    b.WriteString("Return ONLY JSON.")

    return b.String()
}
