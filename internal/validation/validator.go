package validation

type Validator interface {
	Validate() error
}

type StringValidator interface {
	Validator
	GetValue() string
}

type NumberValidator interface {
	Validator
	GetValue() int
}

type nameValidator struct {
	value string
}

func NewNameValidator(value string) StringValidator {
	return &nameValidator{value: value}
}

func (n *nameValidator) Validate() error {
	if len(n.value) > GetConfig().MaxNameLength {
		return ErrNameTooLong
	}
	return nil
}

func (n *nameValidator) GetValue() string {
	return n.value
}

type ageValidator struct {
	value int
}

func NewAgeValidator(value int) NumberValidator {
	return &ageValidator{value: value}
}

func (a *ageValidator) Validate() error {
	cfg := GetConfig()
	if a.value < cfg.MinAge || a.value > cfg.MaxAge {
		return ErrInvalidAge
	}
	return nil
}

func (a *ageValidator) GetValue() int {
	return a.value
}
