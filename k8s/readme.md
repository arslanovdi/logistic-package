Для работы приложения в k8s кластере должны быть установлены:
- [longhorn](https://longhorn.io/)
- [istio ingress/egress](https://istio.io/latest/docs/setup/getting-started/)
- [strimzi cluster operator](https://strimzi.io/)
- [zalando postgres-operator](https://github.com/zalando/postgres-operator)

strimzi должен быть установлен в namespace `kafka`