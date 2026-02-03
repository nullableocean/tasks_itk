# 4. Задание: Анализ кода на Go

Это задание направлено на глубокое понимание работы срезов (interface), их модификации и передачи в функциях Go.  
**Ваша задача:** Определить вывод программы и зафиксировать ответы **в сообщениях коммитов** с кратким объяснением логики.
package main

```
import "fmt"

type MyStruct struct {
	MyInt int
}

func func1() MyStruct {
	return MyStruct{MyInt: 1}
}

func func2() *MyStruct {
	return &MyStruct{}
}

func func3(s *MyStruct) {
	s.MyInt = 333
}

func func4(s MyStruct) {
	s.MyInt = 923
}

func func5() *MyStruct {
	return nil
}

func main() {
	ms1 := func1()
	fmt.Println(ms1.MyInt) // 1

	ms2 := func2()
	fmt.Println(ms2.MyInt) // 0

	func3(ms2)
	fmt.Println(ms2.MyInt) // 333

	func4(ms1)
	fmt.Println(ms1.MyInt) // 1

	ms5 := func5()
	fmt.Println(ms5.MyInt) // panic nil pointer
}
```

---------------------------------------------------------------
# Slice TASK1
# Задание: Анализ кода на Go

Это задание направлено на глубокое понимание работы срезов (slices), их модификации и передачи в функциях Go.  
**Ваша задача:** Определить вывод каждой из предложенных программ и зафиксировать ответы **в сообщениях коммитов** с кратким объяснением логики.
### 1.
```
package main

import "fmt"

type account struct {
	value int
}

func main() {
	s1 := make([]account, 0, 2)
	s1 = append(s1, account{})
	s2 := append(s1, account{})
	acc := &s2[0]
	acc.value = 100
	fmt.Println(s1, s2) // s1: [account{value:100}], s1 - len=1, s2: [account{value:100}, account{value:0}], s2 - len = 2
	s1 = append(s1, account{})
	acc.value += 100
	fmt.Println(s1, s2) // s1: [account{value:200}, account{value:0}], s2: [account{value:200}, account{value:0}]
}
```
-----
2.
```
package main

import "fmt"

func main() {
	slice := make([]string, 0, 5)
	slice = append(slice, "0", "1", "2", "3")
	fmt.Println(slice, len(slice), cap(slice)) // ["0", "1", "2", "3"], 4, 5
	addToSlice1(slice)
	fmt.Println(slice, len(slice), cap(slice)) // ["0", "1", "2", "one"], 4, 5
	addToSlice2(slice)
	fmt.Println(slice, len(slice), cap(slice)) // ["0", "1", "2", "3"], 4, 5
}

func addToSlice1(slice []string) {
	slice = append(slice[1:3], "one") // slice[1:3], len = 2, cap = 3, следующий элемент встанет без аллокаций нового массива на место элемента под индексом 3 в изначальном массиве
}

func addToSlice2(slice []string) {
	slice = append(slice, "two") // элемент не будет виден в переданном слайсе, так как длина у него не изменится
}
```
---
3.
```
package main

import "fmt"

func main() {
	a1 := make([]int, 0, 10)
	a1 = append(a1, []int{1, 2, 3, 4, 5}...) // len = 5
	a2 := append(a1, 6) // a2 len = 6
	a3 := append(a1, 7) // a3 len = 6, поскольку новых аллокаций не было, имеем дело с одним и тем же массивом, a3 перетрет элемент в a2[5]
	fmt.Println(a1, a2, a3) // [1, 2, 3, 4, 5], [1, 2, 3, 4, 5, 7], [1, 2, 3, 4, 5, 7]
}
```
---
4.
```
package main

import "fmt"

func main() {
	a := []int{1, 2, 3} // len 3, cap 3
	b := a[:2] // b = [1,2] len 2, cap 3
	b = append(b, 4) // b = [1,2,4], a = [1,2,4], перетираем значение в изначальном массиве
	fmt.Println(b) // [1,2,4]
	fmt.Println(a) // [1,2,4]
}
```
-----
5.
```
package main

import "fmt"

func main() {
	arr := []int{1, 2, 3}
	src := arr[:1] // [1], len = 1, cap = 3
	foo(src)
	fmt.Println(src) // [1]
	fmt.Println(arr) // [1, 5, 3]
}

func foo(src []int) {
	src = append(src, 5)
}
```
-----
6.
```
package main

import "fmt"

func main() {
	arr := [5]int{1, 2, 3, 4, 5} // len = 5
	bar := arr[1:3] // [2,3], len = 2, cap = 4 (текущая длина + кол-во оставшихся элементов в массиве)
	bar = append(bar, 10, 11, 12, 13) // новый массив для bar, так как вышли за cap, bar = [2,3,10,11,12,13], изначальный массив не меняется, могло быть иначе если бы последовательно аппендили
	fmt.Println(arr, bar) //
}
```
-----
7.
```
package main

import "fmt"

func main() {
	a := []string{"a", "b", "c"}
	b := a[1:2] 
	fmt.Println(b, cap(b), len(b)) // // b = [b], cap = 2, len = 1
	b[0] = "q"
	fmt.Println(a) // [a, q, c]
}
```
---
8.
```
package main

import (
	"fmt"
)

func main() {
	nums := make([]int, 1, 3) // len = 1, cap = 3
	fmt.Println(nums) // [0]
	appendSlice(nums, 1)
	fmt.Println(nums) // [0]
	copySlice(nums, []int{2, 3})
	fmt.Println(nums) // [2], копируем одно значение, так как len = 1, в 0 индекс
	mutateSlice(nums, 1, 4) // panic, out range, вышли индексом за размер
	fmt.Println(nums) //
}

func appendSlice(sl []int, val int) {
	sl = append(sl, val)
}

func copySlice(sl, cp []int) {
	copy(sl, cp)
}

func mutateSlice(sl []int, idx, val int) {
	sl[idx] = val
}
```
---
9.
```
package main

import (
	"fmt"
)

func main() {
	slice := make([]int, 3, 4) [0,0,0]
	appendingSlice(slice[:1]) // передали [0], cap = 4
	fmt.Println(slice) // [0, 1, 0]
}

func appendingSlice(slice []int) {
	slice = append(slice, 1) // добавляем без аллокации на место в изначальном массиве
}
```


