package mapstasks

import (
	"fmt"
	"slices"
	"strings"
)

/*
# 2. Частотный анализ слов

## Описание задания

В этом задании вам необходимо реализовать программу на Go, которая проводит частотный анализ слов в заданном тексте. Ваша программа должна использовать `map` для подсчета количества повторений каждого слова и выводить результат, отсортированный по убыванию частоты.

## Задачи

1. **Реализуйте функцию `WordFrequency(text string) map[string]int`:**
    - Принимает строку `text`.
    - Разбивает строку на слова (например, с помощью `strings.Fields`).
    - Подсчитывает количество повторений каждого слова.
    - Возвращает `map[string]int`, где ключ – слово, а значение – количество его вхождений.

2. **Реализуйте функцию `PrintWordFrequency(freqMap map[string]int)`:**
    - Принимает `map[string]int` с данными о частоте слов.
    - Выводит слова и их количество, отсортированные по убыванию частоты.
*/

func WordFrequency(text string) map[string]int {
	m := make(map[string]int)
	words := strings.Fields(text)

	for _, w := range words {
		if count, ex := m[w]; ex {
			m[w] = count + 1
		} else {
			m[w] = 1
		}
	}

	return m
}

func PrintWordFrequency(freqMap map[string]int) {
	type node struct {
		word  string
		count int
	}

	nodes := make([]node, 0, len(freqMap))
	for k, v := range freqMap {
		nodes = append(nodes, node{word: k, count: v})
	}

	slices.SortFunc(nodes, func(a node, b node) int {
		if a.count == b.count {
			return 0
		}

		if a.count > b.count {
			return 1
		} else {
			return -1
		}
	})

	format := "word: %s, count: %d\n"
	for _, n := range nodes {
		fmt.Printf(format, n.word, n.count)
	}
}
