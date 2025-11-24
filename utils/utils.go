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

	for frame := 1; frame <= 10; frame++ {
		if i >= len(throws) {
			return 0, fmt.Errorf("недостаточно бросков для фрейма %d", frame)
		}

		if throws[i] == 10 {
			if i+2 >= len(throws) {
				return 0, fmt.Errorf("недостаточно бонусных бросков после страйка (фрейм %d)", frame)
			}
			ttl += 10 + throws[i+1] + throws[i+2]
			i++
			continue
		}
		if i+1 >= len(throws) {
			return 0, fmt.Errorf("не хватает второго броска для фрейма %d", frame)
		}
		first := throws[i]
		second := throws[i+1]
		if first < 0 || first > 10 || second < 0 || second > 10 {
			return 0, fmt.Errorf("некорректное значение броска в фрейме %d", frame)
		}
		sum := first + second
		if sum == 10 {
			if i+2 >= len(throws) {
				return 0, fmt.Errorf("недостаточно бонусных бросков после спэа (фрейм %d)", frame)
			}
			ttl += 10 + throws[i+2]
			i += 2
		} else if sum < 10 {
			ttl += sum
			i += 2
		} else {
			return 0, fmt.Errorf("сумма бросков превышает 10 в фрейме %d", frame)
		}
	}

	return ttl, nil
}

func ScrPart(throws []int) int {
	ttl := 0
	i := 0

	for frame := 1; frame <= 10; frame++ {
		if i >= len(throws) {
			break
		}

		if throws[i] == 10 {
			if i+2 >= len(throws) {
				break
			}
			ttl += 10 + throws[i+1] + throws[i+2]
			i++
			continue
		}
		if i+1 >= len(throws) {
			break
		}
		first := throws[i]
		second := throws[i+1]
		sum := first + second
		if sum == 10 {
			if i+2 >= len(throws) {
				break
			}
			ttl += 10 + throws[i+2]
			i += 2
		} else {
			ttl += sum
			i += 2
		}
	}

	return ttl
}
