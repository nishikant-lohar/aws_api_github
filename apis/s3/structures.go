package s3

type data struct {
	bucketID       string `json:"instance_id,omitempty"`
	InstanceType   string `json:"instance_type,omitempty"`
	InstanceStatus string `json:"instance_status,omitempty"`
}

type info struct {
	Success    bool   `json:"success,omitempty"`
	Message    string `json:"message,omitempty"`
	Status     string `json:"status,omitempty"`
	StatusCode int64  `json:"status_code,omitempty"`
}

// CreateInstanceResponse struct
type CreateInstanceResponse struct {
	Data data `json: data`
	Info info `json: info`
}

// CreateInstanceResponse struct
type DescribeInstanceResponse struct {
	Data data `json: data`
	Info info `json: info`
}

// StopInstanceResponse struct
type StopInstanceResponse struct {
	Info info `json: "info"`
}

// TerminateInstanceResponse struct
type TerminateInstanceResponse struct {
	Info info `json: "info"`
}

// StopInstanceInputParams struct
type StopInstanceInputParams struct {
	InstanceID string `json:"instance_id"`
}
