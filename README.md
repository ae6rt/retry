Go retry package with timeout and retry limits.

Usage

```
...

import "github.com/ae6rt/retry"

work := func() error {
   // do stuff
   return nil
}

try := retry.New(3*time.Second, 3, retry.DefaultBackoffFunc)

err := try.Try(work)

if err != nil {
   if retry.IsTimeout(err) {
     fmt.Printf("Timeout\n")
   } else {
     fmt.Printf("Error: %v\n", err)
   }
}
```
