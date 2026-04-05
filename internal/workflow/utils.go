package workflow

func CloneMetadata(metadata map[string]string) map[string]string {
	if metadata == nil {
		return nil
	}

	cloned := make(map[string]string, len(metadata))
	for k, v := range metadata {
		cloned[k] = v
	}

	return cloned
}

func CloneSteps(steps []Step) []Step {
	if steps == nil {
		return nil
	}

	cloned := make([]Step, len(steps))
	for i, step := range steps {
		cloned[i] = Step{
			ID:       step.ID,
			Name:     step.Name,
			Action:   step.Action,
			ActorID:  step.ActorID,
			Status:   step.Status,
			Order:    step.Order,
			Metadata: CloneMetadata(step.Metadata),
		}
	}

	return cloned
}

func CloneStages(stages []Stage) []Stage {
	if stages == nil {
		return nil
	}

	cloned := make([]Stage, len(stages))
	for i, stage := range stages {
		cloned[i] = Stage{
			ID:             stage.ID,
			Name:           stage.Name,
			Steps:          CloneSteps(stage.Steps),
			CompletionMode: stage.CompletionMode,
			Metadata:       CloneMetadata(stage.Metadata),
			Status:         stage.Status,
		}
	}

	return cloned
}
