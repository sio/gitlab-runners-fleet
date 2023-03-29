package gitlab

//
// Manage project-level runners (the only runner type available to free users at gitlab.com)
//

// Execute all maintenance chores:
//  - Purge dead runner entries
//  - Assign new runners to all projects
//func (api API) UpdateRunnerAssignments(tag string) error {
//	var (
//		err error
//		everything any
//		page string
//		projects = make(map[string]void)
//		runners = make(map[string]void)
//	)
//	everything, err = api.GraphQL(query, map[string]string{"cursor": page})
//}
//
//func removeRunner(uid string) {
//}
//
//func assignRunner(runner_uid string, project_uids []string) {
//}

type void struct{}
