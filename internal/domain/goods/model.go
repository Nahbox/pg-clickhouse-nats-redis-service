package goods

type Project struct {
	Id        int
	Name      string
	CreatedAt string
}

type Good struct {
	Id          int    `json:"id"`
	ProjectId   int    `json:"projectId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	Removed     bool   `json:"removed"`
	CreatedAt   string `json:"createdAt"`
}

type Meta struct {
	Total   int `json:"total"`
	Removed int `json:"removed"`
	Limit   int `json:"limit"`
	Offset  int `json:"offset"`
}

type GetResponse struct {
	Meta  Meta   `json:"meta"`
	Goods []Good `json:"goods"`
}

type PatchData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PostData struct {
	Name string `json:"name"`
}

type DeleteResponse struct {
	Id         int  `json:"id"`
	CampaignId int  `json:"campaignId"`
	Removed    bool `json:"removed"`
}

type ReprioritizeResponse struct {
	Priorities []Priorities `json:"priorities"`
}

type Priorities struct {
	Id       int `json:"id"`
	Priority int `json:"priority"`
}
