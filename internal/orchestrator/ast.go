package orchestrator

import (
	"errors"
	"strconv"
	"strings"
)

// узел AST
type Node struct {
	Op    string
	Left  *Node
	Right *Node
}

// функция парсинга выражения в AST
func ParseExpression(expression string) (*Node, error) {
	tokens, err := Tokenize(expression)
	if err != nil {
		return nil, err
	}
	pos := 0
	node, err := ParseExpr(&tokens, &pos)
	if err != nil {
		return nil, err
	}
	// Если после парсинга остались неиспользованные токены, значит, выражение некорректно
	if pos != len(tokens) {
		// Invalid syntax: unexpected tokens remaining
		return nil, errors.New("Invalid data")
	}
	return node, nil
}

// функция токенизации (разбиение на числа, операторы, скобки)
func Tokenize(expression string) ([]string, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	tokens := []string{}
	num := ""
	lastToken := ""

	for _, char := range expression {
		switch {
		case char >= '0' && char <= '9':
			num += string(char)
		case char == '+' || char == '-' || char == '*' || char == '/' || char == '(' || char == ')':
			if num != "" {
				tokens = append(tokens, num)
				lastToken = num
				num = ""
			}
			// Проверка на два арифметических знака подряд
			if (char == '+' || char == '-' || char == '*' || char == '/') &&
				(lastToken == "+" || lastToken == "-" || lastToken == "*" || lastToken == "/") {
				// Invalid syntax: two consecutive operators
				return nil, errors.New("Invalid data")
			}
			// Проверка на отсутствие знака между скобками
			if char == '(' && (lastToken == ")" || (lastToken != "" && isNumber(lastToken))) {
				// Invalid syntax: missing operator between parentheses
				return nil, errors.New("Invalid data")
			}
			tokens = append(tokens, string(char))
			lastToken = string(char)
		default:
			// Invalid syntax: contains invalid characters
			return nil, errors.New("Invalid data")
		}
	}
	if num != "" {
		tokens = append(tokens, num)
	}
	// Проверка, что выражение не заканчивается оператором
	if lastToken == "+" || lastToken == "-" || lastToken == "*" || lastToken == "/" {
		// Invalid syntax: expression cannot end with an operator
		return nil, errors.New("Invalid data")
	}
	return tokens, nil
}

// парсим выражение с учетом приоритетов и скобок
func ParseExpr(tokens *[]string, pos *int) (*Node, error) {
	defer func() {
		if r := recover(); r != nil {
			*pos = len(*tokens) // Останавливаем парсинг
		}
	}()

	node, err := ParseTerm(tokens, pos)
	if err != nil {
		return nil, err
	}

	for *pos < len(*tokens) {
		op := (*tokens)[*pos]
		if op != "+" && op != "-" {
			break
		}
		*pos++

		rightNode, err := ParseTerm(tokens, pos)
		if err != nil {
			return nil, err
		}

		node = &Node{Op: op, Left: node, Right: rightNode}
	}

	return node, nil
}

// учитываем умножение и деление
func ParseTerm(tokens *[]string, pos *int) (*Node, error) {
	node, err := ParseFactor(tokens, pos)
	if err != nil {
		return nil, err
	}

	for *pos < len(*tokens) {
		op := (*tokens)[*pos]
		if op != "*" && op != "/" {
			break
		}
		*pos++

		rightNode, err := ParseFactor(tokens, pos)
		if err != nil {
			return nil, err
		}

		node = &Node{Op: op, Left: node, Right: rightNode}
	}

	return node, nil
}

// числа и выражения в скобках
func ParseFactor(tokens *[]string, pos *int) (*Node, error) {
	if *pos >= len(*tokens) {
		// Invalid syntax: unexpected end of expression
		return nil, errors.New("Invalid data")
	}

	token := (*tokens)[*pos]
	*pos++

	if token == "(" {
		node, err := ParseExpr(tokens, pos)
		if err != nil {
			return nil, err
		}
		if *pos >= len(*tokens) || (*tokens)[*pos] != ")" {
			// Invalid syntax: unbalanced parentheses
			return nil, errors.New("Invalid data")
		}
		*pos++ // Пропускаем ')'
		return node, nil
	}

	// Проверяем, является ли токен числом
	if _, err := strconv.Atoi(token); err != nil {
		// Invalid syntax: expected a number
		return nil, errors.New("Invalid data")
	}

	return &Node{Op: token}, nil
}