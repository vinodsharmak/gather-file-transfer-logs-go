
### Environment variables

To configure the logger use the environment variables:

* LOG_LEVEL=The logger level.

* LOG_INSTANCE=The instance name.

* LOG_FILE_PATH=The path to log file. Will be created inside a user cache directory. Logging into the file will be disabled If it is empty.

### Setter Methods

To configure the logger use the setter methods as well:

* `Logger.SetLogFile(logFilePath string)`

> **logFilePath** is the path to log file. Will be created inside a user
> cache directory. Logging into the file will be disabled If it is
> empty.

* `Logger.SetInstance(instance string)`

> **instance** is the instance name.

* `Logger.SetLevel(level string)`

> **level** is the logger level. One from (debug | info | warning | error).

### Closing the logger

Since you can write to file you should also close the logger.

```golang

Logger.Close()

```

### Sending a log data to the controller

```golang

Logger.SendLogsToController() error

```


To send the log data to the controller you need to set up a sender before closing the logger.

Otherwise, it won't be sent.

```golang

Logger.SetSender(accessToken, url, machinePairID, machineID)

```

Case worker_node request: accessToken value request must be from machine pair

Case receiver request: accessToken value request must be from machine metadata

Case sender request: accessToken value request must be from machine metadata
