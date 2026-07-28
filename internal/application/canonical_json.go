package application

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

const maxCanonicalInteger = int64(1<<53 - 1)

type canonicalParser struct {
	raw []byte
	pos int
}

func parseCanonicalJSON(raw []byte) (any, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("is empty")
	}
	if !utf8.Valid(raw) {
		return nil, fmt.Errorf("contains invalid UTF-8")
	}
	parser := canonicalParser{raw: raw}
	value, err := parser.parseValue()
	if err != nil {
		return nil, err
	}
	parser.skipWhitespace()
	if parser.pos != len(parser.raw) {
		return nil, fmt.Errorf("contains trailing data")
	}
	return value, nil
}

func (p *canonicalParser) parseValue() (any, error) {
	p.skipWhitespace()
	if p.pos >= len(p.raw) {
		return nil, fmt.Errorf("ended before a value")
	}
	switch p.raw[p.pos] {
	case '{':
		return p.parseObject()
	case '[':
		return p.parseArray()
	case '"':
		value, err := p.parseString()
		if err != nil {
			return nil, err
		}
		return norm.NFC.String(value), nil
	case 't':
		return p.parseLiteral("true", true)
	case 'f':
		return p.parseLiteral("false", false)
	case 'n':
		return p.parseLiteral("null", nil)
	default:
		return p.parseInteger()
	}
}

func (p *canonicalParser) parseObject() (map[string]any, error) {
	p.pos++
	result := make(map[string]any)
	rawKeys := make(map[string]struct{})
	normalizedKeys := make(map[string]struct{})
	p.skipWhitespace()
	if p.consume('}') {
		return result, nil
	}
	for {
		p.skipWhitespace()
		if p.pos >= len(p.raw) || p.raw[p.pos] != '"' {
			return nil, fmt.Errorf("object key must be a string")
		}
		rawKey, err := p.parseString()
		if err != nil {
			return nil, err
		}
		if _, exists := rawKeys[rawKey]; exists {
			return nil, fmt.Errorf("contains duplicate object key %q", rawKey)
		}
		rawKeys[rawKey] = struct{}{}
		key := norm.NFC.String(rawKey)
		if _, exists := normalizedKeys[key]; exists {
			return nil, fmt.Errorf("contains object keys that collide after NFC normalization")
		}
		normalizedKeys[key] = struct{}{}

		p.skipWhitespace()
		if !p.consume(':') {
			return nil, fmt.Errorf("object key %q is missing ':'", key)
		}
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		result[key] = value
		p.skipWhitespace()
		if p.consume('}') {
			return result, nil
		}
		if !p.consume(',') {
			return nil, fmt.Errorf("object entries must be separated by ','")
		}
	}
}

func (p *canonicalParser) parseArray() ([]any, error) {
	p.pos++
	result := make([]any, 0)
	p.skipWhitespace()
	if p.consume(']') {
		return result, nil
	}
	for {
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		result = append(result, value)
		p.skipWhitespace()
		if p.consume(']') {
			return result, nil
		}
		if !p.consume(',') {
			return nil, fmt.Errorf("array entries must be separated by ','")
		}
	}
}

func (p *canonicalParser) parseString() (string, error) {
	if !p.consume('"') {
		return "", fmt.Errorf("expected a string")
	}
	var builder strings.Builder
	for p.pos < len(p.raw) {
		current := p.raw[p.pos]
		if current == '"' {
			p.pos++
			return builder.String(), nil
		}
		if current == '\\' {
			p.pos++
			if p.pos >= len(p.raw) {
				return "", fmt.Errorf("string ends after an escape")
			}
			escape := p.raw[p.pos]
			p.pos++
			switch escape {
			case '"', '\\', '/':
				builder.WriteByte(escape)
			case 'b':
				builder.WriteByte('\b')
			case 'f':
				builder.WriteByte('\f')
			case 'n':
				builder.WriteByte('\n')
			case 'r':
				builder.WriteByte('\r')
			case 't':
				builder.WriteByte('\t')
			case 'u':
				value, err := p.parseHexRune()
				if err != nil {
					return "", err
				}
				if utf16.IsSurrogate(value) {
					if value < 0xD800 || value > 0xDBFF || p.pos+2 > len(p.raw) || p.raw[p.pos] != '\\' || p.raw[p.pos+1] != 'u' {
						return "", fmt.Errorf("contains an unpaired Unicode surrogate")
					}
					p.pos += 2
					low, err := p.parseHexRune()
					if err != nil || low < 0xDC00 || low > 0xDFFF {
						return "", fmt.Errorf("contains an invalid Unicode surrogate pair")
					}
					value = utf16.DecodeRune(value, low)
				}
				builder.WriteRune(value)
			default:
				return "", fmt.Errorf("contains invalid escape \\%c", escape)
			}
			continue
		}
		if current < 0x20 {
			return "", fmt.Errorf("contains an unescaped control character")
		}
		runeValue, size := utf8.DecodeRune(p.raw[p.pos:])
		if runeValue == utf8.RuneError && size == 1 {
			return "", fmt.Errorf("contains invalid UTF-8")
		}
		builder.WriteRune(runeValue)
		p.pos += size
	}
	return "", fmt.Errorf("unterminated string")
}

