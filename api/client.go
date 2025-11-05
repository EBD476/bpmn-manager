package api

import (
	"bpmn-manager/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type APIClient struct {
	baseURL    string
	httpClient *http.Client
	authToken  string
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *APIClient) SetAuthToken(token string) {
	c.authToken = token
}

func (c *APIClient) doRequest(method, endpoint string) ([]byte, error) {
	url := c.baseURL + endpoint

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "BPMN-Manager-CLI/1.0")
	if c.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.authToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	return body, nil
}

// func (c *APIClient) GetUserTasks() ([]models.UserTask, error) {

// 	body, err := c.doRequest("GET", "/api/user/tasks")
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Try to parse as direct array first (your API format)
// 	var rawTasks []map[string]interface{}
// 	if err := json.Unmarshal(body, &rawTasks); err == nil {
// 		fmt.Printf("✅ Successfully parsed %d user tasks (direct array format)\n", len(rawTasks))
// 		return c.parseRawTasks(rawTasks), nil
// 	}

// 	// If direct array parsing fails, try other formats
// 	return c.tryAlternativeFormats(body)

// 	// var response models.UserTasksResponse
// 	// if err := json.Unmarshal(body, &response); err != nil {
// 	// 	return nil, fmt.Errorf("failed to parse user tasks: %v", err)
// 	// }

// 	// return response.Tasks, nil
// }

func (c *APIClient) GetUserTasks() ([]models.UserTask, error) {

	// url := "http://192.168.164.150:8086/api/user/tasks"

	// resp, err := http.Get(url)
	// if err != nil {
	// 	fmt.Println("Error making request:", err)
	// 	return nil, nil
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	fmt.Println("Request failed with status:", resp.Status)
	// 	return nil, nil
	// }

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println("Error reading response:", err)
	// 	return nil, nil
	// }

	// // Parse JSON into slice of UserTask
	// var tasks []models.UserTask
	// err = json.Unmarshal(body, &tasks)
	// if err != nil {
	// 	fmt.Println("Error parsing JSON:", err)
	// 	return nil, nil
	// }

	body, err := c.doRequest("GET", "/api/user/tasks")
	if err != nil {
		return nil, err
	}

	var response []models.UserTask
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse user tasks: %v", err)
	}

	return response, nil
}

func (c *APIClient) getStringField(data map[string]interface{}, field string) string {
	if value, exists := data[field]; exists {
		if str, ok := value.(string); ok {
			return str
		}
		// Handle null values
		if value == nil {
			return ""
		}
		// Try to convert other types to string
		return fmt.Sprintf("%v", value)
	}
	return ""
}

func (c *APIClient) GetRunningProcesses() ([]models.RunningProcess, error) {
	// body, err := c.doRequest("GET", "/api/running-processes")
	body, err := c.doRequest("GET", "/api/all-instances")

	if err != nil {
		return nil, err
	}

	var response []models.RunningProcess
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse running processes: %v", err)
	}

	return response, nil
}

func (c *APIClient) GetProcessDetails(processID string) (*models.ProcessDetails, error) {

	endpoint := fmt.Sprintf("/api/%s/details", processID)
	body, err := c.doRequest("GET", endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse process details: %v", err)
		// return nil, err
	}

	var response models.ProcessDetails
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse process details: %v", err)
	}

	return &response, nil
}

func (c *APIClient) GetCompletedProcesses() ([]models.ProcessDetails, error) {
	body, err := c.doRequest("GET", "/api/completed-processes")
	if err != nil {
		return nil, err
	}

	var response []models.ProcessDetails
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse running processes: %v", err)
	}

	return response, nil
}

func (c *APIClient) GetCompletedTasks() ([]models.UserTask, error) {
	body, err := c.doRequest("GET", "/api/completed-tasks")
	if err != nil {
		return nil, err
	}

	var response []models.UserTask
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse running processes: %v", err)
	}

	return response, nil
}

func (c *APIClient) CompleteTask(taskID string, formData models.FormData) error {

	endpoint := fmt.Sprintf("/api/complete-task/%s", taskID)
	url := c.baseURL + endpoint

	// Create a TaskCompletionRequest object (you can adjust the payload based on your API)
	// payload := models.TaskCompletionRequest{
	// 	TaskID:           taskID,
	// 	ReDefineDecision: false,
	// 	DbDecision:       false,
	// 	Status:           "completed", // Assuming the status is what marks the task as completed
	// }

	// Marshal the payload to JSON
	jsonData, err := json.Marshal(formData)
	if err != nil {
		return fmt.Errorf("failed to marshal request payload: %w", err)
	}

	// Send a POST request to the API
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set the content type to JSON
	req.Header.Set("Content-Type", "application/json")

	// Send the request using the client
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Print the response body (for debugging or informational purposes)
	// fmt.Printf("Response Body: %s\n", string(respBody))
	// Check if the request was successful (HTTP Status Code 200-299)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to complete task, status code: %s %s", string(respBody), jsonData)
		// return fmt.Errorf("failed to complete task, status code: %d", resp.StatusCode)
	}

	// if err := json.Unmarshal(body, &response); err != nil {
	// 	return nil, fmt.Errorf("failed to parse process details: %v", err)
	// }

	// Check if the request was successful (HTTP Status Code 200-299)
	// if body.StatusCode < 200 || resp.StatusCode >= 300 {
	// 	return fmt.Errorf("failed to complete task, status code: %d", resp.StatusCode)
	// }

	// Optionally, process the response body (e.g., check success or return data)
	// Here we just print the response for now
	// fmt.Printf("Task completed successfully.")

	return nil
}

func (c *APIClient) StartProcess() error {
	body, err := c.doRequest("GET", "/api/start")
	if err != nil {
		return nil
	}

	var response []models.RunningProcess
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse running processes: %v", err)
	}

	return nil
}
