package oop

import (
	"errors"
	"fmt"
)

/*

# 1. Система управления транспортом (ООП в Go)

Реализуйте иерархию классов транспорта, используя принципы ООП: наследование (композицию), инкапсуляцию и полиморфизм.

## Задание

### Цель
1. Создать базовый интерфейс `Vehicle` с методами:
    - `StartEngine() error` — запускает двигатель.
    - `StopEngine() error` — останавливает двигатель.
    - `GetInfo() string` — возвращает информацию о транспорте.
2. Реализовать три типа транспорта:
    - **Автомобиль** (`Car`):
        - Имеет поле `Brand` (марка) и `EngineOn` (состояние двигателя).
        - Метод `Honk() string` — возвращает "Beep beep!".
    - **Грузовик** (`Truck`):
        - Наследует функциональность `Car`.
        - Добавляет поле `CargoCapacity` (грузоподъемность в тоннах).
        - Переопределяет `Honk()` — возвращает "Honk Honk!".
    - **Электрокар** (`ElectricCar`):
        - Наследует функциональность `Car`.
        - Добавляет поле `BatteryLevel` (уровень заряда в %).
        - Переопределяет `StartEngine()`: запускается только если `BatteryLevel > 5%`.

### Требования
- Используйте **композицию** для наследования (встраивание структур).
- Поля `EngineOn`, `BatteryLevel` и `CargoCapacity` должны быть **инкапсулированы** (не экспортируемы).
- Для работы с полями добавьте методы:
    - `GetEngineStatus() bool` — возвращает состояние двигателя.
    - `GetBatteryLevel() int` — возвращает уровень заряда.
    - `GetCargoCapacity() float64` — возвращает грузоподъемность.
- Напишите unit-тесты, проверяющие:
    - Корректность запуска/остановки двигателя.
    - Полиморфизм через интерфейс `Vehicle`.
    - Уникальное поведение методов (например, `Honk()`).


*/

var (
	ErrEngineAlreadyRunning = errors.New("двигатель уже работает")
	ErrEngineOff            = errors.New("двигатель не запущен")
	ErrLowBattery           = errors.New("низкий заряд батареи")
	ErrInvalidBatteryLevel  = errors.New("неверный заряд батареи")
)

type Vehicle interface {
	StartEngine() error
	StopEngine() error
	GetInfo() string
}

type Car struct {
	Brand    string
	engineOn bool
}

func NewCar(brand string) *Car {
	return &Car{
		Brand: brand,
	}
}

func (c Car) Honk() string {
	return "Beep beep!"
}

func (c Car) GetEngineStatus() bool {
	return c.engineOn
}

func (c *Car) StartEngine() error {
	if c.engineOn {
		return ErrEngineAlreadyRunning
	}

	c.engineOn = true
	return nil
}

func (c *Car) StopEngine() error {
	if !c.engineOn {
		return ErrEngineOff
	}

	c.engineOn = false
	return nil
}

func (c *Car) GetInfo() string {
	return fmt.Sprintf("brand: %s, engine started: %v", c.Brand, c.GetEngineStatus())
}

type Truck struct {
	Car

	CargoCapacity float64
}

func NewTruck(brand string, cargoCapacity float64) *Truck {
	return &Truck{
		Car: *NewCar(brand),

		CargoCapacity: cargoCapacity,
	}
}

func (t Truck) GetCargoCapacity() float64 {
	return t.CargoCapacity
}

func (t Truck) Honk() string {
	return "Honk Honk!"
}

func (t *Truck) GetInfo() string {
	return fmt.Sprintf("cargo capacity: %f, %s", t.GetCargoCapacity(), t.Car.GetInfo())
}

type ElectricCar struct {
	Car

	BatteryLevel int
}

func NewElectricCar(brand string) *ElectricCar {
	return &ElectricCar{
		Car:          *NewCar(brand),
		BatteryLevel: 100,
	}
}

func (c ElectricCar) GetBatteryLevel() int {
	return c.BatteryLevel
}

func (c *ElectricCar) SetBatteryLevel(newLevel int) error {
	if newLevel < 0 || newLevel > 100 {
		return ErrInvalidBatteryLevel
	}

	c.BatteryLevel = newLevel

	return nil
}

func (c *ElectricCar) StartEngine() error {
	if c.GetBatteryLevel() < 5 {
		return ErrLowBattery
	}

	return c.Car.StartEngine()
}

func (c *ElectricCar) GetInfo() string {
	return fmt.Sprintf("battery level: %d, %s", c.GetBatteryLevel(), c.Car.GetInfo())
}
