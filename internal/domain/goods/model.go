package goods

type Project struct {
	Id        int    `json:"-"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type Good struct {
	Id          int    `json:"-"`
	ProjectId   int    `json:"project_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	Removed     bool   `json:"removed"`
	CreatedAt   string `json:"created_at"`
}

type Meta struct {
	Total   int
	Removed int
	Limit   int
	Offset  int
}

type GetResponse struct {
	Meta  Meta
	Goods []Good
}
