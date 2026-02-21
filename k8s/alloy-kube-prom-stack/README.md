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
        traefik["Traefik (NodePort:30000)"]
    end

    subgraph apps["アプリケーション"]
        httpbin["httpbin"]
        mythical["mythical"]
    end

    subgraph observability["オブザーバビリティ基盤"]
        subgraph alloy_ns["Alloy (DaemonSet + Clustering)"]
            alloy_metrics["Metrics Collection\n(ServiceMonitors / PodMonitors / self)"]
            alloy_logs["Log Collection\n(pod logs / k8s events)"]
            alloy_otel["OTel Pipeline\n(receiver → processors → exporters)"]
            alloy_rules["mimir.rules.kubernetes"]
        end

        subgraph beyla_ns["Beyla"]
            beyla["Grafana Beyla (eBPF Auto-instrumentation)"]
        end

        subgraph grafana_ns["Grafana (kube-prometheus-stack)"]
            grafana["Grafana"]
        end

        subgraph kps["kube-prometheus-stack"]
            kube_monitors["ServiceMonitors / PodMonitors"]
            kube_rules["PrometheusRules"]
            exporters["Node/KSM/etc Exporters"]
        end

        subgraph mimir_ns["Mimir"]
            mimir["Mimir (metrics store)"]
            mimir_am["Alertmanager"]
            mimir_ruler["Ruler"]
        end

        subgraph loki_ns["Loki (SimpleScalable)"]
            loki["Loki (log store)"]
        end

        subgraph tempo_ns["Tempo"]
            tempo["Tempo (trace store)"]
            tempo_mg["MetricsGenerator (span-metrics, service-graphs)"]
        end

        subgraph minio_ns["MinIO"]
            minio["MinIO (S3 backend)"]
        end
    end

    %% eBPF instrumentation
    beyla t1@-->|"OTLP (:4318)"| alloy_otel

    %% App traces via OTel SDK
    apps t2@-->|"OTLP (:4317/:4318)"| alloy_otel

    %% Alloy -> metrics
    alloy_metrics m1@-->|"remote_write"| mimir
    alloy_rules m2@-->|"rules sync"| mimir_ruler

    %% Alloy scrape via CRDs
    kube_monitors m3@-->|"ServiceMonitors / PodMonitors"| alloy_metrics
    exporters m4@-->|"scrape"| alloy_metrics
    kube_rules m7@-->|"PrometheusRules"| alloy_rules

    %% Alloy -> logs
    alloy_logs l1@-->|"pod logs (file tail)"| loki
    alloy_logs l2@-->|"k8s events"| loki
    alloy_otel l3@-->|"OTLP logs"| loki

    %% Alloy -> traces
    alloy_otel t3@-->|"OTLP traces"| tempo

    %% Alloy OTel -> metrics
    alloy_otel m6@-->|"OTLP metrics"| mimir

    %% Tempo MetricsGenerator -> Mimir
    tempo_mg m5@-->|"remote_write (span metrics)"| mimir

    %% Storage backends
    mimir s1@-->|"S3 (tsdb / ruler / alertmanager)"| minio
    loki  s2@-->|"S3 (chunks / ruler)"| minio
    tempo s3@-->|"S3 (traces)"| minio

    %% Grafana datasources
    grafana r1@-->|"PromQL"| mimir
    grafana r2@-->|"LogQL"| loki
    grafana r3@-->|"TraceQL"| tempo
    grafana r4@-->|"Alertmanager"| mimir_am

    %% Ingress
    traefik i1@-->|"grafana.k8s.localhost"| grafana
    traefik i2@-->|"alloy.k8s.localhost"| alloy_otel
    traefik i3@-->|"minio.k8s.localhost"| minio
    traefik i4@-->|"httpbin.k8s.localhost"| httpbin

    classDef traces stroke:#9b59b6
    classDef metrics stroke:#f5a623
    classDef logs stroke:#7ed321
    classDef storage stroke:#95a5a6
    classDef read stroke:#85c1e9
    classDef ingress stroke:#3498db

    class t1,t2,t3 traces
    class m1,m2,m3,m4,m5,m6,m7 metrics
    class l1,l2,l3 logs
    class s1,s2,s3 storage
    class r1,r2,r3,r4 read
    class i1,i2,i3,i4 ingress
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
