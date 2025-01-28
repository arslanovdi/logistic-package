Для работы приложения в k8s кластере должны быть установлены:
- [longhorn](https://longhorn.io/)
- [istio ingress/egress](https://istio.io/latest/docs/setup/getting-started/)
- [strimzi cluster operator](https://strimzi.io/)
- [zalando postgres-operator](https://github.com/zalando/postgres-operator)
- указать в файле `argocdApplication.yaml` адрес k8s кластера

strimzi должен быть установлен в namespace `kafka`

fluentbit должен быть установлен в namespace `observability`

После развертывания нужно перезапустить поды fluentbit и добавть gelf input в graylog порт 12201 tcp.