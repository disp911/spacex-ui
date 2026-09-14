package job

import (
	"github.com/disp911/spacex-ui/v2/web/service"
)

// CheckTelemtRunningJob restarts the Telegram proxy if it exits unexpectedly.
type CheckTelemtRunningJob struct {
	telemtService service.TelemtService
}

// NewCheckTelemtRunningJob creates a new Telegram proxy health check job.
func NewCheckTelemtRunningJob() *CheckTelemtRunningJob {
	return new(CheckTelemtRunningJob)
}

// Run restarts telemt when it should be running but is not.
func (j *CheckTelemtRunningJob) Run() {
	j.telemtService.RestartIfCrashed()
}
