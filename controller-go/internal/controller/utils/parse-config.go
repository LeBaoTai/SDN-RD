package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/LeBaoTai/myco-controller/internal/oc"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/goyang/pkg/yang"
	"github.com/openconfig/ygot/ytypes"
)

func buildGNMIPathWithSchema(schema *yang.Entry, tokens []string) (*gnmi.Path, error) {
	var elems []*gnmi.PathElem
	cur := schema
	i := 0
	// Loop trong — chỉ chạy trong phạm vi 1 dòng
	for i < len(tokens) {
		name := tokens[i]
		child, ok := cur.Dir[name]
		if !ok {
			return nil, fmt.Errorf("không tìm thấy node %q trong schema (dòng path: %v)", name, tokens)
		}
		elem := &gnmi.PathElem{Name: name}
		if child.IsList() {
			i++
			if i >= len(tokens) {
				return nil, fmt.Errorf("thiếu giá trị key cho list %q", name)
			}
			elem.Key = map[string]string{child.Key: tokens[i]}
		}
		elems = append(elems, elem)
		cur = child
		i++
	}
	return &gnmi.Path{Elem: elems}, nil
}

func ProcessFlatFile(filepath string, device *oc.Device, schema *ytypes.Schema) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		line = strings.TrimPrefix(line, "set /")
		line = strings.TrimPrefix(line, "set")
		line = strings.TrimSpace(line)

		tokens := strings.Fields(line)
		if len(tokens) < 2 {
			log.Printf("Skip this line: %v, missing path/value: %q", lineNumber, line)
		}
		value := tokens[len(tokens)-1]
		pathTokens := tokens[:len(tokens)-1]

		path, err := buildGNMIPathWithSchema(schema.RootSchema(), pathTokens)
		if err != nil {
			log.Printf("Skip this line: %v", line)
			continue
		}

		if err := ytypes.SetNode(
			schema.RootSchema(),
			device,
			path,
			&gnmi.TypedValue{Value: &gnmi.TypedValue_StringVal{StringVal: value}},
			&ytypes.InitMissingElements{},
		); err != nil {
			log.Printf("Error on line %d (%q): %v", lineNumber, line, err)
			continue
		}
		return scanner.Err()
	}
	return scanner.Err()
}

func ParseBraceConfig(filePath string) (map[string]interface{}, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	root := make(map[string]interface{})
	stack := []map[string]interface{}{root}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		cur := stack[len(stack)-1]

		switch {
		case strings.HasSuffix(line, "{"):
			// mở container mới, ví dụ: "system {" hoặc "information {"
			key := strings.TrimSpace(strings.TrimSuffix(line, "{"))
			newNode := make(map[string]interface{})
			cur[key] = newNode
			stack = append(stack, newNode)

		case line == "}":
			// đóng container, quay lại cấp cha
			if len(stack) > 1 {
				stack = stack[:len(stack)-1]
			}

		default:
			// leaf: "key value" hoặc "key \"value có khoảng trắng\""
			key, val, err := parseLeaf(line)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Bỏ qua dòng lỗi: %q (%v)\n", line, err)
				continue
			}
			cur[key] = val
		}
	}
	return root, scanner.Err()
}

func parseLeaf(line string) (string, interface{}, error) {
	// xử lý giá trị có dấu ngoặc kép, ví dụ: contact "admin@example.com"
	if idx := strings.Index(line, "\""); idx != -1 {
		key := strings.TrimSpace(line[:idx])
		val := strings.Trim(line[idx:], "\"")
		return key, val, nil
	}

	tokens := strings.Fields(line)
	if len(tokens) < 2 {
		return "", nil, fmt.Errorf("dòng thiếu key/value")
	}
	key := tokens[0]
	valStr := strings.Join(tokens[1:], " ")

	// thử convert sang kiểu số/bool nếu phù hợp
	if n, err := strconv.Atoi(valStr); err == nil {
		return key, n, nil
	}
	if b, err := strconv.ParseBool(valStr); err == nil {
		return key, b, nil
	}
	return key, valStr, nil
}
