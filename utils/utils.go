package utils

import (
	"fmt"
)

func Inp(s string) ([]int, error) {
	throws := make([]int, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == ' ' || c == '\t' {
			continue
		}
		// c = byte(unicode.ToLower(rune(c)))

		switch c {
		case 'x', 'X':
			throws = append(throws, 10)
		case '-':
			throws = append(throws, 0)
		case '/':
			if len(throws) == 0 {
				return nil, fmt.Errorf("'/'' без предыдущего символа")
			}
			pr := throws[len(throws)-1]
			if pr == 10 {
				return nil, fmt.Errorf("'/'' не может идти сразу после X")
			}
			throws = append(throws, 10-pr)
		default:
			if c < '0' || c > '9' {
				return nil, fmt.Errorf("некорректный символ %c", c)
			}
			throws = append(throws, int(c-'0'))
		}
	}
	return throws, nil
}

func Scr(throws []int) (int, error) {
	ttl := 0
	i := 0
	for frame := 1; frame <= 9; frame++ {
		if i >= len(throws) {
			return 0, fmt.Errorf("недостаточно бросков для фрейма %d", frame)
		}
		if throws[i] == 10 {
			if i+2 >= len(throws) {
				return 0, fmt.Errorf("недостаточно бонусных бросков после страйка на фрейме %d", frame)
			}
			ttl += 10 + throws[i+1] + throws[i+2]
			i++
			continue
		}
		if i+1 >= len(throws) {
			return 0, fmt.Errorf("не хватает второго броска для фрейма %d", frame)
		}
		tr1 := throws[i]
		tr2 := throws[i+1]
		if tr1 < 0 || tr1 > 10 || tr2 < 0 || tr2 > 10 {
			return 0, fmt.Errorf("некорректное значение броска в фрейме %d", frame)
		}
		sum := tr1 + tr2
		if sum > 10 {
			return 0, fmt.Errorf("сумма бросков превышает 10 в фрейме %d", frame)
		}
		if sum == 10 {
			if i+2 >= len(throws) {
				return 0, fmt.Errorf("недостаточно бонусных бросков после spare на фрейме %d", frame)
			}
			ttl += 10 + throws[i+2]
		} else {
			ttl += sum
		}
		i += 2
	}
	if i >= len(throws) {
		return 0, fmt.Errorf("недостаточно бросков для фрейма 10")
	}
	if throws[i] == 10 {
		if i+2 >= len(throws) {
			return 0, fmt.Errorf("в 10-м фрейме после страйка нужно еще 2 броска")
		}
		ttl += throws[i] + throws[i+1] + throws[i+2]
		return ttl, nil
	}
	if i+1 >= len(throws) {
		return 0, fmt.Errorf("не хватает второго броска для фрейма 10")
	}
	tr1 := throws[i]
	tr2 := throws[i+1]
	if tr1 < 0 || tr1 > 10 || tr2 < 0 || tr2 > 10 {
		return 0, fmt.Errorf("некорректное значение броска в фрейме 10")
	}
	sum := tr1 + tr2
	if sum == 10 {
		if i+2 >= len(throws) {
			return 0, fmt.Errorf("в 10-м фрейме после spare нужен еще 1 бросок")
		}
		ttl += throws[i] + throws[i+1] + throws[i+2]
		return ttl, nil
	}
	if sum > 10 {
		return 0, fmt.Errorf("сумма бросков превышает 10 в фрейме 10")
	}
	ttl += sum
	return ttl, nil
}
