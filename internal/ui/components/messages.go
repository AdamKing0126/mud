package components

import "github.com/google/uuid"

type SelectedMessage struct {
	Selected []map[string]string
}

type SubmitMessage struct {
  SenderId uuid.UUID
  RecipientId uuid.UUID
	Data any
}

