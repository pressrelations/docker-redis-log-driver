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
mkdir rootfs
docker export $ID | tar -x -C rootfs/
docker plugin create docker-redis-log-driver .
docker plugin enable docker-redis-log-driver
```

## Usage

### Basic usage

Run a container using this plugin:

```
docker run --log-driver docker-redis-log-driver --log-opt redis-address=some.redis.server:6379 --log-opt redis-password=secure --log-opt redis-list=logs alpine date
```

Observe the Redis list named `logs` in database `0` of your Redis instance. You should see a JSON like the following:

```
{
  "message": "Thu Jun  8 14:36:01 UTC 2017",
  "container_id": "20cd880a54679f26c85edc53a7fbe7079d7e5b88a883f245a0c215d0afd3e600",
  "container_name": "hungry_cori",
  "container_created": "2017-06-08T14:36:01.08719318Z",
  "image_id": "sha256:baa5d63471ead618ff91ddfacf1e2c81bf0612bfeb1daf00eb0843a41fbfade3",
  "image_name": "alpine",
  "command": "date",
  "tag": "20cd880a5467",
  "extra": {},
  "host": "hostname",
  "timestamp": "2017-06-08T14:36:01.64922286Z"
}
```

### Advanced usage

This example shows the usage of

* Custom log tags (c.f. https://docs.docker.com/engine/admin/logging/log_tags/)
* Container label logging
* Container environment variable logging

```
docker run --label foo=abc --label bar=xyz -e SOME_ENV_VAR=foobar --log-driver docker-redis-log-driver --log-opt redis-address=some.redis.server:6379 --log-opt redis-password=secure --log-opt redis-list=logs --log-opt "tag={{.ImageName}}/{{.Name}}/{{.ID}}" --log-opt labels=foo,bar --log-opt env=SOME_ENV_VAR alpine date
```

Observe the Redis list named `logs` in database `0` of your Redis instance. You should see a JSON like the following:

```
{
  "message": "Thu Jun  8 14:45:22 UTC 2017",
  "container_id": "35fb6802c95dce77147b89f595ed7528ebf364c8a795ef1468f52464bfd8fb23",
  "container_name": "trusting_tesla",
  "container_created": "2017-06-08T14:45:22.253419923Z",
  "image_id": "sha256:baa5d63471ead618ff91ddfacf1e2c81bf0612bfeb1daf00eb0843a41fbfade3",
  "image_name": "alpine",
  "command": "date",
  "tag": "alpine/trusting_tesla/35fb6802c95d",
  "extra": {
    "SOME_ENV_VAR": "foobar",
    "bar": "xyz",
    "foo": "abc"
  },
  "host": "workstation",
  "timestamp": "2017-06-08T14:45:22.625446437Z"
}
```

## Uninstall

To uninstall, please make sure that no containers are still using this plugin. After that, disable and remove the plugin like this:

```
docker plugin disable docker-redis-log-driver
docker plugin rm docker-redis-log-driver
```
