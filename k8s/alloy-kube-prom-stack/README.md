# alloy + kube-prom-stack

https://grafana.com/docs/alloy/latest/

以下をやってみる
- [x] alloy
  - [x] self monitoring
  - [x] prometheus.operator.podmonitors
  - [x] prometheus.operator.servicemonitors
  - [x] mimir.rules.kubernetes
  - [x] cluster label
  - [x] k8s events
  - [x] pod logs
  - [x] otel traces, metrics, logs
- kube-prom-stack
  - [x] without prometheus or alertmanager
  - [x] crds
  - [x] exporters & service monitors
  - [x] rules
  - [x] grafana
  - [x] dashboards
- [x] mimir
  - [x] ruler
  - [x] alertmanager
- [x] loki
- [x] tempo
- [x] beyla


## 構成図

```mermaid
graph TB
    subgraph ingress["Ingress"]
        traefik["Traefik\n(NodePort:30000)"]
    end

    subgraph apps["アプリケーション"]
        httpbin["httpbin"]
        mythical["mythical"]
    end

    subgraph observability["オブザーバビリティ基盤"]
        subgraph alloy_ns["Alloy (DaemonSet + Clustering)"]
            alloy["Grafana Alloy"]
        end

        subgraph beyla_ns["Beyla"]
            beyla["Grafana Beyla\n(eBPF Auto-instrumentation)"]
        end

        subgraph kps["kube-prometheus-stack"]
            grafana["Grafana"]
            kube_sm["ServiceMonitors / PodMonitors / PrometheusRules"]
            exporters["Node/KSM/etc Exporters"]
        end

        subgraph mimir_ns["Mimir"]
            mimir["Mimir\n(metrics store)"]
            mimir_am["Alertmanager"]
            mimir_ruler["Ruler"]
        end

        subgraph loki_ns["Loki (SimpleScalable)"]
            loki["Loki\n(log store)"]
        end

        subgraph tempo_ns["Tempo"]
            tempo["Tempo\n(trace store)"]
            tempo_mg["MetricsGenerator\n(span-metrics, service-graphs)"]
        end

        subgraph minio_ns["MinIO"]
            minio["MinIO\n(S3 backend)"]
        end
    end

    %% eBPF instrumentation
    beyla -->|"OTLP metrics+traces\n:4318"| alloy

    %% App traces via OTel SDK
    apps -->|"OTLP gRPC/HTTP\n:4317/:4318"| alloy

    %% Alloy -> metrics
    alloy -->|"prometheus.remote_write\n(Prometheus RW v2)"| mimir
    alloy -->|"mimir.rules.kubernetes\n(sync PrometheusRules)"| mimir_ruler

    %% Alloy scrape via CRDs
    kube_sm -->|"ServiceMonitors / PodMonitors"| alloy
    exporters -->|"scrape"| alloy

    %% Alloy -> logs
    alloy -->|"pod logs\n(file tail)"| loki
    alloy -->|"k8s events"| loki
    alloy -->|"otelcol.exporter.loki\n(OTLP logs)"| loki

    %% Alloy -> traces
    alloy -->|"otelcol.exporter.otlphttp\n(OTLP traces)"| tempo

    %% Tempo MetricsGenerator -> Mimir
    tempo_mg -->|"prometheus.remote_write\n(span metrics)"| mimir

    %% Storage backends
    mimir -->|"S3 (tsdb / ruler / alertmanager)"| minio
    loki  -->|"S3 (chunks / ruler)"| minio
    tempo -->|"S3 (traces)"| minio

    %% Grafana datasources
    grafana -->|"PromQL"| mimir
    grafana -->|"LogQL"| loki
    grafana -->|"TraceQL"| tempo
    grafana -->|"Alertmanager"| mimir_am

    %% Ingress
    traefik -->|"grafana.k8s.localhost"| grafana
    traefik -->|"alloy.k8s.localhost"| alloy
    traefik -->|"minio.k8s.localhost"| minio
    traefik -->|"httpbin.k8s.localhost"| httpbin
```

## TODO
- readme
  - 構成図的なやつを claude に書かせる
以下は issue にして、いったんこれは終わりにする
- otel
  - beyla の metric は scrape と otel のどっちがよい？
  - k8s 系の属性が metric についてこない
    - scrape するなら `k8s_*` が付いてくるけど
      - `k8s_` prefix を外したやつがあった方が他と統一できてよいかも
    - otel metric の場合はどう扱う？
  - metric と trace を紐付けたい
    - 何かあったはず
- ingress
  - traefik.enabled を追加したい
  - host 名を変えられるようにしたい
    - helmfile の env にして、helmfile で set すればよさそうに思う
- mimir.rules.kubernetes
  - external label: cluster
