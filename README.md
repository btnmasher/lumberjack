# lumberjack

Simple/Stupid Multi-Backend Logger for Go

Named Lumberjack for obvious punny reasons. (Get it.. Logger... Lumberjack... GET IT?!)

![Deal With It](http://i.imgur.com/o3bP95Y.png)

## Why

Because I like my logs formatted a certain way, and was tired of duplicating my efforts to setup the boilerplate to log in a particular way. So I made this project for my own purposes to have a quick and easy way to set up logging with levels, formatting, runtime information, and multiple backends.

## State of the Project

It is currently in an early state with very basic console printing backend. It's not intended for use in production, but I probably will use it anyway because I'm a masochist.

## What's Planned

##### More Backends!

- BSON
- SMTP
- Syslog
- Who knows?

##### Custom Formatting!

- Using Go's Templating Engine
- COLOR!! (Who doesn't like color?)

## How to use this garbage?

```Go
    import "github.com/btnmasher/lumberjack"

    ...

    //Setup the logger instance
    logger = lumberjack.NewLoggerWithDefaults()

    //Add a level (defaults included: INFO, WARN, ERROR, CRITICAL, FATAL)
    logger.AddLevel(lumberjack.DEBUG)

    logger.Info("Holy Shit!")
    logger.Debugf("%s did a thing.", "thing")
```

With the defaults, the verbosity level of the print backend is set to print extra information (the file and line number) on ERROR or higher

The output looks like this:

```Bash
    2015/08/17 12:23:57 (INFO) @ main.main(): Holy Shit!
    2015/08/17 12:23:57 (DEBUG) @ main.main() main.go:10: thing did a thing.
```

In the future, to add Backends, they simply need to implement the interface:

`lumberjack.Backend`

with the requiremed method:

`Log(*LogEntry)`

then add to the logger:

```Go
    /* Specifying the name allows you to have
    multiple copies of the same backend with
    different settings */
    logger.AddBackend("somename", &SomeBackend{})
````

##### Http Backend?

You can specify a basic HTTP backend to POST log entries formatted in JSON. Implemented in the backend, is a buffering mechanism to buffer an arbitrarily defined number of log entries for a given time period. The buffer will be sent via HTTP POST as a JSON array either when the buffer is filled, or the time interval elapses (whichever occurs first).

```Go
    /* Add an HTTP JSON POST backend with a message
    buffer of 10 per 5sec */
    hb := lumberjack.NewHttpClientBackend(
    "http://www.example.com:1234", 10, time.Second*5)

    //Don't forget to add the backend!
    logger.AddBackend("http", hb)

    //Defer the closing of the HTTP Backend's goroutine.
    defer close(hb.Stop)
```

So given the above example, once 10 log entries are sent to the backend, it will HTTP POST them to the specified URL. Or, if 5 seconds elapses, whatever is currently in the buffer will be sent without waiting to fill.

## Want to Contribute?

Send me a pull request, I'll probably merge it. But let's be honest, who's going to use this drivel? :P