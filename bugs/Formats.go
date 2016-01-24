package bugs
import (
    "encoding/json"
)


func (b Bug) ToJSONString() (string, error) {
	bJSONStruct := struct {
		Identifier  string `json:",omitempty"`
		Title       string
		Description string
		Status      string   `json:",omitempty"`
		Priority    string   `json:",omitempty"`
		Milestone   string   `json:",omitempty"`
		Tags        []string `json:",omitempty"`
	}{
		Identifier:  b.Identifier(),
		Title:       b.Title(""),
		Description: b.Description(),
		Status:      b.Status(),
		Priority:    b.Priority(),
		Milestone:   b.Milestone(),
		Tags:        b.StringTags(),
	}

	bJSON, err := json.Marshal(bJSONStruct)
    if err != nil {
        return "", err
    }
	return string(bJSON), nil
}
