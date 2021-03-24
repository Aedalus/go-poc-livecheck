# Livecheck

A quick go PoC to use buffered channels as a weighted semaphore. 

Takes in a manifest file like so

```yaml
http_livechecks:
  - name: google
    endpoint: https://google.com
    expected:
      code: 200
  - name: bad-google
    endpoint: https://bad-google.com
    expected:
      code: 200
  - name: stack-overflow
    endpoint: https://stackoverflow.com
    expected:
      code: 200
```

Run it with a given concurrency

```shell script
go run . -concurrency 1
go run . -concurrency 10
```

It will make a http get to each item in the manifest, and check that the status code is as expected

```
Found 4 livechecks
Running with concurrency: 1

google
        success
bad-google
        fail: Get "https://bad-google.com": dial tcp: lookup bad-google.com: No address associated with hostname
stack-overflow
        success
```

The combination of using buffered channels as semaphores + a simple wait group allows the lib to easily control how many goroutines to launch at once.