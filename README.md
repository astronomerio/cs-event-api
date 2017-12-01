# cs-event-api

[![CircleCI](https://circleci.com/gh/astronomerio/cs-event-api.svg?style=svg&circle-token=04ac38b355e4a32d61a0b286f7031adf7dab2c11)](https://circleci.com/gh/astronomerio/cs-event-api)

## Configuraton
The binary looks for a config file called `config.yml` in the current working directory.

An example of what that should look is:
```yaml
IngestionHandler: kafka
PrometheusEnabled: true
HealthCheckEnabled: true
PProfEnabled: true
StreamName: main

KafkaBrokers: 192.168.99.100:9092
KafkaTopic: main
```