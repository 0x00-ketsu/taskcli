package layout

// Validate task name
func validateTaskName(taskName string) bool {
	if len(taskName) < 3 {
		statusView.showForSeconds("[red]Task title should be at least 3 characters.", 5)
		return false
	}

	return true
}
