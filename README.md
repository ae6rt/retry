Go retry package with timeout and retry limits.

Usage

```
retry := New(3*time.Second, 3)
work := func() error {
   return nil
}

err := retry.Try(work)
if err != nil {
   if retry.IsTimeout(err) {
     fmt.Printf("Timeout\n")
   } else {
     fmt.Printf("Error: %v\n", err)
   }
}
```
