Go retry package with timeout and retry limits.

Usage

```
...

import "github.com/ae6rt/retry"

work := func() error {
   // do stuff
   return nil
}

r := retry.New(3, retry.DefaultBackoffFunc)

err := r.Try(work)

if err != nil {
   fmt.Printf("Error: %v\n", err)
}
```
