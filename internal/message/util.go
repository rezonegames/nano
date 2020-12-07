package message

import "nano/internal/env"

func Serialize(v interface{}) ([]byte, error) {
	if data, ok := v.([]byte); ok {
		return data, nil
	}
	data, err := env.Serializer.Marshal(v)
	if err != nil {
		return nil, err
	}
	return data, nil
}
