package goods

type Repository interface {
	Add(data *Good) (*Good, error)
	Get(limit, offset int) (*GetResponse, error)
	Update(data *Good) (*Good, error)
	Delete(id, projectId int) (*DeleteResponse, *Good, error)
	UpdatePriority(id, projectId, newPriority int) (*ReprioritizeResponse, []Good, error)
}
