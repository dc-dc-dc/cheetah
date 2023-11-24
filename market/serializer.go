package market

import (
	"encoding/json"
	"fmt"
)

type SerializableReceiverGenerator func() SerializableReceiver

var (
	serializableMap = map[string]SerializableReceiverGenerator{}
)

func RegisterSerializableReceiver(key string, generator SerializableReceiverGenerator) {
	serializableMap[key] = generator
}

func GetSerializableReceivers(receivers []MarketReceiver) []SerializableReceiver {
	res := make([]SerializableReceiver, 0)
	for _, receiver := range receivers {
		switch receiver.(type) {
		case SerializableReceiver:
			{
				serializableReceiver := receiver.(SerializableReceiver)
				res = append(res, serializableReceiver)
				break
			}
		}
	}
	return res
}

type serializableReceiverJSON struct {
	Key string `json:"key"`
	Raw []byte `json:"raw"`
}

func DeserializeReceivers(raw []byte) ([]MarketReceiver, error) {
	// extract all the serializable receivers
	serializableReceivers := make([]serializableReceiverJSON, 0)
	if err := json.Unmarshal(raw, &serializableReceivers); err != nil {
		return nil, err
	}

	res := make([]MarketReceiver, 0, len(serializableReceivers))
	for _, receiver := range serializableReceivers {
		generator, ok := serializableMap[receiver.Key]
		if !ok {
			return nil, fmt.Errorf("unknown receiver type: %s", receiver.Key)
		}
		receiverInstance := generator()
		if err := json.Unmarshal(receiver.Raw, receiverInstance); err != nil {
			return nil, err
		}
		res = append(res, receiverInstance)
	}
	return res, nil
}

func SerializeReceivers(receivers ...MarketReceiver) ([]byte, error) {
	serializableReceivers := GetSerializableReceivers(receivers)
	res := make([]serializableReceiverJSON, len(serializableReceivers))
	for i, receiver := range serializableReceivers {
		raw, err := json.Marshal(receiver)
		if err != nil {
			return nil, err
		}
		res[i] = serializableReceiverJSON{receiver.PrefixKey(), raw}
	}

	return json.Marshal(res)
}
