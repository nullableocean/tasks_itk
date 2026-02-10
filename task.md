DEFER
# Задание: Анализ кода на Go  

### 1

**Ваша задача:** Определить вывод программы и зафиксировать ответы **в сообщениях коммитов** с кратким объяснением логики.
```
package main  
  
import (  
    "fmt"  
)  
  
func main() {  
    fmt.Println("start")  
    for i := 1; i < 4; i++ {  
       defer fmt.Println(i)  // start end 3 2 1, defer выполняет по принципу lifo 
    }  
    fmt.Println("end")  
}
```

# Задание: Анализ кода на Go  

**Ваша задача:** Определить вывод программы и зафиксировать ответы **в сообщениях коммитов** с кратким объяснением логики.  
Исправьте код так, чтобы defer вывел обновленное значение value.
```
package main  
  
import "fmt"  
  
func main() { 
    value := 123  
     // defer fmt.Println(value)
    defer func(v *int){ fmt.Println(*v) }(&value) 
    changeValue(&value)  
}  
func changeValue(value *int) {  
    *value = 456  
}
```
3.
**Ваша задача:** Определить вывод программы и зафиксировать ответы **в сообщениях коммитов** с кратким объяснением логики.
```
package main

import (
	"errors"
	"fmt"
)

func main() {
    println("Case 1")
    case1()
    println()
    println()

    println("Case 2")
    case2()
    println()
    println()

    println("Case 3")
    case3()
    println()
    println()

}

func case1() {
    helperWithDefer := func(isError bool) error {
        var retVal error

        defer func() {
            retVal = errors.New("Extra error") // возвращаемое значение не изменится
        }()

        if isError {
            retVal = errors.New("Default error")
        }

        return retVal // nil || Default error
    }

    helperWithoutDefer := func(isError bool) error {
        var retVal error

        if isError {
            retVal = errors.New("Default error")
        }

        return retVal // nil || Default error
    }

    fmt.Println("\twithout:") 
    fmt.Println(helperWithoutDefer(false))
    fmt.Println(helperWithoutDefer(true))
    fmt.Println("\twith:")
    fmt.Println(helperWithDefer(false))
    fmt.Println(helperWithDefer(true))
	/*
		whitout
		nil
		Default error
		with
		nil
		Default error
	   */
}

func case2() {
    helperWithDefer := func(isError bool) (retVal error) {
        defer func() {
            retVal = errors.New("Extra error") // возвращаемое значение изменяется, тк именованный возврат
        }()

        if isError {
            retVal = errors.New("Default error")
        }

        return
    }

    helperWithoutDefer := func(isError bool) (retVal error) {
        if isError {
            retVal = errors.New("Default error")
        }

        return
    }

    fmt.Println("\twithout:")
    fmt.Println(helperWithoutDefer(false))
    fmt.Println(helperWithoutDefer(true))
    fmt.Println("\twith:")
    fmt.Println(helperWithDefer(false))
    fmt.Println(helperWithDefer(true))
	/*
		whitout
		nil
		Default error
		with
		Extra error
		Extra error
	   */
}

func case3() {
    helperWithDefer := func(isError bool) (retVal error) {
        defer func() {
            retVal = errors.New("First Error")
        }()

        defer func() {
            retVal = errors.New("Second Error")
        }()

        if isError {
            retVal = errors.New("Default error")
        }

        return
    }

    helperWithoutDefer := func(isError bool) (retVal error) {
        if isError {
            retVal = errors.New("Default error")
        }

        return
    }

    fmt.Println("\twithout:")
    fmt.Println(helperWithoutDefer(false))
    fmt.Println(helperWithoutDefer(true))
    fmt.Println("\twith:")
    fmt.Println(helperWithDefer(false))
    fmt.Println(helperWithDefer(true))
    
    /*
		whitout
		nil
		Default error
		with
		First Error
		First Error
	   */
}
```
-------------------------------------------------------------




ERRORS

# Задание: Возврат ошибки без дополнительных пакетов  

