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

## Installation

Requires at least Docker 17.05 since that brings support for log driver plugins.
