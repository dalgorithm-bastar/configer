package util

import "strings"

// GetPrefix 去除空格并判断是否存在”/“,没有则添加
func GetPrefix(input string) string {
    if input == "" {
        return "/"
    }
    prefix := strings.TrimSpace(input)
    if prefix[len(prefix)-1] != '/' {
        prefix = prefix + "/"
    }
    return prefix
}

// Join 使用sep连接input
func Join(sep string, input ...string) string {
    return strings.Join(input, sep)
}

// ContainforSliceInOrder 检测字符串切片中是否按序包含目标字符串（不允许重复）
func ContainforSliceInOrder(inputSlice []string, targetString ...string) bool {
    position := 0
    for _, target := range targetString {
        contain := false
        for i, referString := range inputSlice {
            if target == referString {
                contain = true
                if i < position {
                    return false
                }
                position = i
            }
            if contain {
                break
            }
        }
        if !contain {
            return false
        }
    }
    return true
}
