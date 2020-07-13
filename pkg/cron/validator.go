package cron

import "context"

type Validator struct{}

func NewValidator() *Validator {
return &Validator{}
}

func (v *Validator) Validate(ctx context.Context, interval string) error {
_, err := parseIntervalToCron(interval)
return err
}