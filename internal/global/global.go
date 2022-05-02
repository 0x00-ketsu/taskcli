package global

import (
	"github.com/0x00-ketsu/taskcli/internal/repository"
	"github.com/rivo/tview"
)

// Maintain variables which will be quoted in different packages
var (
	App    *tview.Application
	Layout *tview.Flex
	TaskRepo repository.Task
)
