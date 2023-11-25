package market

import (
	"encoding/json"
	"fmt"
)

type ReceiverGenerator func() MarketReceiver

var (
	serializableMap = map[string]ReceiverGenerator{}
)

func RegisterSerializableReceiver(key string, generator ReceiverGenerator) {
	serializableMap[key] = generator
}

func GetSerializableReceivers(receivers []MarketReceiver) []SerializableReceiver {
	res := make([]SerializableReceiver, 0, len(receivers))
	for _, receiver := range receivers {
		if serializableReceiver, ok := receiver.(SerializableReceiver); ok {
			res = append(res, serializableReceiver)
		}
	}
	return res
}

type serializableReceiverJSON struct {
	Key string `json:"key"`
	Raw []byte `json:"raw,omitempty"`
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
		if _, ok := receiverInstance.(SerializableDataReceiver); ok {
			fmt.Printf("trying to unmarshal %s\n", receiver.Key)
			if err := json.Unmarshal(receiver.Raw, receiverInstance); err != nil {
				return nil, fmt.Errorf("failed to deserialize receiver %s: %s", receiver.Key, err)
			}
		}
		res = append(res, receiverInstance)
	}
	return res, nil
}

func SerializeReceivers(receivers ...MarketReceiver) ([]byte, error) {
	serializableReceivers := GetSerializableReceivers(receivers)
	res := make([]serializableReceiverJSON, len(serializableReceivers))
	var err error
	for i, receiver := range serializableReceivers {
		var raw []byte = nil
		if serializableDataReceiver, ok := receiver.(SerializableDataReceiver); ok {
			raw, err = json.Marshal(serializableDataReceiver)
			if err != nil {
				return nil, err
			}
		}
		res[i] = serializableReceiverJSON{receiver.PrefixKey(), raw}
	}

	return json.Marshal(res)
}
