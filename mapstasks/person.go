package mapstasks

import "fmt"

/*

# 1. Задание: Работа с map в Go
## Описание
В этом задании вам нужно реализовать функции для работы с map в Go.
Вам предстоит создать, заполнить и обработать map, а затем выполнить некоторые операции с ним.
## Задачи
1. Создайте map, где ключ - это строка (имя человека), а значение - его возраст.
2. Добавьте в map несколько записей.
3. Реализуйте функцию `GetAge(name string) int`, которая возвращает возраст человека по его имени.
4. Реализуйте функцию `DeletePerson(name string)`, которая удаляет запись из map.
5. Реализуйте функцию `PrintAll()`, которая выводит все записи в map.

*/

var persons map[string]int

func init() {
	persons = make(map[string]int)

	AddPerson("Иван", 178)
	AddPerson("Саша", 128)
}

// Добавление записей
func AddPerson(name string, age int) {
	persons[name] = age
}

// Получение возраста
func GetAge(name string) int {
	return persons[name]
}

// Удаление записи
func DeletePerson(name string) {
	delete(persons, name)
}

// Вывод всех записей
func PrintAll() {
	for k, v := range persons {
		fmt.Printf("Имя: %s, Возраст: %d\n", k, v)
	}
}
