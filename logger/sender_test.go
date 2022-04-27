package logger

import (
	"encoding/json"
	"testing"
)

func TestRequest(t *testing.T) {
	t.Run("test empty MachineID", func(t *testing.T) {
		reqData := request{
			Instance:   "worker_node",
			LogContent: "content",
			MachineID:  "",
		}
		expectedResult := `{"instance":"worker_node","log_content":"content"}`

		requestBody, err := json.Marshal(reqData)
		if err != nil {
			t.Error(err)
		}
		requestBodyStr := string(requestBody)

		if requestBodyStr != expectedResult {
			t.Fatalf("want: %s, got: %s", expectedResult, requestBodyStr)
		}
	})
}
