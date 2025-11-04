package storage

import (
	"bpmn-manager/models"
	"encoding/json"
	"os"
	"path/filepath"
)

type Storage struct {
	dataDir string
}

func NewStorage(dataDir string) *Storage {
	os.MkdirAll(dataDir, 0755)
	return &Storage{dataDir: dataDir}
}

func (s *Storage) SaveProcess(process *models.Process) error {
	filename := filepath.Join(s.dataDir, "process_"+process.ID+".json")
	data, err := json.MarshalIndent(process, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (s *Storage) LoadProcess(processID string) (*models.Process, error) {
	filename := filepath.Join(s.dataDir, "process_"+processID+".json")
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var process models.Process
	err = json.Unmarshal(data, &process)
	return &process, err
}

func (s *Storage) ListProcesses() ([]*models.Process, error) {
	files, err := filepath.Glob(filepath.Join(s.dataDir, "process_*.json"))
	if err != nil {
		return nil, err
	}

	var processes []*models.Process
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		var process models.Process
		if err := json.Unmarshal(data, &process); err == nil {
			processes = append(processes, &process)
		}
	}
	return processes, nil
}

func (s *Storage) SaveInstance(instance *models.ProcessInstance) error {
	filename := filepath.Join(s.dataDir, "instance_"+instance.ID+".json")
	data, err := json.MarshalIndent(instance, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (s *Storage) LoadInstance(instanceID string) (*models.ProcessInstance, error) {
	filename := filepath.Join(s.dataDir, "instance_"+instanceID+".json")
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var instance models.ProcessInstance
	err = json.Unmarshal(data, &instance)
	return &instance, err
}

func (s *Storage) ListInstances() ([]*models.ProcessInstance, error) {
	files, err := filepath.Glob(filepath.Join(s.dataDir, "instance_*.json"))
	if err != nil {
		return nil, err
	}

	var instances []*models.ProcessInstance
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		var instance models.ProcessInstance
		if err := json.Unmarshal(data, &instance); err == nil {
			instances = append(instances, &instance)
		}
	}
	return instances, nil
}
