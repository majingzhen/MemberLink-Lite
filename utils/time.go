package utils

import (
	"fmt"
	"time"
)

// ParseTimeRange 解析时间范围
func ParseTimeRange(startTimeStr, endTimeStr string) (*time.Time, *time.Time, error) {
	var startTime, endTime *time.Time
	var err error

	// 支持的时间格式
	timeFormats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02",
	}

	// 解析开始时间
	if startTimeStr != "" {
		var parsedTime time.Time
		for _, format := range timeFormats {
			parsedTime, err = time.Parse(format, startTimeStr)
			if err == nil {
				break
			}
		}
		if err != nil {
			return nil, nil, fmt.Errorf("开始时间格式错误: %s", startTimeStr)
		}
		startTime = &parsedTime
	}

	// 解析结束时间
	if endTimeStr != "" {
		var parsedTime time.Time
		for _, format := range timeFormats {
			parsedTime, err = time.Parse(format, endTimeStr)
			if err == nil {
				break
			}
		}
		if err != nil {
			return nil, nil, fmt.Errorf("结束时间格式错误: %s", endTimeStr)
		}

		// 如果只有日期，设置为当天的23:59:59
		if len(endTimeStr) == 10 { // YYYY-MM-DD 格式
			parsedTime = parsedTime.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
		endTime = &parsedTime
	}

	// 验证时间范围
	if startTime != nil && endTime != nil && startTime.After(*endTime) {
		return nil, nil, fmt.Errorf("开始时间不能晚于结束时间")
	}

	return startTime, endTime, nil
}

// FormatTime 格式化时间
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatTimePtr 格式化时间指针
func FormatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
