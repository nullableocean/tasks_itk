package mapstasks

import (
	"fmt"
)

/*

# 3. Работа с map в Go: Фильтрация и инвертирование

Это задание направлено на освоение продвинутых операций с map в Go, включая фильтрацию по значениям и инвертирование ключей/значений.

## Задания

### Фильтрация по значению: `FilterByValue`

**Задача**:
Реализуйте функцию `FilterByValue`, которая фильтрует элементы map, оставляя только те, чьи значения присутствуют в разрешённом списке.

**Требования**:
- Функция должна принимать:
    - Исходную map типа `map[int]string`.
    - Список разрешённых значений типа `[]string`.
- Возвращает новую map, содержащую только элементы с значениями из списка.
- Исходная map не должна изменяться.
- Эффективная проверка значений (используйте set для оптимизации)`make(map[string]struct{}`.

*/

func FilterByValue(m map[int]string, allowedValues []string) map[int]string {
	newMap := make(map[int]string, len(m))
	allowedSet := make(map[string]struct{}, len(allowedValues))

	for _, v := range allowedValues {
		allowedSet[v] = struct{}{}
	}

	for k, v := range m {
		if _, ex := allowedSet[v]; ex {
			newMap[k] = v
		}
	}

	return newMap
}

func InvertMap(m map[string]int) (map[int]string, error) {
	invertMap := make(map[int]string, len(m))

	for k, v := range m {
		if _, ex := invertMap[v]; ex {
			return invertMap, fmt.Errorf("not unique value: %v", v)
		}

		invertMap[v] = k
	}

	return invertMap, nil
}
