## v2h
convert value to hex

## usage
```cassandraql
Format(v uint64) string
Parse(s string) (uint64, error)
```

## eg
```cassandraql
DECIMAL   HEX   VEX    ZERO-PADDED
0         0     00     0000000000000000
1         1     01     0000000000000001
2         2     02     0000000000000002
10        a     0a     000000000000000a
15        f     0f     000000000000000f
16        10    110    0000000000000010
288       120   2120   0000000000000120
```

## inspired
https://github.com/benbjohnson/vex