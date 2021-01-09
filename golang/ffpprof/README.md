# ffpprof
ffpprof is a easy-to-setup pprof profile data collection library for development time.
it highly depends on the standard library packages such as runtime/pprof and the exiting
tools such as go tool pprof.

## usage
```cassandraql
import ffpprof

autopprof.Capture(autopprof.CPUProfile{
    Duration: 30 * time.Second,
})
```

