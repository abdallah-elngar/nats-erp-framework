package scheduler

import (
	"fmt"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

// CronSchedule مساعد لإنشاء جدولة Cron
type CronSchedule struct {
	schedule string
}

// NewCronSchedule ينشئ جدولة Cron جديدة
func NewCronSchedule() *CronSchedule {
	return &CronSchedule{}
}

// EverySecond كل ثانية
func (cs *CronSchedule) EverySecond() *CronSchedule {
	cs.schedule = "* * * * * *"
	return cs
}

// EveryMinute كل دقيقة
func (cs *CronSchedule) EveryMinute() *CronSchedule {
	cs.schedule = "0 * * * * *"
	return cs
}

// EveryHour كل ساعة
func (cs *CronSchedule) EveryHour() *CronSchedule {
	cs.schedule = "0 0 * * * *"
	return cs
}

// EveryDay كل يوم
func (cs *CronSchedule) EveryDay() *CronSchedule {
	cs.schedule = "0 0 0 * * *"
	return cs
}

// EveryWeek كل أسبوع
func (cs *CronSchedule) EveryWeek() *CronSchedule {
	cs.schedule = "0 0 0 * * 0"
	return cs
}

// EveryMonth كل شهر
func (cs *CronSchedule) EveryMonth() *CronSchedule {
	cs.schedule = "0 0 0 1 * *"
	return cs
}

// EveryYear كل سنة
func (cs *CronSchedule) EveryYear() *CronSchedule {
	cs.schedule = "0 0 0 1 1 *"
	return cs
}

// At ساعة محددة
func (cs *CronSchedule) At(hour, minute int) *CronSchedule {
	cs.schedule = fmt.Sprintf("0 %d %d * * *", minute, hour)
	return cs
}

// DailyAt يومياً في ساعة محددة
func (cs *CronSchedule) DailyAt(hour, minute int) *CronSchedule {
	return cs.At(hour, minute)
}

// EveryMinutes كل عدد محدد من الدقائق
func (cs *CronSchedule) EveryMinutes(minutes int) *CronSchedule {
	cs.schedule = fmt.Sprintf("0 */%d * * * *", minutes)
	return cs
}

// EveryHours كل عدد محدد من الساعات
func (cs *CronSchedule) EveryHours(hours int) *CronSchedule {
	cs.schedule = fmt.Sprintf("0 0 */%d * * *", hours)
	return cs
}

// Custom جدولة مخصصة
func (cs *CronSchedule) Custom(schedule string) *CronSchedule {
	cs.schedule = schedule
	return cs
}

// Build يبني الجدولة
func (cs *CronSchedule) Build() string {
	return cs.schedule
}

// Validate يتحقق من صحة الجدولة
func (cs *CronSchedule) Validate() error {
	if cs.schedule == "" {
		return fmt.Errorf("schedule is empty")
	}

	parts := strings.Split(cs.schedule, " ")
	if len(parts) != 6 {
		return fmt.Errorf("invalid cron format, expected 6 parts, got %d", len(parts))
	}

	return nil
}

// NextRun يحسب وقت التشغيل التالي
func (cs *CronSchedule) NextRun(from time.Time) (time.Time, error) {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(cs.schedule)
	if err != nil {
		return time.Time{}, err
	}

	return schedule.Next(from), nil
}
