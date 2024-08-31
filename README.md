# concurrency-cna

Docker, Kubernetes, Prometheus, Jaeger, and many other great tools are written in Go. I was curious why, and the reason might be that Go is the cloud-native language par excellence. As demonstrated in this example, running parallel and concurrent tasks doesn’t require much code or ceremony.
```bash
Hello AWS Community!
Getting https://example.com
Response size: 1256
Getting https://google.com
Response size: 20469
Getting https://github.com/jorgetovar
Response size: 196096
Time taken 3 URLs 1.876308708s (Sync)

Getting https://github.com/jorgetovar
Getting https://example.com
Getting https://google.com
Response size: 1256 for URL https://example.com
Response size: 196096 for URL https://github.com/jorgetovar
Response size: 20527 for URL https://google.com
Time taken 3 URLs 373.382042ms (Goroutines & Channels)
```bash
