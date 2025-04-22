package contributors

import (
	"fmt"

	"contrib.rocks/apps/api/go/model"
)

func createContributorsJSONCacheKey(r *model.Repository) string {
	return fmt.Sprintf("contributors-json-cache/v1.2/%s--%s.json", r.Owner, r.RepoName)
}