### 1

Реализуйте функцию `handle() error`, которая возвращает ошибку, **не используя дополнительные пакеты** (кроме уже импортированного `fmt`).  
  
## Требования  
1. Функция `handle()` должна возвращать тип `error`.  
2. Запрещено подключать пакеты, кроме `fmt`.  
3. Ошибка должна содержать текст (не `nil`).


```
package main

import "fmt"

func handle() error {
   return fmt.Errorf("error with fmt")
}

func main() {
    err := handle()
    if err != nil {
        fmt.Println(handle())
    }
}
```

# Кастомные ошибки в Go  

### 2

Реализуйте три типа кастомных ошибок:  
1. Простую ошибку через `errors.New()`.  
2. Ошибку с форматированием через `fmt.Errorf()`.  
3. Структуру `MyError`, реализующую интерфейс `error`.  

  
## Требования  
1. **Простая ошибка**:  
   - Создайте функцию `SimpleError() error`, возвращающую ошибку с текстом "простая ошибка".  

1. **Форматированная ошибка**:  
   - Создайте функцию `FormattedError(age int) error`, которая возвращает ошибку в формате: "ошибка: возраст %d недопустим".  
   - Добавьте оборачивание ошибки с `%w`.  

1. **Структура `MyError`**:  
   - Реализуйте метод `Error() string`.  
   - Добавьте поле `Code int` для кода ошибки.  
   - Создайте функцию `StructError() error`, возвращающую `MyError{Code: 404, Msg: "не найдено"}`.
   - 
  ```go
    package main

    import (
        "errors"
        "fmt"
    )

    var ErrTimeStarted = errors.New("time started")
    var ErrTimeStoped = fmt.Errorf("time stoped")
    var ErrUnacceptableAge = errors.New("некорректный возраст")

    func SimpleError() error {
        return errors.New("простая ошибка")
    }

    func FormattedError(age int) error {
        return fmt.Errorf("%w: возраст %d недопустим", ErrUnacceptableAge, age)
    }

    type MyError struct {
        Msg string
        Code int
    }

    func (e MyError) Error() string {
        return e.Msg
    }

    func StructError() error {
        return MyError{
            Code: 404,
            Msg: "не найдено",
        }
    }
```

# Анализ цепочек ошибок в Go error.Is

### 3

Напишите функцию `ProcessError(err error)`, которая:
1. Проверяет, содержит ли ошибка в своей цепочке ошибку типа `TimeoutError` (кастомный тип).
2. Если содержит, выводит "Требуется повторная попытка".
3. Проверяет, содержит ли ошибка в своей цепочке ошибку `ErrNotFound` (стандартная ошибка).
4. Если содержит, выводит "Ресурс не найден".
5. Для остальных ошибок выводит "Неизвестная ошибка".



## Требования
1. Определите кастомные ошибки:

```
   var (
       ErrNotFound   = errors.New("ресурс не найден")
       TimeoutError = errors.New("таймаут операции")
   )
```
2. Создайте функцию SimulateRequest() error, которая:
   - В 50% случаев возвращает TimeoutError, обёрнутую в fmt.Errorf("запрос не выполнен: %w", TimeoutError).
   - В 30% случаев возвращает ErrNotFound, обёрнутую в fmt.Errorf("ошибка: %w", ErrNotFound).
   - В 20% случаев возвращает новую ошибку "неизвестная ошибка".
   - Реализуйте логику анализа ошибок в ProcessError.
3. Реализуйте логику анализа ошибок в ProcessError.


```go
    var (
    ErrNotFound   = errors.New("ресурс не найден")
    TimeoutError = errors.New("таймаут операции")
    )

    func SimulateRequest() error {
        n := rand.Intn(100)
        
        if n > 50 {
            return fmt.Errorf("запрос не выполнен: %w", TimeoutError)
        }
        
        if n > 20 {
            return fmt.Errorf("ошибка: %w", ErrNotFound)
        }
        
        return fmt.Errorf("неизвестная ошибка")
    }

    func ProcessError(err error) {
        if errors.Is(err, TimeoutError) {
            fmt.Println("Требуется повторная попытка")
        } else if errors.Is(err, ErrNotFound) {
            fmt.Println("Ресурс не найден")
        } else {
            fmt.Println("Неизвестная ошибка")
        }
    }
```
-------------------------------------------------------


