package store

import "encoding/json"

func Encode(value any) ([]byte, error) {
	return json.Marshal(value)
}

func Decode(data []byte, value any) error {
	return json.Unmarshal(data, value)
}

func SaveJSON(s *FileStore, key string, value any) error {
	data, err := Encode(value)
	if err != nil {
		return err
	}
	return s.Save(key, data)
}

func LoadJSON(s *FileStore, key string, value any) error {
	data, err := s.Load(key)
	if err != nil {
		return err
	}
	return Decode(data, value)
}
