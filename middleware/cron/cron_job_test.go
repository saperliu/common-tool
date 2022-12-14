package cron
import (
	"fmt"
	"github.com/robfig/cron"
	"testing"
	"time"
)

func task() {
	fmt.Println("I am a running job.")
}

func taskWithParams(a int, b string) {
	t:= time.Now()
	fmt.Printf("time %v",t )
	fmt.Println(a, b)
}

func TestNewCron(t *testing.T)  {
	c := cron.New()
	err:=c.AddFunc("0/10 * * * * *", func() { fmt.Println("Every hour on the 10 sec") })
	t.Log("----",err)
	c.Start()
	for{
		time.Sleep(10 * time.Second)
	}
}
func TestSecond(*testing.T) {
	defaultScheduler.Every(10).Second().Do(task)
	//defaultScheduler.Every(10).Second().Do(taskWithParams, 1, "hello")
	//Every(10).Second().Do(taskWithParams, 1, "hello")
	
	defaultScheduler.Start()
	for{
		time.Sleep(10 * time.Second)
	}


	//// Do jobs without params
	//gocron.Every(1).Second().Do(task)
	//gocron.Every(2).Seconds().Do(task)
	//gocron.Every(1).Minute().Do(task)
	//gocron.Every(2).Minutes().Do(task)
	//gocron.Every(1).Hour().Do(task)
	//gocron.Every(2).Hours().Do(task)
	//gocron.Every(1).Day().Do(task)
	//gocron.Every(2).Days().Do(task)
	//
	//// Do jobs on specific weekday
	//gocron.Every(1).Monday().Do(task)
	//gocron.Every(1).Thursday().Do(task)
	//
	//// function At() take a string like 'hour:min'
	//gocron.Every(1).Day().At("10:30").Do(task)
	//gocron.Every(1).Monday().At("18:30").Do(task)
	//
	//// remove, clear and next_run
	//_, time := gocron.NextRun()
	//fmt.Println(time)
	//
	//gocron.Remove(task)
	//gocron.Clear()
	//
	//// function Start start all the pending jobs
	//<- gocron.Start()
	//
	//// also , you can create a your new scheduler,
	//// to run two scheduler concurrently
	//s := gocron.NewScheduler()
	//s.Every(3).Seconds().Do(task)
	//<- s.Start()
}

// This is a basic test for the issue described here: https://github.com/jasonlvhit/gocron/issues/23
func TestScheduler_Weekdays(t *testing.T) {
	scheduler := NewScheduler()

	job1 := scheduler.Every(1).Monday().At("23:59")
	job2 := scheduler.Every(1).Wednesday().At("23:59")
	job1.Do(task)
	job2.Do(task)
	t.Logf("job1 scheduled for %s", job1.NextScheduledTime())
	t.Logf("job2 scheduled for %s", job2.NextScheduledTime())
	if job1.NextScheduledTime() == job2.NextScheduledTime() {
		t.Fail()
		t.Logf("Two jobs scheduled at the same time on two different weekdays should never run at the same time.[job1: %s; job2: %s]", job1.NextScheduledTime(), job2.NextScheduledTime())
	}
}

// This ensures that if you schedule a job for today's weekday, but the time is already passed, it will be scheduled for
// next week at the requested time.
func TestScheduler_WeekdaysTodayAfter(t *testing.T) {
	scheduler := NewScheduler()

	now := time.Now()
	timeToSchedule := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()-1, 0, 0, time.Local)

	job := callTodaysWeekday(scheduler.Every(1)).At(fmt.Sprintf("%02d:%02d", timeToSchedule.Hour(), timeToSchedule.Minute()))
	job.Do(task)
	t.Logf("job is scheduled for %s", job.NextScheduledTime())
	if job.NextScheduledTime().Weekday() != timeToSchedule.Weekday() {
		t.Fail()
		t.Logf("Job scheduled for current weekday for earlier time, should still be scheduled for current weekday (but next week)")
	}
	nextWeek := time.Date(now.Year(), now.Month(), now.Day()+7, now.Hour(), now.Minute()-1, 0, 0, time.Local)
	if !job.NextScheduledTime().Equal(nextWeek) {
		t.Fail()
		t.Logf("Job should be scheduled for the correct time next week.")
	}
}

// This is to ensure that if you schedule a job for today's weekday, and the time hasn't yet passed, the next run time
// will be scheduled for today.
func TestScheduler_WeekdaysTodayBefore(t *testing.T) {
	scheduler := NewScheduler()

	now := time.Now()
	timeToSchedule := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()+1, 0, 0, time.Local)

	job := callTodaysWeekday(scheduler.Every(1)).At(fmt.Sprintf("%02d:%02d", timeToSchedule.Hour(), timeToSchedule.Minute()))
	job.Do(task)
	t.Logf("job is scheduled for %s", job.NextScheduledTime())
	if !job.NextScheduledTime().Equal(timeToSchedule) {
		t.Fail()
		t.Logf("Job should be run today, at the set time.")
	}
}

func Test_formatTime(t *testing.T) {
	tests := []struct {
		name     string
		args     string
		wantHour int
		wantMin  int
		wantErr  bool
	}{
		{
			name:     "normal",
			args:     "16:18",
			wantHour: 16,
			wantMin:  18,
			wantErr:  false,
		},
		{
			name:     "normal",
			args:     "6:18",
			wantHour: 6,
			wantMin:  18,
			wantErr:  false,
		},
		{
			name:     "notnumber",
			args:     "e:18",
			wantHour: 0,
			wantMin:  0,
			wantErr:  true,
		},
		{
			name:     "outofrange",
			args:     "25:18",
			wantHour: 25,
			wantMin:  18,
			wantErr:  true,
		},
		{
			name:     "wrongformat",
			args:     "19:18:17",
			wantHour: 0,
			wantMin:  0,
			wantErr:  true,
		},
		{
			name:     "wrongminute",
			args:     "19:1e",
			wantHour: 19,
			wantMin:  0,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHour, gotMin, err := formatTime(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("formatTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHour != tt.wantHour {
				t.Errorf("formatTime() gotHour = %v, want %v", gotHour, tt.wantHour)
			}
			if gotMin != tt.wantMin {
				t.Errorf("formatTime() gotMin = %v, want %v", gotMin, tt.wantMin)
			}
		})
	}
}

// utility function for testing the weekday functions *on* the current weekday.
func callTodaysWeekday(job *Job) *Job {
	switch time.Now().Weekday() {
	case 0: job.Sunday()
	case 1: job.Monday()
	case 2: job.Tuesday()
	case 3: job.Wednesday()
	case 4: job.Thursday()
	case 5: job.Friday()
	case 6: job.Saturday()
	}
	return job
}