package postgres

type Projects struct {
	Id        int    `json:"-"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type Goods struct {
	Id          int    `json:"-"`
	ProjectId   int    `json:"project_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	Removed     bool   `json:"removed"`
	CreatedAt   string `json:"created_at"`
}
