package filewords

import (
	"bufio"
	"os"
)

/*
Задание: Параллельный подсчет слов в файлах с использованием паттерна Fan-Out

### Цель задания
Реализовать параллельную обработку текстовых файлов с использованием паттерна **Fan-Out**, чтобы ускорить подсчет слов в каждом файле.

### Описание задачи
Есть директория с текстовыми файлами. Нужно:
1. Прочитать все файлы.
2. Распределить их обработку между несколькими горутинами.
3. Подсчитать количество слов в каждом файле.
4. Вывести общую статистику.

### Требования
- Использовать паттерн **Fan-Out** для распределения задач.
- Обработка каждого файла должна выполняться в отдельной горутине.
- Результаты должны агрегироваться в основном потоке.

*/

type WordCountResult struct {
	File      string
	WordCount int
	Err       error
}

func CountWordsInFiles(paths []string) []<-chan WordCountResult {
	outChannels := make([]<-chan WordCountResult, 0, len(paths))

	for _, p := range paths {
		resCh := countJob(p)
		outChannels = append(outChannels, resCh)
	}

	return outChannels
}

func countJob(path string) <-chan WordCountResult {
	out := make(chan WordCountResult)

	go func() {
		defer close(out)

		count, err := countWordInFile(path)
		out <- WordCountResult{
			File:      path,
			WordCount: count,
			Err:       err,
		}
	}()

	return out
}

func countWordInFile(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}

	count := 0

	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)

	for scan.Scan() {
		count++
	}

	return count, nil
}
