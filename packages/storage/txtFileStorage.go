package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

//type history struct {
//	value string
//	updatedAt time.Time
//}

type fileContent struct {
	PreviousValues []string `json:"previous_values"` // TODO: make it from type history
	LatestValue    string   `json:"latest_value"`
}

type txtFileStorage struct {
	mutex sync.Mutex
}

func (f *txtFileStorage) Save(key string, value []byte) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	fileContent := &fileContent{}

	previousValues, err := f.Get(key)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("-text storage: Error getting previousValues", err)
		return err
	}

	err = json.Unmarshal(previousValues, fileContent)
	if err != nil {
		fmt.Println("-text storage: error unmarshalling previousValues:", err)
		return err
	}

	if fileContent.LatestValue != "" {
		fileContent.PreviousValues = append(fileContent.PreviousValues, fileContent.LatestValue)
	}
	fileContent.LatestValue = string(value)

	binaryContent, err := json.Marshal(fileContent)
	if err != nil {
		fmt.Println("-text storage: error marshalling previousValues:", err)
		return err
	}

	err = os.WriteFile(keyToFilename(key), binaryContent, os.ModePerm)
	if err != nil {
		fmt.Println("-text storage: error writing previousValues:", err)
		return err
	}
	return nil
}

func (f *txtFileStorage) Get(key string) ([]byte, error) {
	fileContent, err := os.ReadFile(keyToFilename(key))
	if err != nil {
		if os.IsNotExist(err) {
			return []byte("{}"), nil
		}
		fmt.Println("-text storage: error reading value with key", key, ":", err)
		return nil, err
	}
	return fileContent, nil
}

func keyToFilename(str string) string {
	return fmt.Sprintf("%s.txt", str)
}