SLICE TASK 2
# Задание: Анализ и исправление кода на Go

Это задание направлено на понимание работы срезов, функций и передачи данных в Go.  
**Ваша задача:**
1. **Проанализировать вывод программ** и объяснить поведение кода.
2. **Исправить код** так, чтобы достигался корректный результат (в некоторых случаях требуется несколько решений)


### 1.// Версия 1.21
```
package main

import (
	"fmt"
)

func main() {
	var numbers []*int
	for _, value := range []int{10, 20, 30, 40} { // в старых версиях в range value не копируется, а используется один адресс памяти
		// numbers = append(numbers, &value)
		
		num := value
		numbers = append(numbers, &num)
	}
	for _, number := range numbers {
		fmt.Println("d", *number)
	}
}
```
----
### 2.
```
package main

import (
	"fmt"
	"strings"
)

func chengeSlice(arr []string) {
	arr[0] = "Goodbye"
}

func appendSomeData(arr []string) []string {
	return append(arr, "!")
}

func main() {
	someSlice := []string{"Hello", "World"}
	chengeSlice(someSlice)
	someSlice = appendSomeData(someSlice) // fix: возвращаем слайс с добавленным значением
	fmt.Println(strings.Join(someSlice, ""))
}
```
----
### 3.
```
package main

import "fmt"

func test(testSlice []string) {
	testSlice = append(testSlice, "Пока")
}
func main() {
	testSlice := make([]string, 0, 3)
	testSlice = append(testSlice, "Привет")
	testSlice = append(testSlice, "Привет")
	test(testSlice) // для фикса нужно возвращать слайс с новым значением
	fmt.Println(testSlice) // [Привет, привет]
}
```
----
### 4.
```
package main

import "fmt"

func main() {
	first := []int{10, 20, 30, 40}
	second := make([]*int, len(first)) 
	for i, v := range first {
	    // возможный фикс в стархы версиях
	    // n := v
	    // second[i] = &n
		second[i] = &v
	}
	fmt.Println(*second[0], *second[1]) // 10 20
}
```
----
### 5.
```
package main

import (
	"fmt"
)

func main() {
	slice := make([]string, 3, 4)
	fmt.Println(slice) // ["","",""]

	appendSlice(slice)
	fmt.Println(slice) // ["","",""]

	mutareSlice(slice)
	fmt.Println(slice) // ["vasya","",""]
}

func appendSlice(slice []string) {
	slice = append(slice, "privet")
}
func mutareSlice(slice []string) {
	slice[0] = "vasya"
}
```



