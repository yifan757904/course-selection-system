package model

import (
	"errors"
	"fmt"
	"time"
)

// 学期常量
const (
	FirstSemester  = "1" // 第一学期
	SecondSemester = "2" // 第二学期
)

// SemesterConfig 学期配置
type SemesterConfig struct {
	FirstSemesterStart  time.Month // 第一学期开始月份 (如2月)
	FirstSemesterEnd    time.Month // 第一学期结束月份 (如6月)
	SecondSemesterStart time.Month // 第二学期开始月份 (如9月)
	SecondSemesterEnd   time.Month // 第二学期结束月份 (如次年1月)
}

// DefaultSemesterConfig 默认学期配置（可根据实际情况调整）
var DefaultSemesterConfig = SemesterConfig{
	FirstSemesterStart:  3, // 2月开学
	FirstSemesterEnd:    8, // 6月结束
	SecondSemesterStart: 9, // 9月开学
	SecondSemesterEnd:   2, // 次年1月结束
}

// GetSemesterByDate 根据日期获取学期
func GetSemesterByDate(date time.Time, config SemesterConfig) string {
	month := date.Month()
	year := date.Year()

	switch {
	// 第二学期范围判断（跨年处理）
	case month >= config.SecondSemesterStart || month <= config.SecondSemesterEnd:
		if month >= config.SecondSemesterStart {
			return fmt.Sprintf("%d-%s", year, SecondSemester)
		}
		return fmt.Sprintf("%d-%s", year-1, SecondSemester) // 跨年处理

	// 第一学期范围判断
	case month >= config.FirstSemesterStart && month <= config.FirstSemesterEnd:
		return fmt.Sprintf("%d-%s", year, FirstSemester)
	default:
		return ""
	}
}

// ValidateSemester 验证学期格式
func ValidateSemester(semester string) error {
	if len(semester) != 6 || semester[4] != '-' {
		return errors.New("invalid semester format, should be YYYY-S")
	}

	year := semester[:4]
	term := semester[5:]

	if term != FirstSemester && term != SecondSemester {
		return errors.New("invalid term, should be 1 or 2")
	}

	if _, err := time.Parse("2006", year); err != nil {
		return fmt.Errorf("invalid year: %v", err)
	}

	return nil
}
