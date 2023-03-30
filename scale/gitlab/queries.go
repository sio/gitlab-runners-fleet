package gitlab

import _ "embed"

//go:embed queries/runners_by_project.gql
var queryRunnersByProject string