# SLICE TASK 4
### Работа со слайсами в Go

Этот проект демонстрирует различные способы работы со слайсами в Go, включая очистку, обнуление и особенности внутренней структуры.

**Ваша задача:** Определить вывод каждого случая и зафиксировать ответы **в сообщениях коммитов** с кратким объяснением логики.
package main

```
import (
	"fmt"
	"unsafe"
)

func main() {
	//1
	first := []int{1, 2, 3, 4, 5}
	first = nil
	fmt.Println("first:", first, ":", len(first), ":", cap(first)) // []:0:0

	//2
	second := []int{1, 2, 3, 4, 5}
	second = second[:0]
	fmt.Println("second:", second, ":", len(second), ":", cap(second)) // []:0:5

	//3
	third := []int{1, 2, 3, 4, 5}
	clear(third)
	fmt.Println("third:", third, ":", len(third), ":", cap(third)) // [0,0,0,0,0,]:5:5

	//4
	fourth := []int{1, 2, 3, 4, 5}
	clear(fourth[1:3])
	fmt.Println("fourth:", fourth, ":", len(fourth), ":", cap(fourth)) // [1,0,0,4,5]:5:5

	//5
	slice := make([]int, 3, 6)
	array := [3]int(slice[:3])
	slice[0] = 10

	fmt.Println("slice = ", slice, len(slice), cap(slice)) // [10, 0, 0], 3, 6
	fmt.Println("array =", array, len(array), cap(array)) // [0,0,0], 3, 3

	//6 В каких случаях Slice пустой или нулевой
	//1
	var data []string
	fmt.Println("var data []string:") // nil
	fmt.Printf("\tempty=%t nil=%t size=%d data=%p\n", len(data) == 0, data == nil, unsafe.Sizeof(data), unsafe.SliceData(data)) // nil слайс
	//2
	data = []string(nil)
	fmt.Println("data = []string(nil):")
	fmt.Printf("\tempty=%t nil=%t size=%d data=%p\n", len(data) == 0, data == nil, unsafe.Sizeof(data), unsafe.SliceData(data)) // nil слайс
	//3
	data = []string{}
	fmt.Println("data = []string{}")
	fmt.Printf("\tempty=%t nil=%t size=%d data=%p\n", len(data) == 0, data == nil, unsafe.Sizeof(data), unsafe.SliceData(data)) // пустой, не nil, инициализированный (аллоцированный)
	//4
	data = make([]string, 0)
	fmt.Println("data =make([]string,0)")
	fmt.Printf("\tempty=%t nil=%t size=%d data=%p\n", len(data) == 0, data == nil, unsafe.Sizeof(data), unsafe.SliceData(data)) // пустой, не nil, инициализированный (аллоцированный)

	empty := struct{}{}
	fmt.Println("empty struct address ", unsafe.Pointer(&empty)) // пустая структура и пустой слайс аллоцированный слайс смотрят в один в тот же адресс в памяти
}
```