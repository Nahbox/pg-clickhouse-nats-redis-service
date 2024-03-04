package clickhouse

type Logs struct {
	Id          int    `json:"-"`
	ProjectId   int    `json:"ProjectId"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Priority    int    `json:"Priority"`
	Removed     bool   `json:"Removed"`
	EventTime   string `json:"EventTime"`
}