GENERICS
# Универсальный стек (LIFO) на Go с дженериками  
  
Реализация простого обобщенного стека с методами `Push`, `Pop`, `Peek` и `IsEmpty`. Используются дженерики Go (1.18+).  
  
## Требования  
  
1. **Структура `Stack[T]`**:  
   - Поле `elements` для хранения элементов (слайс `[]T`).  
   - Методы:  
        - `NewStack[T]() *Stack[T]`: конструктор.  
        - `Push(value T)`: добавление элемента в стек.  
        - `Pop() (T, bool)`: удаление и возврат верхнего элемента (с проверкой на пустоту).  
        - `Peek() (T, bool)`: возврат верхнего элемента без удаления.  
        - `IsEmpty() bool`: проверка стека на пустоту.  
  
2. **Дополнительно**:  
   - Гарантировать безопасность операций (например, `Pop` на пустом стеке возвращает `false`).  
   - Использовать слайс для эффективного добавления/удаления элементов.


```
Решение в ./stack 
```

----------------------------------------------------
Interface
1.
# Универсальный потокобезопасный кэш с TTL, очисткой и JSON-сериализацией  
  
Реализация in-memory кэша на Go с расширенными возможностями: автоматическое удаление ключей, очистка и сериализация данных.  
  
---  
  
## **Основные возможности**  
  
**TTL (Time-To-Live)**    
- Автоматическое удаление ключей по истечении времени жизни.   
  
**Очистка кэша**    
- Мгновенное удаление всех данных одной командой.    
  
**Сериализация в JSON**    
- Преобразование актуальных данных в JSON-формат.   
  
**Потокобезопасность**    
- Использование `sync.RWMutex` для конкурентного доступа.   
  
**Универсальное хранение**    
- Поддержка любых типов данных через `interface{}`.  
---  
  
## **Методы**  
### **Базовые операции**  
  
 - `Set(key string, value interface{}, ttl time.Duration)` Добавляет значение с указанным TTL   
 - `Get(key string) (interface{}, bool)`  Возвращает значение (с проверкой TTL)   
 - `Delete(key string)`  Удаляет конкретный ключ  
 - `Exists(key string) bool` Проверяет наличие непросроченного ключа   
  
### **Расширенные функции**  
 - `Clear()` Полностью очищает кэш   
 - `ToJSON() ([]byte, error)` Сериализует данные в JSON   
 - `GetAs[T any](key string) (T, error)`  Типизированное получение  
  

```
Решение в ./cachettl
```
 
--------------------
# Задание: Анализ кода на Go

Это задание направлено на глубокое понимание работы срезов (interface), их модификации и передачи в функциях Go.  
**Ваша задача:** Определить вывод каждой из предложенных программ и зафиксировать ответы **в сообщениях коммитов** с кратким объяснением логики.
1.
```
package main  
  
import (  
    "fmt"  
)  
  
type MyError struct {  
    data string  
}  
  
func (m *MyError) Error() string {  
    return m.data  
}  
func foo(i int) error {  
    var err *MyError  
    if i > 5 {  
       err = &MyError{data: "i>5"}  
    }  
    return err  
}  
func main() {  
    err := foo(4)  
    if err != nil {  
       fmt.Println("oops")  // oops, структура интерфейса error будет не nil, так как хранит данные о типе (*MyError)
    } else {  
       fmt.Println("ok")  
    }  
}
```

