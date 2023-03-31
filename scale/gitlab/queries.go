package gitlab

import _ "embed"

//go:embed queries/runners_by_project.gql
var queryRunnersByProject string

//go:embed queries/runner_assign.gql
var mutationRunnerAssignment string

//go:embed queries/runner_remove.gql
var mutationRunnerRemove string
