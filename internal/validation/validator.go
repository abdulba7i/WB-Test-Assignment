package validation

// Validator определяет интерфейс для валидации
type Validator interface {
	Validate() error
}

// StringValidator определяет интерфейс для валидации строк
type StringValidator interface {
	Validator
	GetValue() string
}

// NumberValidator определяет интерфейс для валидации чисел
type NumberValidator interface {
	Validator
	GetValue() int
}

// nameValidator реализует валидацию для имени
type nameValidator struct {
	value string
}

// NewNameValidator создает новый валидатор имени
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

// ageValidator реализует валидацию для возраста
type ageValidator struct {
	value int
}

// NewAgeValidator создает новый валидатор возраста
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
