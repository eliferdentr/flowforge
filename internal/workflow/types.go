package workflow

import "time"

type CompletionMode int

const (
	CompletionModeAll CompletionMode = iota
	CompletionModeAny
	CompletionModeSequential
)

type StepStatus string

const (
	StepStatusPending   StepStatus = "PENDING"
	StepStatusCompleted StepStatus = "COMPLETED"
	StepStatusRejected  StepStatus = "REJECTED"
	StepStatusSkipped   StepStatus = "SKIPPED"
)

type StageStatus string

const (
	StageStatusPending    StageStatus = "PENDING"
	StageStatusInProgress StageStatus = "IN_PROGRESS"
	StageStatusCompleted  StageStatus = "COMPLETED"
	StageStatusSkipped    StageStatus = "SKIPPED"
)

type ActionType string
type State string

// smallest unit of work in a workflow
type Step struct {
	ID       int
	Name     string
	Action   ActionType
	ActorID  string
	Status   StepStatus
	Order    int
	Metadata map[string]string
}

// groups one or more steps in the same phase of a workflow
type Stage struct {
	ID             int
	Name           string
	Steps          []Step
	CompletionMode CompletionMode
	Metadata       map[string]string
	Status         StageStatus
}

/*
DRAFT + SUBMIT -> WAITING_MANAGER
WAITING_MANAGER + APPROVE -> WAITING_FINANCE
WAITING_MANAGER + REJECT -> REJECTED
*/
type Transition struct {
	FromState           State
	ToState             State
	Action              ActionType
	RequiredStageStatus StageStatus
}

//represents the template/definition of a workflow
type WorkflowDefinition struct {
	ID           int
	Name         string
	InitialState State
	States       []State
	Actions      []ActionType
	Stages       []Stage
	Transitions  []Transition
}

//represents one running workflow process
type WorkflowInstance struct {
	ID                int
	DefinitionID      int
	Name              string
	Stages            []Stage
	CurrentState      State
	CurrentStageIndex int
	History           []HistoryEntry
	Metadata          map[string]string
}

// keeps the log of what happened in the workflow
type HistoryEntry struct {
	StepID    int
	ActorID   string
	Action    ActionType
	FromState State
	ToState   State
	Timestamp time.Time
	Metadata  map[string]string
}
