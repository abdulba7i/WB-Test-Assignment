package model

import (
	"fmt"
	"l0/internal/validation"
)

// Actor представляет основную модель актера
type Actor struct {
	name validation.StringValidator
	age  validation.NumberValidator
}

// NewActor создает нового актера с валидацией
func NewActor(name string, age int) (*Actor, error) {
	actor := &Actor{
		name: validation.NewNameValidator(name),
		age:  validation.NewAgeValidator(age),
	}

	if err := actor.Validate(); err != nil {
		return nil, fmt.Errorf("invalid actor: %w", err)
	}

	return actor, nil
}

// Validate выполняет валидацию всех полей актера
func (a *Actor) Validate() error {
	if err := a.name.Validate(); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}
	if err := a.age.Validate(); err != nil {
		return fmt.Errorf("invalid age: %w", err)
	}
	return nil
}

// Name возвращает имя актера
func (a *Actor) Name() string {
	return a.name.GetValue()
}

// Age возвращает возраст актера
func (a *Actor) Age() int {
	return a.age.GetValue()
}
