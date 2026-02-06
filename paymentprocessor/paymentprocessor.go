package paymentprocessor

import (
	"errors"

	rand "math/rand/v2"
)

/*
## Задание

### Цель
1. Создать интерфейс `PaymentProcessor` с методом `ProcessPayment(amount float64) error`.
2. Реализовать интерфейс для трех Банков:
    - **Sberbank**
    - **Tbank**
    - **Alfabank**

### Требования
- Каждый провайдер должен иметь уникальный идентификатор (например, `APIKey`).
- Метод `ProcessPayment` должен:
    - Возвращать `nil`, если сумма платежа положительная.
    - Возвращать ошибку `ErrInvalidAmount`, если сумма ≤ 0.
    - Возвращать ошибку `ErrProviderUnavailable`, если провайдер недоступен (заглушка). Сделать рандомный шанс, что банк недоступен.
*/

var (
	ErrInvalidAmount       = errors.New("некорректная сумма платежа")
	ErrProviderUnavailable = errors.New("провайдер недоступен")
)

type PaymentProcessor interface {
	ProcessPayment(amount float64) error
}

type ApiKey string

type paymentProvider struct {
	key ApiKey
}

func newProvider(key ApiKey) paymentProvider {
	return paymentProvider{
		key: key,
	}
}

func (p paymentProvider) GetApiKey() ApiKey {
	return p.key
}

func (p paymentProvider) validateAmount(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	return nil
}

func (p paymentProvider) checkProviderAvailable() error {
	n := rand.IntN(10)

	if n >= 7 {
		return ErrProviderUnavailable
	}

	return nil
}

type Sberbank struct {
	paymentProvider
}

func NewSberbank(key ApiKey) *Sberbank {
	return &Sberbank{
		paymentProvider: newProvider(key),
	}
}

func (p *Sberbank) ProcessPayment(amount float64) error {
	if err := p.validateAmount(amount); err != nil {
		return err
	}

	if err := p.checkProviderAvailable(); err != nil {
		return err
	}

	return nil
}

type Tbank struct {
	paymentProvider
}

func NewTbank(key ApiKey) *Tbank {
	return &Tbank{
		paymentProvider: newProvider(key),
	}
}

func (p *Tbank) ProcessPayment(amount float64) error {
	if err := p.validateAmount(amount); err != nil {
		return err
	}

	if err := p.checkProviderAvailable(); err != nil {
		return err
	}

	return nil
}

type Alfabank struct {
	paymentProvider
}

func NewAlfabank(key ApiKey) *Alfabank {
	return &Alfabank{
		paymentProvider: newProvider(key),
	}
}

func (p *Alfabank) ProcessPayment(amount float64) error {
	if err := p.validateAmount(amount); err != nil {
		return err
	}

	if err := p.checkProviderAvailable(); err != nil {
		return err
	}

	return nil
}
