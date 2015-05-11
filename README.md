Go retry package with timeout and retry limits.

Usage

```
r := New(3*time.Second, 3)
f := func() error {
   return nil
}

err := r.Try(f)
if err != nil {
   if retry.isTimeout(err) {
     fmt.Printf("Timeout\n")
   } else {
     fmt.Printf("Error: %v\n", err)
   }
}
```