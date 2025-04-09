package orchestrator

import (
	infra_component "github.com/Michaelpalacce/go-btva/internal/run/components/infra"
	software_component "github.com/Michaelpalacce/go-btva/internal/run/components/software"
	"github.com/Michaelpalacce/go-btva/internal/run/state"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

type TaskFunc func() error

// Orchestrator is a struct that orchestrates the tasks collections and order of execution of them
type Orchestrator struct {
	OS    *os.OS
	State *state.State

	SoftwareTasks []TaskFunc
	InfraTasks    []TaskFunc
	EnvTasks      []TaskFunc
	FinalTasks    []TaskFunc

	components struct {
		infraComponent    infra_component.InfraComponent
		softwareComponent software_component.SoftwareComponent
	}
}

// NewOrchestrator will return a new Orchestrator that is used to contain and execute tasks
func NewOrchestrator(os *os.OS, state *state.State) *Orchestrator {
	orchestrator := &Orchestrator{OS: os, State: state}

	orchestrator.components.infraComponent = *infra_component.NewInfraComponent(os, state)
	orchestrator.components.softwareComponent = *software_component.NewSoftwareComponent(os, state)

	return orchestrator
}

// runTaskOption accepts an orchestrator and is supposed to modify the state and add tasks to it
type runTaskOption func(o *Orchestrator) error

// Reset will reset all the task collections to empty
func (o *Orchestrator) Reset() {
	o.SoftwareTasks = make([]TaskFunc, 0)
	o.InfraTasks = make([]TaskFunc, 0)
	o.EnvTasks = make([]TaskFunc, 0)
	o.FinalTasks = make([]TaskFunc, 0)
}

// Tasks can be used to modify the current execution stack of the orhecstrator
func (o *Orchestrator) Tasks(options ...runTaskOption) error {
	for _, option := range options {
		if err := option(o); err != nil {
			return err
		}
	}

	return nil
}

// Run executes all the tasks added to the orchestrator in order specified within
// Will remove all tasks at the end if successful. Otherwise, tasks will be kept so they can be re-executed
func (o *Orchestrator) Run(options ...runTaskOption) error {
	if err := o.Tasks(options...); err != nil {
		return err
	}

	allTasks := make([][]TaskFunc, 0)
	allTasks = append(allTasks, o.SoftwareTasks, o.InfraTasks, o.EnvTasks, o.FinalTasks)

	for _, taskCollection := range allTasks {
		for _, task := range taskCollection {
			if err := task(); err != nil {
				return err
			}
		}
	}

	o.Reset()

	return nil
}