func (p *canonicalParser) parseHexRune() (rune, error) {
	if p.pos+4 > len(p.raw) {
		return 0, fmt.Errorf("contains an incomplete Unicode escape")
	}
	value, err := strconv.ParseUint(string(p.raw[p.pos:p.pos+4]), 16, 16)
	if err != nil {
		return 0, fmt.Errorf("contains an invalid Unicode escape")
	}
	p.pos += 4
	return rune(value), nil
}

func (p *canonicalParser) parseInteger() (any, error) {
	start := p.pos
	if p.consume('-') && p.pos >= len(p.raw) {
		return nil, fmt.Errorf("contains an incomplete number")
	}
	if p.pos >= len(p.raw) || p.raw[p.pos] < '0' || p.raw[p.pos] > '9' {
		return nil, fmt.Errorf("contains an invalid value")
	}
	if p.raw[p.pos] == '0' {
		p.pos++
		if p.pos < len(p.raw) && p.raw[p.pos] >= '0' && p.raw[p.pos] <= '9' {
			return nil, fmt.Errorf("number contains a leading zero")
		}
	} else {
		for p.pos < len(p.raw) && p.raw[p.pos] >= '0' && p.raw[p.pos] <= '9' {
			p.pos++
		}
	}
	if p.pos < len(p.raw) && (p.raw[p.pos] == '.' || p.raw[p.pos] == 'e' || p.raw[p.pos] == 'E') {
		return nil, fmt.Errorf("only schema-declared integers are supported")
	}
	value, err := strconv.ParseInt(string(p.raw[start:p.pos]), 10, 64)
	if err != nil || value < -maxCanonicalInteger || value > maxCanonicalInteger {
		return nil, fmt.Errorf("integer exceeds the safe integer range")
	}
	if value == 0 {
		return int64(0), nil
	}
	return value, nil
}

func (p *canonicalParser) parseLiteral(literal string, value any) (any, error) {
	if p.pos+len(literal) > len(p.raw) || string(p.raw[p.pos:p.pos+len(literal)]) != literal {
		return nil, fmt.Errorf("contains an invalid literal")
	}
	p.pos += len(literal)
	return value, nil
}

func (p *canonicalParser) skipWhitespace() {
	for p.pos < len(p.raw) {
		switch p.raw[p.pos] {
		case ' ', '\t', '\r', '\n':
			p.pos++
		default:
			return
		}
	}
}

func (p *canonicalParser) consume(expected byte) bool {
	if p.pos >= len(p.raw) || p.raw[p.pos] != expected {
		return false
	}
	p.pos++
	return true
}

func canonicalJSON(value any) []byte {
	return appendCanonicalJSON(nil, value)
}

func appendCanonicalJSON(target []byte, value any) []byte {
	switch typed := value.(type) {
	case nil:
		return append(target, "null"...)
	case bool:
		if typed {
			return append(target, "true"...)
		}
		return append(target, "false"...)
	case int64:
		return strconv.AppendInt(target, typed, 10)
	case string:
		return appendCanonicalString(target, typed)
	case []any:
		target = append(target, '[')
		for index, element := range typed {
			if index > 0 {
				target = append(target, ',')
			}
			target = appendCanonicalJSON(target, element)
		}
		return append(target, ']')
	case map[string]any:
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(left, right int) bool {
			return utf16Less(keys[left], keys[right])
		})
		target = append(target, '{')
		for index, key := range keys {
			if index > 0 {
				target = append(target, ',')
			}
			target = appendCanonicalString(target, key)
			target = append(target, ':')
			target = appendCanonicalJSON(target, typed[key])
		}
		return append(target, '}')
	default:
		panic(fmt.Sprintf("unsupported canonical JSON type %T", value))
	}
}

func appendCanonicalString(target []byte, value string) []byte {
	const hexDigits = "0123456789abcdef"
	target = append(target, '"')
	for _, runeValue := range value {
		switch runeValue {
		case '"', '\\':
			target = append(target, '\\', byte(runeValue))
		case '\b':
			target = append(target, '\\', 'b')
		case '\t':
			target = append(target, '\\', 't')
		case '\n':
			target = append(target, '\\', 'n')
		case '\f':
			target = append(target, '\\', 'f')
		case '\r':
			target = append(target, '\\', 'r')
		default:
			if runeValue < 0x20 {
				target = append(target, '\\', 'u', '0', '0', hexDigits[(runeValue>>4)&0xF], hexDigits[runeValue&0xF])
				continue
			}
			target = utf8.AppendRune(target, runeValue)
		}
	}
	return append(target, '"')
}

func utf16Less(left, right string) bool {
	leftUnits := utf16.Encode([]rune(left))
	rightUnits := utf16.Encode([]rune(right))
	limit := len(leftUnits)
	if len(rightUnits) < limit {
		limit = len(rightUnits)
	}
	for index := 0; index < limit; index++ {
		if leftUnits[index] != rightUnits[index] {
			return leftUnits[index] < rightUnits[index]
		}
	}
	return len(leftUnits) < len(rightUnits)
}
