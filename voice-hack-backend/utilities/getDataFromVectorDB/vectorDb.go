package getdatafromvectordb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	pc "github.com/pinecone-io/go-pinecone/pinecone"
)

// ===============================
//        CONFIG
// ===============================
var (
	LLM_GATEWAY_KEY = "sk-gDgbFkH2KH-HKvX2a4uEXw" // Replace with your key
	PINECONE_KEY    = "pcsk_6VcVxq_QvKbhRAvWRg8iAJDUdY6STdzHaRckmwRTsqPrEQs7HfuudesdDsAFQiUw1kqovh"
	PINECONE_HOST   = "https://insights-fv7quol.svc.aped-4627-b74a.pinecone.io"
	PINECONE_INDEX  = "insights"
	NAMESPACE       = "default"
)

// ===============================
//      EMBEDDING FUNCTION
// ===============================
func GetEmbedding(text string) ([]float32, error) {
	url := "https://imllm.intermesh.net/v1/embeddings"
	body := map[string]interface{}{
		"model": "google/gemini-embedding-001",
		"input": text,
	}
	bodyBytes, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+LLM_GATEWAY_KEY)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	var respJSON struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
		Error interface{} `json:"error"`
	}

	json.Unmarshal(data, &respJSON)

	if respJSON.Error != nil {
		return nil, fmt.Errorf("embedding failed: %v", respJSON.Error)
	}

	if len(respJSON.Data) > 0 {
		return respJSON.Data[0].Embedding, nil
	}

	return nil, fmt.Errorf("no embedding returned")
}

// ===============================
//      VECTOR DB STRUCT
// ===============================
type VectorDB struct {
	Client *pc.Client
	Index  *pc.IndexConnection
}

// NewVectorDB initializes Pinecone client and index connection
func NewVectorDB() *VectorDB {
	// ctx := context.Background()

	client, err := pc.NewClient(pc.NewClientParams{
		ApiKey: PINECONE_KEY,
		Host:   PINECONE_HOST,
	})
	if err != nil {
		log.Fatalf("failed to create Pinecone client: %v", err)
	}

	idxConn, err := client.Index(pc.NewIndexConnParams{
		Host:      PINECONE_HOST,
		Namespace: NAMESPACE,
	})
	if err != nil {
		log.Fatalf("failed to connect to Pinecone index: %v", err)
	}

	return &VectorDB{
		Client: client,
		Index:  idxConn,
	}
}

// QueryText queries Pinecone by a text string and returns metadata as string list
func (vdb *VectorDB) QueryText(query string, topK int) ([]string, error) {
	// Step 1: get embedding
	embedding, err := GetEmbedding(query)
	if err != nil {
		return nil, fmt.Errorf("embedding failed: %v", err)
	}

	// Step 2: query Pinecone
	resp, err := vdb.Index.QueryByVectorValues(context.Background(), &pc.QueryByVectorValuesRequest{
		Vector:          embedding,
		TopK:            uint32(topK),
		IncludeMetadata: true,
		IncludeValues:   false,
	})
	if err != nil {
		return nil, fmt.Errorf("pinecone query failed: %v", err)
	}

	// Step 3: collect metadata as list of strings
	results := []string{}
	for i, match := range resp.Matches {
		if match.Vector != nil && match.Vector.Metadata != nil {
			summary := ""
			transcript := ""

			if val, ok := match.Vector.Metadata.Fields["summary"]; ok {
				summary = val.GetStringValue()
			}
			if val, ok := match.Vector.Metadata.Fields["transcript"]; ok {
				transcript = val.GetStringValue()
			}

			resultStr := fmt.Sprintf("sample_call_%d\nsummary: %q\ntranscript: %q", i+1, summary, transcript)
			results = append(results, resultStr)
		} else {
			// fallback if metadata missing
			results = append(results, "sample_call_"+strconv.Itoa(i+1)+": metadata missing")
		}
	}

	return results, nil
}


func GetSampleCalls(text string, topK int) (string, string) {
	vdb := NewVectorDB()

	results, err := vdb.QueryText(text, topK)
	if err != nil {
		return "", fmt.Sprintf("Query failed: %v", err)
	}
	ans := ""
	for _, res := range results {
		ans += res + "\n-------------------------------\n"
	}
	return ans, ""
}

