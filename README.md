## Configuration

### Environment variables

* LOG_LEVEL=The logger level.
* LOG_INSTANCE=The instance name.
* LOG_FILE_PATH=The path to log file. Will be created inside a user cache directory. Logging into the file will be disabled If it is empty.

### Closing the logger

Since you can write to file you should also close the logger.
```golang
Logger.Close()
```

### Sending a log data to the controller

To send the log data to controller you need to set up a sender before the logger close.
Otherwise, it won't be sent.

```golang
Logger.SetSender(accessToken, url, machinePairID)
```