2.
```
package main  
  
import (  
    "fmt"  
)  
  
type errorString struct {  
    s string  
}  
  
func (e errorString) Error() string {  
    return e.s  
}  
  
func checkErr(err error) {  
    fmt.Println(err == nil)  
}  
  
func main() {  
    var e1 error  
    checkErr(e1)  // true
  
    var e *errorString  
    checkErr(e)  // false
  
    e = &errorString{}  
    checkErr(e)  // false
  
    e = nil  
    checkErr(e)  // false, в функции получаем интерфейс который будет знать о типе у e
}
```

3.
```
package main

import "fmt"

type CustomError struct {
	message string
}

func (e *CustomError) Error() string {
	return e.message
}

func returnError(flag bool) error {
	if flag {
		return &CustomError{"Что-то пошло не так"}
	}
	var err *CustomError
	return err
}

func main() {
	err1 := returnError(true)
	err2 := returnError(false)

	fmt.Println("err1 == nil:", err1 == nil) // false
	fmt.Println("err2 == nil:", err2 == nil) // false

}
```

--------------------------------------------------------
PANICS
# Паники и их обработка в Go

### 1.

В этом задании вам нужно работать с `panic` и `recover`, чтобы понять их механику.

## Требования

1. **Функция с паникой**
    - Реализуйте функцию `CausePanic()`, которая вызывает `panic("что-то пошло не так!")`.

2. **Обработка паники с `recover`**
    - Создайте функцию `HandlePanic()`, которая:
        - вызывает `CausePanic()`,
        - использует `defer` + `recover` для перехвата паники,
        - при панике выводит `"Паника перехвачена: ..."`, но программа продолжает выполняться.

3. **Запуск с демонстрацией поведения**
    - В `main()` вызовите `CausePanic()` напрямую и посмотрите, что произойдет.
    - Затем вызовите `HandlePanic()` и убедитесь, что паника обработана.

```go
func CausePanic() {
	panic("что-то пошло не так!")
}

func HandlePanic() {
	defer func(){
		if r := recover(); r != nil {
			fmt.Println("паника перехвачена: ", r)
		}
	}()

	CausePanic()
}

```


## Обработка паники при делении на ноль

### 2

В этом задании вы научитесь использовать `panic` и `recover` для безопасного выполнения кода.

## Задание

1. Реализуйте функцию `SafeDivide(a, b int) int`, которая:
    - Вызывает `panic("деление на ноль")`, если `b == 0`.
    - Обрабатывает панику с `recover()` и возвращает `0` вместо аварийного завершения.

2. В `main()` протестируйте вызовы:
```go
   SafeDivide(10, 2) // Ожидаемый результат: 5
   SafeDivide(10, 0) // Ожидаемый результат: 0 (без паники)
   
func SafeDivide(a, b int) (out int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic: ", r)
			out = 0
		}
	}()
	
	if b == 0 {
		panic("деление на ноль")
	}
	
	return a/b
}
```


## Обработка паники в многоуровневых вызовах функций

### 3

В этом задании вы разберетесь, как `panic` поднимается по стеку вызовов и как `recover` может ее обработать.

### Задание

1. Реализуйте три функции:
    - `Level1()`, которая вызывает `Level2()`.
    - `Level2()`, которая вызывает `Level3()`.
    - `Level3()`, которая вызывает `panic("ошибка в Level3")`.

2. В `Level1()` используйте `recover()` для перехвата паники и выведите:
   Паника обработана на уровне 1: ошибка в Level3

3.  В `Level2()` добавьте `defer`, который печатает `"Завершаем Level2"`, чтобы убедиться, что `defer` выполняется даже при панике.

4. В `main()` вызовите `Level1()` и убедитесь, что программа не завершилась аварийно.

```
func Level1() {
	defer func(){
		if r := recover(); r != nil {
			fmt.Println("Паника обработана на уровне 1:", r)
		}
	}()

	Level2()
}

func Level2() {
	defer func(){
		fmt.Println("Завершаем Level2")
	}()

	Level3()

}

func Level3() {
	panic("ошибка в Level3")
}
/*
print:

Завершаем Level2
Паника обработана на уровне 1: ошибка в Level3

*/
```