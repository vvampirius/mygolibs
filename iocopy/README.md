# iocopy
It's just my simple library for copy bytes from io.Reader to io.Writer with chunks, report progress, and ability to Cancel (Context).

## Example:

```go
package main

import (
   "context"
   "log"
   "os"
   "github.com/vvampirius/mygolibs/iocopy"
   "time"
)

func main() {
   devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0777)
   if err != nil { log.Fatal(err.Error()) }

   devRandom, err := os.Open(`/dev/urandom`)
   if err != nil { log.Fatal(err.Error()) }

   ctx, _ := context.WithTimeout(context.Background(), 4 * time.Second)
   c := make(chan iocopy.Report, 0)
   go iocopy.Copy(devNull, devRandom, 1024, 1 * time.Second, c, ctx)

   active := true
   for active {
      select {
      case v := <-c:
         log.Printf("%v %f", v, v.Speed(time.Second))
         if v.Error != nil { active = false }
      case <-ctx.Done():
         active = false
      }
   }
}
```

