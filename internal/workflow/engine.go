package workflow

import (
	"fmt"
	"strings"
)

//service/manager/orchestrator
/*
the aim of the engine is to:
1- start the workflow:
	-takes the initial state, puts the stages into instance and sets the active stage
2- assess the transitions:
	- checks if the action is valid given the current state:
		-what's the current state
		-what's the action
		-is there a sutiable transition for it
		-is the stage status satisfied
3- applies the action:
	-continues the system after a step is completed
		-finds the step
		-updates its status
		-checks if the stage is completed
		-updates state if needed
		-logs into history
4- maintains the consistency of the workflow
	-prevents random data from being written everywhere.

For example:

a completed step shouldn't be completed again
a transition that doesn't exist shouldn't be executed
a step outside the current stage shouldn't be processed by mistake

*/

type Engine struct {
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) StartWorkflow(def WorkflowDefinition) (WorkflowInstance, error) {
	wfi := WorkflowInstance{}
	if len(def.Stages) < 1 {
		return wfi, fmt.Errorf("The stages of the workflow can not be empty!")
	}
	if strings.TrimSpace(def.Name) == "" {
		return wfi, fmt.Errorf("The name of the workflow can not be empty!")
	}
	if def.InitialState == "" {
		return wfi, fmt.Errorf("Initial state can not be empty!")
	}
	wfi.CurrentState = def.InitialState
	wfi.Stages = CloneStages(def.Stages)
	wfi.Stages[0].Status = StageStatusInProgress
	wfi.DefinitionID = def.ID
	wfi.CurrentStageIndex = 0
	wfi.History = make([]HistoryEntry, 0)
	wfi.Metadata = make(map[string]string)
	wfi.Name = def.Name
	return wfi, nil
}

func (e *Engine) FindTransition(def WorkflowDefinition, currentState State, action ActionType, currentStageStatus StageStatus) (Transition, bool, error) {
	tr := Transition{}
	if len(def.Transitions) < 1 {
		return tr, false, fmt.Errorf("Length of the transitions of the workflow with id %d can not be less than 1", def.ID)
	}
	if currentState == "" {
		return tr, false, fmt.Errorf("The current state of the workflow with id %d can not be empty", def.ID)
	}
	if action == "" {
		return tr, false, fmt.Errorf("The action type of the workflow with id %d can not be empty", def.ID)
	}
	// if currentStageStatus == "" {
	// 	return tr, false, fmt.Errorf("The current stage status of the workflow with id %d can not be empty", def.ID)
	// }

	for _, tr := range def.Transitions {
		if tr.FromState == currentState && tr.Action == action {
			if tr.RequiredStageStatus == "" {
				//ok
				return tr, true, nil
			}
			if tr.RequiredStageStatus == currentStageStatus {
				//ok
				return tr, true, nil
			}
		}
	}
	return tr, false, nil
}

func (e *Engine) ApplyAction(instance *WorkflowInstance, def WorkflowDefinition, stepID int, actorID string, action ActionType) error {
	// TODO
	if instance == nil {
		return fmt.Errorf("Workflow instance can not be nil")
	}
	if len(instance.Stages) < 1 {
		return fmt.Errorf("Length of the stages of the workflow with id %d can not be less than 1", instance.ID)
	}
	if len(instance.Stages) <= instance.CurrentStageIndex {
		return fmt.Errorf("Length of the stages of the workflow with id %d is %d and the CurrentStageIndex value id %d", instance.ID, len(instance.Stages), instance.CurrentStageIndex)
	}

	stage := &instance.Stages[instance.CurrentStageIndex]

	for i := 0; i < len(stage.Steps); i++ {
		if stepID != stage.Steps[i].ID {
			continue
		} //if does not match then keep searching

		//check if the action is valid for this step
		if action != stage.Steps[i].Action {
			return fmt.Errorf("Action for step with id %d (action: %s) is not same with the given action: %s", stepID, stage.Steps[i].Action, action)
		}
		//check if the actor ids are same
		if actorID != stage.Steps[i].ActorID {
			return fmt.Errorf("ActorId for step with id %d (ActorId: %s) is not same with the given actorId: %s", stepID, stage.Steps[i].ActorID, actorID)
		}
		//if everything is true then check if the status is suitable to update this step
		transition, isFound, err := e.FindTransition(def, instance.CurrentState, action, stage.Status)
		if err != nil {
			return err
		}
		if !isFound {
			return fmt.Errorf("Could not find a suitable transition for step with id %d for the action %s", stepID, action)
		}
		//update step
		// stage.Steps[i].Status = burası ne olacak

	}
	return nil

}

func (e *Engine) UpdateStageStatus(stage *Stage) {
	// TODO
}

func (e *Engine) AddHistoryEntry(instance *WorkflowInstance, stepID int, actorID string, action ActionType, fromState State, toState State) {
	// TODO
}
