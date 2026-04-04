# Code Blocks

## Python

```python
def fibonacci(n: int) -> int:
    if n <= 1:
        return n
    return fibonacci(n - 1) + fibonacci(n - 2)

for i in range(10):
    print(f"fib({i}) = {fibonacci(i)}")
```

## Go

```go
package main

import "fmt"

func main() {
    ch := make(chan int, 10)
    go func() {
        for i := 0; i < 10; i++ {
            ch <- i * i
        }
        close(ch)
    }()
    for v := range ch {
        fmt.Println(v)
    }
}
```

## JavaScript

```javascript
const fetchData = async (url) => {
  const res = await fetch(url);
  return res.json();
};
```

## YAML

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  replicas: 3
```

## Bash

```bash
#!/bin/bash
for file in *.md; do
  unimd2pdf -i "$file"
done
```

## SQL

```sql
SELECT u.name, COUNT(o.id) AS orders
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
GROUP BY u.name
ORDER BY orders DESC;
```
