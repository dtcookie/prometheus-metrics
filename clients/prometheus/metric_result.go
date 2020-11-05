package prometheus

import (
	"encoding/json"
	"strconv"
)

// MetricResult has no documentation
type MetricResult struct {
	Status string           `json:"status"`
	Data   MetricResultData `json:"data"`
}

// MetricResultData has no documentation
type MetricResultData struct {
	ResultType string         `json:"resultType"`
	Result     []MetricRecord `json:"result"`
}

// UnmarshalJSON has no documentation
func (mrd *MetricResultData) UnmarshalJSON(data []byte) error {
	var err error
	rawMessages := map[string]json.RawMessage{}
	if err = json.Unmarshal(data, &rawMessages); err != nil {
		return err
	}
	if err = json.Unmarshal(rawMessages["resultType"], &mrd.ResultType); err != nil {
		return err
	}
	if mrd.ResultType != "vector" {
		return nil
	}
	return json.Unmarshal(rawMessages["result"], &mrd.Result)
}

// MetricRecord has no documentation
type MetricRecord struct {
	Dimensions map[string]string `json:"metric"`
	Value      Measurement       `json:"value"`
}

// Measurement has no documentation
type Measurement struct {
	Timestamp uint64  `json:"timestamp"`
	Value     float64 `json:"value"`
}

// UnmarshalJSON has no documentation
func (m *Measurement) UnmarshalJSON(data []byte) error {
	var err error
	rawMessages := []json.RawMessage{}
	if err = json.Unmarshal(data, &rawMessages); err != nil {
		return err
	}
	var timeStamp float64
	if err = json.Unmarshal(rawMessages[0], &timeStamp); err != nil {
		return err
	}
	m.Timestamp = uint64(timeStamp)
	var strValue string
	if err = json.Unmarshal(rawMessages[1], &strValue); err != nil {
		return err
	}
	if m.Value, err = strconv.ParseFloat(strValue, 64); err != nil {
		return err
	}
	return nil
}

/*
{
	"status": "success",
	"data": {
		"resultType": "vector",
		"result": [
			{
				"metric": {
					"__name__": "node_zfs_abd_struct_size",
					"instance": "localhost:9100",
					"job": "node"
				},
				"value": [
					1604401209.823,
					"0"
				]
			}
		]
	}
}
*/
