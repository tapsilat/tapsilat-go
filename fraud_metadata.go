package tapsilat

import (
	"maps"
	"slices"
)

// FraudMetadataPrefix marks order metadata that is sent to fraud screening: "fraud.check_in_date" reaches the
// fraud rules as the check_in_date attribute. Other metadata is not sent.
const FraudMetadataPrefix = "fraud."

// FraudMetadata turns attributes into order metadata entries with the fraud prefix, in key order.
// Names must match ^[a-z][a-z0-9_]{0,63}$; an order carries at most 50 of them, each value at most 512 characters.
func FraudMetadata(attributes map[string]string) []OrderMetadata {
	metadata := make([]OrderMetadata, 0, len(attributes))
	for _, name := range slices.Sorted(maps.Keys(attributes)) {
		metadata = append(metadata, OrderMetadata{Key: FraudMetadataPrefix + name, Value: attributes[name]})
	}
	return metadata
}
