package contributors

import "contrib.rocks/libs/go/model"

var _ error = &RepositoryNotFoundError{}

type RepositoryNotFoundError struct {
	Repository *model.Repository
}

func (e *RepositoryNotFoundError) Error() string {
	return "Repository not found: " + e.Repository.String()
}
