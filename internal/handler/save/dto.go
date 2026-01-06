package save

type (
	Request struct {
		Key   string `json:"key"`
		Value any    `json:"value"`
	}

	Response struct {
		Key   string `json:"key"`
		Value any    `json:"value"`
	}
)
