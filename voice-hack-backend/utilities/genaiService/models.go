package genaiService

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"google.golang.org/genai"
)

func ListAvailableModels(ginCtx *gin.Context, client *genai.Client) (availableModels []string, err error) {
	modelList, modelListErr := client.Models.List(ginCtx, &genai.ListModelsConfig{})
	if modelListErr != nil {
		err = modelListErr
		return
	}
	for _, model := range modelList.Items {
		for _, action := range model.SupportedActions {
			if action == "generateContent" {
				availableModels = append(availableModels, model.Name)
			}
		}
	}
	if len(availableModels) == 0 {
		err = fmt.Errorf("no suitable model found")
		return
	}
	return
}

func IsModelAvailable(ginCtx *gin.Context, client *genai.Client, model string) (err error) {
	availableModels, availableModelsErr := ListAvailableModels(ginCtx, client)
	if availableModelsErr != nil {
		return availableModelsErr
	}
	for _, m := range availableModels {
		if "models/"+model == m {
			return nil
		}
	}
	return fmt.Errorf("model is not available")
}
