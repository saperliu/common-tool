package cron

import (
	"fmt"
	"github.com/robfig/cron"
)

func newCronWorker(r func()) {
	c := cron.New()
	err := c.AddFunc("0 30 * * * *", r)
	if err != nil {
		fmt.Printf(" start worker error %v ", err)
	}
	c.Start()
	c.Stop() // Stop the scheduler (does not stop any jobs already running).
}

func exp() {
	c := cron.New()
	c.AddFunc("0 30 * * * *", func() { fmt.Println("Every hour on the half hour") })
	c.AddFunc("@hourly", func() { fmt.Println("Every hour") })
	c.AddFunc("@every 1h30m", func() { fmt.Println("Every hour thirty") })
	c.Start()

	// Funcs are invoked in their own goroutine, asynchronously.

	// Funcs may also be added to a running Cron
	c.AddFunc("@daily", func() { fmt.Println("Every day") })

	// Inspect the cron job entries' next and previous run times.
	//inspect(c.Entries())
	//Field name   | Mandatory? | Allowed values  | Allowed special characters
	//----------   | ---------- | --------------  | --------------------------
	//Seconds      | Yes        | 0-59            | * / , -
	//Minutes      | Yes        | 0-59            | * / , -
	//Hours        | Yes        | 0-23            | * / , -
	//Day of month | Yes        | 1-31            | * / , - ?
	//Month        | Yes        | 1-12 or JAN-DEC | * / , -
	//Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?

	c.Stop() // Stop the scheduler (does not stop any jobs already running).
}
