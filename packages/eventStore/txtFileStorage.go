package eventStore

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
	Events       []string `json:"events"` // TODO: make it from type history
	CurrentState string   `json:"current_state"`
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
		fmt.Println("-eventStore: Error getting previousValues", err)
		return err
	}

	err = json.Unmarshal(previousValues, fileContent)
	if err != nil {
		fmt.Println("-eventStore: error unmarshalling previousValues:", err)
		return err
	}

	if fileContent.CurrentState != "" {
		fileContent.Events = append(fileContent.Events, fileContent.CurrentState)
	}
	fileContent.CurrentState = string(value)

	binaryContent, err := json.Marshal(fileContent)
	if err != nil {
		fmt.Println("-eventStore: error marshalling previousValues:", err)
		return err
	}

	err = os.WriteFile(keyToFilename(key), binaryContent, os.ModePerm)
	if err != nil {
		fmt.Println("-eventStore: error writing previousValues:", err)
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
		fmt.Println("-eventStore: error reading value with key", key, ":", err)
		return nil, err
	}
	return fileContent, nil
}

func keyToFilename(str string) string {
	return fmt.Sprintf("%s.txt", str)
}
