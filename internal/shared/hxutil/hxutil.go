package hxutil

import (
	"encoding/json"
	"maps"
)

// JoinTriggers is helper function when we send to client multiple events to trigger.
// It takes JSON string or just string and joins theme into JSON. If it recieves just
// string e.g. "example" it will put it in JSON as "example": "".
func JoinTriggers(triggers ...string) string {
	events := map[string]any{}

	for _, v := range triggers {
		vMap := map[string]any{}
		err := json.Unmarshal([]byte(v), &vMap)
		if err != nil {
			events[v] = ""
			continue
		}

		maps.Copy(events, vMap)
	}

	eventsStr, _ := json.Marshal(events)
	return string(eventsStr)
}
