package insightsGenerateController

import (
	"fmt"
	"net/http"
	insightsGenerateModel "voice-hack-backend/modules/tripPlanner/model/insightsGenerateModel"
	urlMedia "voice-hack-backend/utilities/urlMedia"

	"github.com/gin-gonic/gin"
)


func FinalSummaryGenerate(ginCtx *gin.Context) {
	// Bind Input
	apiInputParam, bindErr := BindInputParams(ginCtx)
	apiResponse := insightsGenerateModel.ApiResponse{}

	// Always create logs at the end
	defer func() {
		insightsGenerateModel.CreateApplicationLogs(ginCtx, apiInputParam, apiResponse)
	}()

	if bindErr != nil {
		apiResponse.Code = http.StatusBadRequest
		apiResponse.Status = "Failure"
		apiResponse.Error = bindErr.Error()
		ReturnApiResponse(ginCtx, http.StatusBadRequest, apiResponse)
		return
	}

	// Generate final summary from stored JSON
	resp, respErr := insightsGenerateModel.FinalSummary(apiInputParam)
	if respErr != nil {
		apiResponse.Code = http.StatusInternalServerError
		apiResponse.Status = "Failure"
		apiResponse.Error = respErr.Error()
		ReturnApiResponse(ginCtx, http.StatusInternalServerError, apiResponse)
		return
	}

	// Safety check: ensure resp is not nil
	if resp == nil {
		apiResponse.Code = http.StatusOK
		apiResponse.Status = "Success"
		apiResponse.Error = "No insights generated"
		ReturnApiResponse(ginCtx, http.StatusOK, apiResponse)
		return
	}

	// Prepare response
	apiResponse.Response.Locations = []insightsGenerateModel.Ensights{*resp}
	apiResponse.Code = http.StatusOK
	apiResponse.Status = "Success"

	ReturnApiResponse(ginCtx, http.StatusOK, apiResponse)
}

func FindDestination(ginCtx *gin.Context) {
	apiInputParam, bindErr := BindInputParams(ginCtx)
	apiResponse := insightsGenerateModel.ApiResponse{}
	defer func() {
		insightsGenerateModel.CreateApplicationLogs(ginCtx, apiInputParam, apiResponse)
	}()

	if bindErr != nil {
		apiResponse.Code = http.StatusBadRequest
		apiResponse.Status = "Failure"
		apiResponse.Error = bindErr.Error()
		ReturnApiResponse(ginCtx, http.StatusBadRequest, apiResponse)
		return
	}
	// Call URL Media to get text from call recording URL
	apiInputParam.TrascriptionURLTxt = make([]string, len(apiInputParam.CallData))

	for index, callDate := range apiInputParam.CallData {
		setInputParam := urlMedia.SetInputParamForTranscribeAPI(callDate.CallRecordingURL, apiInputParam.ExecutiveID, fmt.Sprint(apiInputParam.Glid))
		fmt.Println(setInputParam)
		getTextFromCallURl, err := urlMedia.CallTranscribeAPI(setInputParam)
		if err != "" {
			apiResponse.Code = http.StatusInternalServerError
			apiResponse.Status = "Failure"
			apiResponse.Error = err
			ReturnApiResponse(ginCtx, http.StatusBadRequest, apiResponse)
			return
		}
		fmt.Print(getTextFromCallURl.TranscriptionURL)
		text, errors := urlMedia.GetTextFromURL(getTextFromCallURl.TranscriptionURL)

		if errors != nil {
			apiResponse.Code = http.StatusInternalServerError
			apiResponse.Status = "Failure"
			apiResponse.Error = "GettextFrom url failed" + errors.Error()
			ReturnApiResponse(ginCtx, http.StatusBadRequest, apiResponse)
			return
		}
		apiInputParam.TrascriptionURLTxt[index] = text

		// fmt.Print(text)
	}

	userQuery := insightsGenerateModel.GenerateUserQuery(apiInputParam)
	// resp, respErr := insightsGenerateModel.GenerateInsights(ginCtx, userQuery, apiInputParam)
	resp, respErr := insightsGenerateModel.GenerateInsightsFromLLM(ginCtx, userQuery, apiInputParam)
	if respErr != nil {
		apiResponse.Code = http.StatusInternalServerError
		apiResponse.Status = "Success"
		apiResponse.Error = respErr.Error()
		ReturnApiResponse(ginCtx, http.StatusBadRequest, apiResponse)
		return
	}
	insightsGenerateModel.StoreInsightsToJSON(resp.Locations)

	apiResponse.Code = http.StatusOK
	apiResponse.Status = "Success"
	apiResponse.Response = resp
	ReturnApiResponse(ginCtx, http.StatusOK, apiResponse)
}

func BindInputParams(ginCtx *gin.Context) (InputParams insightsGenerateModel.ApiInputParams, err error) {
	bindErr := ginCtx.ShouldBindBodyWithJSON(&InputParams)
	return InputParams, bindErr
}

func ReturnApiResponse(ginCtx *gin.Context, apiCode int, apiResponse insightsGenerateModel.ApiResponse) {
	ginCtx.JSON(apiCode, apiResponse)
}
