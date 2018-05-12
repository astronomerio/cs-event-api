# event-api

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
