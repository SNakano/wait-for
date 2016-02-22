# wait-for
Wait for connect to a server at command-line

# Basic usage

```
$ ./wait-for -h
Usage: ./wait-for [-w host:port...] <command>
Version: 0.0.1
  -t int
        maximum time allowed for connection (default 300)
  -w value
        wait for host[s]. format: host:port (default [])
  <command>
        execute command
```

# Usage with docker-compose

Dockerfile
```Dockerfile
FROM ubuntu:latest
ADD https://github.com/SNakano/wait-for/releases/download/v0.0.1/wait-for_linux_amd64 /wait-for
RUN chmod o+x /wait-for
```

docker-compose.yml
```yaml
hello-world:
  build: .
  links:
    - els
    - mq
  entrypoint: /wait-for -w els:9200 -w mq:5672
  command: echo "hello world"

els:
  image: elasticsearch

mq:
  image: rabbitmq
```
