# slog-gorm

## Key features

- compatible with any `slog.Handler`, which allows you to keep control on
  the format of your logs.
- can define a threshold to identify and log the slow queries.
- can log all SQL messages or just the errors if you prefer.
- can define a custom `slog.Level` for errors, slow queries or the other logs.


## Requirement

- `golang >= 1.21`

## Usage



```golang
import (
    "log/slog"
    "os"

    "gorm.io/driver/sqlite"
    "gorm.io/gorm"

    slogGorm "github.com/apus-run/sea-kit/gormx/log"
)

// Create an slog-gorm instance
gormLogger := slogGorm.New() // use slog.Default() by default


// GORM: Globally mode
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
    Logger: gormLogger,
})

// GORM: Continuous session mode
tx := db.Session(&Session{Logger: gormLogger})
tx.First(&user)
tx.Model(&user).Update("Age", 18)
```

```golang
   import (
      "gorm.io/gorm"
      "gorm.io/gorm/clause"
      "gorm.io/gorm/logger"

       slogGorm "github.com/apus-run/sea-kit/gormx/log"
    )
	    
   var gormLogger logger.Interface
	if gormTraceAll {
		gormLogger = slogGorm.New(slogGorm.WithLogger(slog.Default()), slogGorm.WithTraceAll())
	} else {
		gormLogger = slogGorm.New(slogGorm.WithLogger(slog.Default()))
	}

	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: false,
		SkipDefaultTransaction:                   true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to open database")
	}
```

### With your `slog.Logger`

The following example shows you how to use a specific `slog.Logger` with `slog-gorm`:

```golang
// With your slog.Logger
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

// Also, you can set specific attributes to distinguish between your application logs and gorm logs
// logger = logger.With(slog.String("log_type", "database"))

gormLogger := slogGorm.New(
    slogGorm.WithLogger(logger), // Optional, use slog.Default() by default
    slogGorm.WithTraceAll(), // trace all messages 
    slogGorm.SetLogLevel(DefaultLogType, slog.Level(32)), // Define the default logging level
)
```

### Use your custom `slog.Level`

As some loggers *(e.g. syslog)* have their own logging levels, `slog-gorm` lets you
use them to ensure the consistency of your logs and make them easier to understand.

You can set the logging level for these log types:

| Type                        | Description                          | Default           |
|-----------------------------|--------------------------------------|-------------------|
| `slogGorm.ErrorLogType`     | For SQL errors                       | `slog.LevelError` |
| `slogGorm.SlowQueryLogType` | For slow queries                     | `slog.LevelWarn`  |
| `slogGorm.DefaultLogType`   | For other messages *(default level)* | `slog.LevelInfo`  |

Example:

```golang
const (
    LOG_EMERG   = slog.Level(0)
    // ...
    LOG_ERR     = slog.Level(3)
    LOG_WARNING = slog.Level(4)
    LOG_NOTICE  = slog.Level(5)
    // ...
    LOG_DEBUG   = slog.Level(7)
)

logger := slog.New(syslogHandler)

gormLogger := slogGorm.New(
    slogGorm.WithLogger(logger),

    // Set logging level for SQL errors
    slogGorm.SetLogLevel(slogGorm.ErrorLogType, LOG_ERR)

    // Set logging level for slow queries
    slogGorm.SetLogLevel(slogGorm.SlowQueryLogType, LOG_NOTICE)

    // Set logging level for other messages (default level)
    slogGorm.SetLogLevel(slogGorm.DefaultLogType, LOG_DEBUG)
)
```

### Other options

```golang
customLogger := sloggorm.New(
	slogGorm.WithSlowThreshold(500 * time.Millisecond), // to identify slow queries

	slogGorm.WithRecordNotFoundError(), // don't ignore not found errors

	slogGorm.WithSourceField("origin"), // instead of "file" (by default)

	slogGorm.WithErrorField("err"),     // instead of "error" (by default)
)
```

By default, the slow queries and SQL errors are logged, but you can ignore all SQL messages with `WithIgnoreTrace()`.

```
customLogger := sloggorm.New(
    slogGorm.WithIgnoreTrace(), // disable the tracing of SQL queries by the logger.
)
```
