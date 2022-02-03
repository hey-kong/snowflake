# snowflake
A simple Twitter snowflake generator implemented in Golang.

snowflake format:
```
//
// |   symbol   |  timestamp  |  worker id  |     seq     |
// |------------|-------------|-------------|-------------|
// |<-- 1bit -->|<-- 41bit -->|<-- 10bit -->|<-- 12bit -->|
//
```

* The ID as a whole is an int64.
* 1 bit are used to ensure that the generated ID is a positive number.
* 41 bits are used to store a timestamp with millisecond precision, using a custom epoch.
* 10 bits are used to store a node id - a range from 0 through 1023.
* 12 bits are used to store a sequence number - a range from 0 through 4095.
