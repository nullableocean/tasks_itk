package shape

import "math"

/*
## Задание 1: Реализация интерфейса `Shape`

### Описание
1. Создать интерфейс `Shape` с методами:
    - `Area() float64` — возвращает площадь фигуры.
    - `Perimeter() float64` — возвращает периметр (длину окружности для круга).
2. Реализовать интерфейс для двух фигур:
    - **Круг** (`Circle`), задаваемый радиусом.
    - **Прямоугольник** (`Rectangle`), задаваемый шириной и высотой.

### Требования
- Для `Circle`:
    - Площадь: `π * r²`
    - Периметр: `2 * π * r`
- Для `Rectangle`:
    - Площадь: `width * height`
    - Периметр: `2 * (width + height)`
- Используйте константу `math.Pi`.

*/

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct {
	r float64
}

func (c *Circle) Area() float64 {
	return math.Pi * math.Pow(c.r, 2)
}

func (c *Circle) Perimeter() float64 {
	return math.Pi * c.r * 2
}

type Rectangle struct {
	w, h float64
}

func (rec *Rectangle) Area() float64 {
	return rec.w * rec.h
}

func (rec *Rectangle) Perimeter() float64 {
	return (rec.w + rec.h) * 2
}
