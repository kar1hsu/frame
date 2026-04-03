package tasks

import (
	"context"
	"encoding/json"
	"fmt"
)

type EmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func HandleEmailSend(ctx context.Context, payload []byte) error {
	var p EmailPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return fmt.Errorf("unmarshal email payload: %w", err)
	}

	// TODO: replace with real email sending logic
	fmt.Printf("[task] sending email to=%s subject=%s\n", p.To, p.Subject)
	return nil
}
