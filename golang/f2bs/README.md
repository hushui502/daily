# f2bs
file to byte slice

# usage
```cassandraql
Usage of file2byteslice:
  -buildtags string
        build tags
  -compress
        use gzip compression
  -input string
        input filename
  -output string
        output filename
  -package string
        package name (default "main")
  -var string
        variable name (default "_")
```

# use
```cassandraql
f2bs -input INPUT_FILE_NAME -output OUTPUT_FILE_NAME -package PACKAGE_NAME -var VARIABLE_NAME
```