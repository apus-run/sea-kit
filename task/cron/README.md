## cron

Scheduled task library encapsulated on [cron v3](github.com/robfig/cron).

### Cron 时间表语法

```
# ┌───────────── 分钟 (0 - 59)
# │ ┌───────────── 小时 (0 - 23)
# │ │ ┌───────────── 月的某天 (1 - 31)
# │ │ │ ┌───────────── 月份 (1 - 12)
# │ │ │ │ ┌───────────── 周的某天 (0 - 6)（周日到周一；在某些系统上，7 也是星期日）
# │ │ │ │ │
# │ │ │ │ │
# │ │ │ │ │
# * * * * *
```

```
# ┌───────────── 分钟 (0 - 59)
# │ ┌───────────── 小时 (0 - 23)
# │ │ ┌───────────── 月的某天 (1 - 31)
# │ │ │ ┌───────────── 月份 (1 - 12)
# │ │ │ │ ┌───────────── 周的某天 (0 - 6) （周日到周一；在某些系统上，7 也是星期日）
# │ │ │ │ │
# │ │ │ │ │
# │ │ │ │ │
# * * * * *
```

<!-- 
| Entry 	| Description   | Equivalent to |
| ------------- | ------------- |-------------  |
| @yearly (or @annually) | Run once a year at midnight of 1 January | 0 0 1 1 * |
| @monthly               | Run once a month at midnight of the first day of the month | 0 0 1 * * |
| @weekly                | Run once a week at midnight on Sunday morning | 0 0 * * 0 |
| @daily (or @midnight)  | Run once a day at midnight | 0 0 * * * |
| @hourly                | Run once an hour at the beginning of the hour | 0 * * * * |
-->
| 输入                      | 描述                          | 相当于         |
| -------------             | -------------                 |-------------   |
| @yearly (or @annually)    | 每年 1 月 1 日的午夜运行一次  | 0 0 1 1 *      |
| @monthly                  | 每月第一天的午夜运行一次      | 0 0 1 * *      |
| @weekly                   | 每周的周日午夜运行一次        | 0 0 * * 0      |
| @daily (or @midnight)     | 每天午夜运行一次              | 0 0 * * *      |
| @hourly                   | 每小时的开始一次              | 0 * * * *      |

<!--  
For example, the line below states that the task must be started every Friday at midnight, as well as on the 13th of each month at midnight:
-->
例如，下面这行指出必须在每个星期五的午夜以及每个月 13 号的午夜开始任务：

`0 0 13 * 5`

<!--  
To generate CronJob schedule expressions, you can also use web tools like [crontab.guru](https://crontab.guru/).
-->
要生成 CronJob 时间表表达式，你还可以使用 [crontab.guru](https://crontab.guru/) 之类的 Web 工具。

### Example of use

```go
package main

import (
    "fmt"
    "time"
	"github.com/apus-run/sea-kit/task/cron"
)

var task1 = func() {
     fmt.Println("this is task1")
     fmt.Println("running task list:", cron.GetRunningTasks())
}

var taskOnce = func() {
	taskName := "taskOnce"
    fmt.Println("this is taskOnce")
    cron.DeleteTask(taskName)
}

func main() {
    err := cron.Init()
    if err != nil {
        panic(err)
    }
	
    cron.Run([]*cron.Task{
        {
            Name:     "task1",
            TimeSpec: "@every 2s",
            Fn:       task1,
        },
        {
            Name:     "taskOnce",
            TimeSpec: "@every 5s",
            Fn:       taskOnce,
        },
    }...)

    time.Sleep(time.Minute)

    // delete task1
    cron.DeleteTask("task1")

    // view running tasks
    fmt.Println("running task list:", cron.GetRunningTasks())
}
```