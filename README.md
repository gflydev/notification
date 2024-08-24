# gFly Notification

    Copyright © 2023, gFly
    https://www.gFly.dev
    All rights reserved.

### Usage

Install
```bash
go get -u github.com/gflydev/notification@v1.0.0
```

Quick usage `main.go`
```go
import (
    _ "github.com/gflydev/cache/redis"
    "github.com/gflydev/cache"
)

// Set cache 15 days
if err = cache.Set(key, value, time.Duration(15*24*3600) * time.Second); err != nil {
    log.Errorf("Error %q", err)
}

val, err := cache.Get(key)
if err != nil {
    log.Errorf("Error %q", err)
}

// Delete 
if err = cache.Del(key); err != nil {
    log.Errorf("Error %q", err)
}
```