package dynatrace

import "encoding/json"

type ingestResponse struct {
	LinesOk      int         `json:"linesOk"`
	LinesInvalid int         `json:"linesInvalid"`
	ErrorMessage interface{} `json:"error"`
}

func (e ingestResponse) Error() string {
	data, _ := json.Marshal(&e)
	return string(data)
}
