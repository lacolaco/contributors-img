package contributors

import (
	"fmt"

	"contrib.rocks/libs/goutils/model"
)

func createContributorsJSONCacheKey(r *model.Repository) string {
	return fmt.Sprintf("contributors-json-cache/v1.2/%s--%s.json", r.Owner, r.RepoName)
}
