package models

import (
	"encoding/xml"
	"time"
)

// RawUserTask matches your exact API response format
type RawUserTask struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Assignee interface{} `json:"assignee"` // Use interface{} to handle null
}

// API Response Structures
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type UserTasksResponse struct {
	Tasks []UserTask `json:"tasks"`
}

type RunningProcessesResponse struct {
	Processes []RunningProcess `json:"processes"`
}

type ProcessDetailsResponse struct {
	Process ProcessDetails `json:"process"`
}

// BPMN Types
type Process struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Version     string     `json:"version"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Activities  []Activity `json:"activities"`
}

type Activity struct {
	ID          string    `json:"id"`
	TaskID      string    `json:"task_id"`
	ProcessID   string    `json:"process_id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"` // task, gateway, event
	Description string    `json:"description"`
	Assignee    string    `json:"assignee"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"` // pending, in-progress, completed
	CreatedAt   time.Time `json:"created_at"`
}

type ProcessInstance struct {
	ID              string            `json:"id"`
	ProcessID       string            `json:"process_id"`
	ProcessName     string            `json:"process_name"`
	CurrentActivity string            `json:"current_activity"`
	Status          string            `json:"processStatus"`
	StartedAt       time.Time         `json:"started_at"`
	Variables       map[string]string `json:"variables"`
}

// API Specific Types
type UserTask struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	ProcessID         string    `json:"processInstanceId"`
	ProcessName       string    `json:"process_name"`
	ActivityName      string    `json:"activity_name"`
	Assignee          string    `json:"assignee"`
	DueDate           time.Time `json:"due_date"`
	Status            string    `json:"processStatus"`
	Priority          string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	TaskDefinitionKey string    `json:"taskDefinitionKey"`
}

type RunningProcess struct {
	InstanceID           string    `json:"instance_id"`
	ProcessID            string    `json:"id"`
	ProcessDefinitionId  string    `json:"processDefinitionId"`
	ProcessDefinitionKey string    `json:"processDefinitionKey"`
	ProcessName          string    `json:"process_name"`
	CurrentActivity      string    `json:"current_activity"`
	StartTime            time.Time `json:"startTime"`
	Duration             string    `json:"duration"`
	Status               string    `json:"processStatus"`
}

type ProcessDetails struct {
	ID                   string                 `json:"id"`
	ProcessDefinitionId  string                 `json:"processDefinitionId"`
	ProcessDefinitionKey string                 `json:"processDefinitionKey"`
	StartTime            string                 `json:"startTime"`
	EndTime              string                 `json:"endTime"`
	CurrentVariables     map[string]interface{} `json:"currentVariables"`
	Name                 string                 `json:"name"`
	Description          string                 `json:"description"`
	Version              string                 `json:"version"`
	CreatedAt            time.Time              `json:"created_at"`
	Duration             int                    `json:"duration"`
	Activities           []ProcessActivity      `json:"activities"`
	Instances            []ProcessInstance      `json:"instances"`
	Statistics           ProcessStatistics      `json:"statistics"`
}

type ProcessActivity struct {
	ID          string    `json:"activityId"`
	TaskId      string    `json:taskId`
	Name        string    `json:"activityName"`
	Type        string    `json:"activityType"`
	Description string    `json:"description"`
	Assignee    string    `json:"assignee"`
	Order       int       `json:"order"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Duration    int       `json:"duration"`
}

type ProcessStatistics struct {
	TotalInstances    int    `json:"total_instances"`
	Running           int    `json:"running"`
	Completed         int    `json:"completed"`
	Failed            int    `json:"failed"`
	AvgCompletionTime string `json:"avg_completion_time"`
}

// TaskCompletionRequest struct for task completion payload (if needed)
type TaskCompletionRequest struct {
	TaskID           string `json:"taskID"`
	Status           string `json:"status"` // Example: "completed"
	ReDefineDecision bool   `json:"reDefineDecision"`
	DbDecision       bool   `json:"dbDecision"`
}

// BPMN XML Structure (simplified)
type Definitions struct {
	XMLName   xml.Name      `xml:"definitions"`
	Processes []BPMNProcess `xml:"process"`
}

type BPMNProcess struct {
	XMLName xml.Name `xml:"process"`
	ID      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
	Tasks   []Task   `xml:"task"`
}

type Task struct {
	XMLName xml.Name `xml:"task"`
	ID      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
}

type FormData struct {
	ReDefineDecision  bool   `json:"reDefineDecision"`
	DbDecision        bool   `json:"dbDecision"`
	BusinessApproved  bool   `json:"businessApproved"`
	TechnicalApproved bool   `json:"technicalApproved"`
	OperationApproved bool   `json:"operationApproved"`
	Comment           string `json:"comment"`
	Message           string `json:"message"`
}
