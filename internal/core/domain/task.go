package domain

import (
	"fmt"
	"time"

	core_errors "github.com/DenisAstrakhan/api-server/internal/core/errors"
)

type Task struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAT    time.Time
	CompletedAt  *time.Time
	AuthorUserID int
}

func NewTask(
	id int,
	version int,
	title string,
	description *string,
	completed bool,
	createdAT time.Time,
	completedAt *time.Time,
	authorUserID int,
) Task {
	return Task{
		ID:           id,
		Version:      version,
		Title:        title,
		Description:  description,
		Completed:    completed,
		CreatedAT:    createdAT,
		CompletedAt:  completedAt,
		AuthorUserID: authorUserID,
	}
}

func NewTaskUninitialized(
	title string,
	description *string,
	authorUserID int,
) Task {
	return NewTask(
		UninitializedID,
		UninitializedVersion,
		title,
		description,
		false,
		time.Now(),
		nil,
		authorUserID,
	)
}

func (t *Task) Validate() error {
	titleLen := len([]rune(t.Title))
	if titleLen < 1 || titleLen > 100 {
		return fmt.Errorf("invalid `Title` len: %d: %w", titleLen, core_errors.ErrInvalidArgument)
	}

	if t.Description != nil {
		descriptionLen := len([]rune(*t.Description))
		if descriptionLen < 1 || descriptionLen > 1000 {
			return fmt.Errorf("invalid `Description` len: %d: %w", titleLen, core_errors.ErrInvalidArgument)
		}
	}

	if t.Completed {
		if t.CompletedAt == nil {
			return fmt.Errorf("'CompletedAt' can't be 'nil if 'Completed'=='true': %w", core_errors.ErrInvalidArgument)
		}
		if t.CompletedAt.Before(t.CreatedAT) {
			return fmt.Errorf("'CompletedAt' can't be before 'CreatedAT': %w", core_errors.ErrInvalidArgument)
		}
	} else {
		if t.CompletedAt != nil {
			return fmt.Errorf("'CompletedAt' must be 'nil if 'Completed'=='false': %w", core_errors.ErrInvalidArgument)
		}
	}
	return nil
}

func (t *Task) CompletionDuration() *time.Duration {
	if !t.Completed {
		return nil
	}

	if t.CompletedAt == nil {
		return nil
	}
	// вычитаем из CompletedAt CreatedAT
	duration := t.CompletedAt.Sub(t.CreatedAT)
	return &duration
}

type TaskPatch struct {
	Title       Nullable[string] `json:"title"`
	Description Nullable[string] `json:"description"`
	Completed   Nullable[bool]   `json:"completed"`
}

func NewTaskPatch(title Nullable[string], description Nullable[string], completed Nullable[bool]) TaskPatch {
	return TaskPatch{
		Title:       title,
		Description: description,
		Completed:   completed,
	}
}

func (p *TaskPatch) Validate() error {
	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf("'Title' can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}
	if p.Completed.Set && p.Completed.Value == nil {
		return fmt.Errorf("'Completed' can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}

func (t *Task) ApplyPatch(p TaskPatch) error {
	if err := p.Validate(); err != nil {
		return fmt.Errorf("validate task patch: %w", err)
	}
	tmpTask := *t
	if p.Title.Set {
		tmpTask.Title = *p.Title.Value
	}
	if p.Description.Set {
		tmpTask.Description = p.Description.Value
	}
	if p.Completed.Set {
		tmpTask.Completed = *p.Completed.Value
		if tmpTask.Completed {
			completedAt := time.Now()
			tmpTask.CompletedAt = &completedAt
		} else {
			tmpTask.CompletedAt = nil
		}
	}

	if err := tmpTask.Validate(); err != nil {
		return fmt.Errorf("validate patched task: %w", err)
	}
	*t = tmpTask
	return nil
}
