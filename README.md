# docker-redis-log-driver

Redis log driver for Docker. This is heavily inspired by https://github.com/cpuguy83/docker-log-driver-test.

## Features

* Send containers stdout/stderr to a Redis list
* Log lines come as JSON messages with lots of meta data
  * Container ID (`container_id`)
  * Container name (`container_name`)
  * Container creation date (`container_created`)
  * Image ID (`image_id`)
  * Image name (`image_name`)
  * Command including `ENTRYPOINT` and arguments (`command`)
  * Log tag as provided via `--log-tag` option (`tag`)
  * Extra information as defined via `--log-opt labels=` or `--log-opt env=` (`extra`)
  * Host that container runs on (`host`)
  * Timestamp when log was generated (`timestamp`)
* Configure logging setup either globally through Docker `config.json` or per container (`--log-opt` style)
* Customizable Redis connection timeouts
* Automatic Redis reconnects

## Install

Requires at least Docker 17.05 since that brings support for log driver plugins. Doesn't work on Windows currently because
log driver plugins aren't supported by Docker for Windows.

```
git clone https://github.com/pressrelations/docker-redis-log-driver
cd docker-redis-log-driver
docker build -t docker-redis-log-driver .
ID=$(docker create docker-redis-log-driver true)
rm -rf rootfs
mkdir rootfs
docker export $ID | tar -x -C rootfs/
docker plugin create docker-redis-log-driver .
docker plugin enable docker-redis-log-driver
```

## Usage

The basic usage looks like this:

```
docker run --log-driver docker-redis-log-driver --log-opt redis-address=redis.domain.com:6379 --log-opt redis-password=secure --log-opt redis-list=logs alpine date
```

Observe the Redis list named `logs` in database 0 of your Redis instance. You should see a JSON like the following:

```
{
  "message": "Thu Jun  8 13:40:30 UTC 2017\\r",
  "container_id": "6ec9a6890823800a22db411a80c15f6a0642c7983832cb3138679d68c019ea47",
  "container_name": "musing_liskov",
  "container_created": "2017-06-08T14:08:45.196778472Z",
  "image_id": "sha256:baa5d63471ead618ff91ddfacf1e2c81bf0612bfeb1daf00eb0843a41fbfade3",
  "image_name": "alpine",
  "command": "sh date",
  "tag": "musing_liskov",
  "host": "workstation",
  "timestamp": "2017-06-08T14:09:19.485778425Z"
}
```

## Uninstall

To uninstall, please make sure that no containers are still using this plugin. After that, disable and remove the plugin like this:

```
docker plugin disable docker-redis-log-driver
docker plugin rm docker-redis-log-driver
```
