package dto

type WSFrame struct {
	Event string         `json:"event"`
	ID    string         `json:"id,omitempty"`
	Data  map[string]any `json:"data,omitempty"`
}
