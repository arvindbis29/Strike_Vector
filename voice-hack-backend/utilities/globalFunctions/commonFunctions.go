package globalFunctions

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetFileAndLineNo(skipValue int) (string, int) {
	_, file, line, ok := runtime.Caller(skipValue)

	if !ok {
		file = "file_not_found"
		line = 0
	}

	return file, line
}

func GetCurrentTimeInMs() string {
	currentMicrosec := int(time.Now().UnixNano() / int64(time.Microsecond))
	currentSec := int(time.Now().UnixNano() / int64(time.Second))
	leftOverLength := len(ConvertValueToString(currentMicrosec)) - len(ConvertValueToString(currentSec))

	currentMs := (currentMicrosec - ConvertValueToInt(ConvertValueToString(currentSec)+strings.Repeat("0", leftOverLength))) / 1000

	currentTime := time.Now().Format("15:04:05") + "." + fmt.Sprintf("%03d", currentMs)

	return currentTime

}

func ConvertValueToString(value any) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}

func ConvertValueToInt(inputValue any) int {

	var convertedVal int
	switch val := inputValue.(type) { //in go switch dont require break statement so no need to apply
	case int64:
		convertedVal = int(val)
	case int32: //rune and int32 are same
		convertedVal = int(val)
	case int16:
		convertedVal = int(val)
	case int8:
		convertedVal = int(val)
	case int:
		convertedVal = val
	case uint:
		convertedVal = int(val)
	case uint64:
		convertedVal = int(val)
	case uint32:
		convertedVal = int(val)
	case uint16:
		convertedVal = int(val)
	case uint8: //byte and uint8 are same
		convertedVal = int(val)
	case []byte: //byte and uint8 are same
		if strToFloatVal, err := strconv.Atoi(string(val)); err == nil {
			convertedVal = strToFloatVal
		} else {
			convertedVal = 0
		}
	case string:
		if strToFloatVal, err := strconv.Atoi(val); err == nil {
			convertedVal = strToFloatVal
		} else {
			convertedVal = 0
		}
	case float64:
		return int(val)
	case float32:
		return int(val)
	case json.Number:
		if numVal, err := val.Int64(); err == nil {
			convertedVal = int(numVal)
		} else if numVal, err := val.Float64(); err == nil {
			convertedVal = int(numVal)
		} else {
			convertedVal = 0
		}
	case bool:
		if val {
			convertedVal = 1
		} else {
			convertedVal = 0
		}
	default:
		convertedVal = 0
	}
	return convertedVal
}

func ConvertJsonValToString(inputValue any) string {

	if inputValue == nil {
		return ""
	}

	var result string

	switch val := inputValue.(type) {
	case string:
		result = val
	case []byte:
		result = string(val)
	case int:
		result = strconv.Itoa(val)
	default:
		result = ConvertValueToString(val)
	}
	return result

}

func WriteJsonLogs(ginCtx *gin.Context, fileName string, logData map[string]any) {

	if ginCtx != nil {
		if hostServer := ginCtx.Request.Host; hostServer != "" {
			logData["server_host"] = hostServer
		}
	}

	if fileName == "" {
		fileName = "trip_planner_logs"
	}
	logFilePath := "centralLogging/"
	if err := os.MkdirAll(logFilePath, 0755); err != nil {
		fmt.Printf("failed to create log dir: %v\n", err)
		return
	}

	finalPath := filepath.Join(logFilePath, fileName+"_"+time.Now().Format("2006_01_02")+".log")

	// the data in the file

	fileHandler, fileOpenErr := os.OpenFile(finalPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if fileOpenErr != nil {
		// return fileOpenErr.Error()
		return
	}
	defer fileHandler.Close()

	finalLogStr, err := ConvertValueToJson(logData)
	if err != "" {
		fmt.Printf("failed to convert log data: %v\n", err)
		return
	}

	if _, err := fmt.Fprintf(fileHandler, "%s\n", string(finalLogStr)); err != nil {
		fmt.Printf("failed to write log: %v\n", err)
	}

}

func ConvertValueToJson(inputValue any) ([]byte, string) {
	res, err := json.Marshal(inputValue)

	if err != nil {
		return []byte{}, err.Error()
	}
	return res, ""
}

func ExtractJson(rawMessage string) (string, error) {

	startIndex := strings.IndexAny(rawMessage, "{[")

	if startIndex == -1 {
		return "", errors.New("no json start Index found")
	}

	endIndex := strings.LastIndexAny(rawMessage, "]}")

	if endIndex == -1 {
		return "", fmt.Errorf("no json last index found")

	}

	newRawMessage := rawMessage[startIndex : endIndex+1]

	var js any
	if err := json.Unmarshal([]byte(newRawMessage), &js); err != nil {
		return "", fmt.Errorf("nvalid JSON after trimming: %w", err)
	}

	return newRawMessage, nil
}
