package devicecontrol

import (
	"errors"
	"fmt"
	"strings"
)

/*
## Задание

Реализуйте систему управления различными устройствами, используя интерфейсы и методы с указателями.

### Цель
1. Создать интерфейс `Device` с методами:
    - `UpdateOS(version string) error` — обновляет ОС устройства.
    - `GetInfo() string` — возвращает информацию об устройстве.
2. Реализовать интерфейс для трех устройств:
    - **Смартфон** (`Smartphone`)
    - **Ноутбук** (`Laptop`)
    - **Умные часы** (`Smartwatch`)

### Требования
- Каждое устройство должно иметь:
    - Поле `OSVersion string` (текущая версия ОС).
    - Поле `Model string` (модель устройства).
- Методы:
    - `UpdateOS`:
        - Обновляет `OSVersion`.
        - Возвращает ошибку `ErrUnsupported`, если обновление невозможно.
    - `GetInfo`:
        - Возвращает строку в формате: `"Модель: [модель], ОС: [версия]"`.
- **Специфичные правила**:
    - Смартфон нельзя обновить, если текущая версия ОС ≥ `"12.0"`.
    - Ноутбук поддерживает только версии ОС с префиксом `"Windows"`.
    - Умные часы нельзя обновить, если новая версия короче 5 символов.

*/

var (
	ErrUnsupported = errors.New("обновление недоступно")
)

type Device interface {
	UpdateOS(version string) error
	GetInfo() string
}

type OsInfo struct {
	OSVersion string
	Model     string
}

func (info OsInfo) String() string {
	format := "Модель: [%s], ОС: [%s]"
	return fmt.Sprintf(format, info.Model, info.OSVersion)
}

type device struct {
	info OsInfo
}

func (d device) GetInfo() string {
	return d.info.String()
}

type Smartphone struct {
	device
}

func NewSmartphone(info OsInfo) *Smartphone {
	return &Smartphone{
		device: device{
			info: info,
		},
	}
}

func (dev *Smartphone) UpdateOS(version string) error {
	if err := dev.canUpdate(); err != nil {
		return err
	}

	dev.info.OSVersion = version

	return nil
}

var (
	// curVersionValidateReg = regexp.MustCompile(`(1[2-9]+|[2-9]\d+|\d{3,})(\.\d+)?`)
	maxSmartphoneVersion = "12.0"
)

// Смартфон нельзя обновить, если текущая версия ОС ≥ `"12.0"`
func (dev *Smartphone) canUpdate() error {
	return dev.canUpdateCurrentVersion(maxSmartphoneVersion)
}

type Laptop struct {
	device
}

func NewLaptop(info OsInfo) *Laptop {
	return &Laptop{
		device: device{
			info: info,
		},
	}
}

func (dev *Laptop) UpdateOS(version string) error {
	if err := dev.validateUpdatingVersion(version); err != nil {
		return err
	}

	dev.info.OSVersion = version

	return nil
}

// Ноутбук поддерживает только версии ОС с префиксом `"Windows"
func (dev *Laptop) validateUpdatingVersion(version string) error {
	if !strings.HasPrefix(version, "Windows") {
		return ErrUnsupported
	}

	return nil
}

type Smartwatch struct {
	device
}

func NewSmartwatch(info OsInfo) *Smartwatch {
	return &Smartwatch{
		device: device{
			info: info,
		},
	}
}

func (dev *Smartwatch) UpdateOS(version string) error {
	if err := dev.validateUpdatingVersion(version); err != nil {
		return err
	}

	dev.info.OSVersion = version

	return nil
}

// Умные часы нельзя обновить, если новая версия короче 5 символов.
func (dev *Smartwatch) validateUpdatingVersion(version string) error {
	if len(version) < dev.getMinLenForUpdate() {
		return ErrUnsupported
	}

	return nil
}

func (dev *Smartwatch) getMinLenForUpdate() int {
	return 5
}
