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
            alloy["Grafana Alloy"]
        end

        subgraph beyla_ns["Beyla"]
            beyla["Grafana Beyla (eBPF Auto-instrumentation)"]
        end

        subgraph kps["kube-prometheus-stack"]
            grafana["Grafana"]
            kube_sm["ServiceMonitors / PodMonitors / PrometheusRules"]
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
    beyla t1@-->|"OTLP metrics+traces (:4318)"| alloy

    %% App traces via OTel SDK
    apps t2@-->|"OTLP gRPC/HTTP (:4317/:4318)"| alloy

    %% Alloy -> metrics
    alloy m1@-->|"prometheus.remote_write (RW v2)"| mimir
    alloy m2@-->|"mimir.rules.kubernetes"| mimir_ruler

    %% Alloy scrape via CRDs
    kube_sm m3@-->|"ServiceMonitors / PodMonitors"| alloy
    exporters m4@-->|"scrape"| alloy

    %% Alloy -> logs
    alloy l1@-->|"pod logs (file tail)"| loki
    alloy l2@-->|"k8s events"| loki
    alloy l3@-->|"OTLP logs"| loki

    %% Alloy -> traces
    alloy t3@-->|"OTLP traces"| tempo

    %% Tempo MetricsGenerator -> Mimir
    tempo_mg m5@-->|"prometheus.remote_write (span metrics)"| mimir

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
    traefik i2@-->|"alloy.k8s.localhost"| alloy
    traefik i3@-->|"minio.k8s.localhost"| minio
    traefik i4@-->|"httpbin.k8s.localhost"| httpbin

    classDef traces stroke:#9b59b6
    classDef metrics stroke:#f5a623
    classDef logs stroke:#7ed321
    classDef storage stroke:#95a5a6
    classDef read stroke:#85c1e9
    classDef ingress stroke:#3498db

    class t1,t2,t3 traces
    class m1,m2,m3,m4,m5 metrics
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
