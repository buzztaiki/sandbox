# alloy + kube-prom-stack

https://grafana.com/docs/alloy/latest/

以下をやってみる
- nodes: 2
- alloy
  - prometheus.exporter.self
  - prometheus.operator.podmonitors
  - prometheus.operator.servicemonitors
- kube-prom-stack
  - without prometheus or alertmanager
  - crds
  - exporters & service monitors
  - rules
  - grafana
  - dashboards
- mimir
- loki
