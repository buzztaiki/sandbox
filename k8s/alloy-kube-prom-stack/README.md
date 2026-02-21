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

### Metrics

```mermaid
graph LR
    subgraph apps["アプリケーション"]
        httpbin["httpbin"]
        mythical["mythical"]
    end

    subgraph beyla_ns["Beyla"]
        beyla["Grafana Beyla (eBPF Auto-instrumentation)"]
    end

    subgraph kps["kube-prometheus-stack"]
        kube_monitors["ServiceMonitors / PodMonitors"]
        kube_rules["PrometheusRules"]
        exporters["Node/KSM/etc Exporters"]
    end

    subgraph alloy_ns["Alloy (DaemonSet + Clustering)"]
        alloy_metrics["Metrics Collection"]
        alloy_rules["mimir.rules.kubernetes"]
        alloy_otel["OTel Pipeline"]
    end

    subgraph mimir_ns["Mimir"]
        mimir["Mimir (metrics store)"]
        mimir_am["Alertmanager"]
        mimir_ruler["Ruler"]
    end

    subgraph tempo_ns["Tempo"]
        tempo_mg["MetricsGenerator (span-metrics, service-graphs)"]
    end

    subgraph minio_ns["MinIO"]
        minio["MinIO (S3 backend)"]
    end

    subgraph grafana_ns["Grafana"]
        grafana["Grafana"]
    end

    subgraph ingress["Ingress"]
        traefik["Traefik (NodePort:30000)"]
    end

    beyla m1@-->|"OTLP (:4318)"| alloy_otel
    apps m2@-->|"OTLP (:4317/:4318)"| alloy_otel
    alloy_metrics m3@-.->|"scrape"| kube_monitors
    alloy_metrics m4@-.->|"scrape"| exporters
    alloy_rules m5@-.->|"PrometheusRules"| kube_rules
    alloy_metrics m6@-->|"remote_write"| mimir
    alloy_rules m7@-->|"rules sync"| mimir_ruler
    alloy_otel m8@-->|"OTLP metrics"| mimir
    mimir_ruler x1@-.->|"query"| mimir
    mimir_ruler x2@-->|"recording rules"| mimir
    mimir_ruler x3@-->|"alerts"| mimir_am
    tempo_mg m9@-->|"remote_write (span metrics)"| mimir
    mimir s1@-->|"S3 (tsdb / ruler / alertmanager)"| minio
    grafana r1@-.->|"PromQL"| mimir
    grafana r2@-.->|"Alertmanager"| mimir_am
    traefik i1@-->|"grafana.k8s.localhost"| grafana
    traefik i2@-->|"alloy.k8s.localhost"| alloy_otel
    traefik i3@-->|"minio.k8s.localhost"| minio

    classDef metrics stroke:#f5a623
    classDef storage stroke:#95a5a6
    classDef read stroke:#c0392b
    classDef ingress stroke:#3498db

    class m1,m2,m3,m4,m5,m6,m7,m8,m9,x1,x2,x3 metrics
    class s1 storage
    class r1,r2 read
    class i1,i2,i3 ingress
```

### Logs

```mermaid
graph LR
    subgraph apps["アプリケーション"]
        httpbin["httpbin"]
        mythical["mythical"]
    end

    subgraph alloy_ns["Alloy (DaemonSet + Clustering)"]
        alloy_logs["Log Collection (pod logs / k8s events)"]
        alloy_otel["OTel Pipeline"]
    end

    subgraph loki_ns["Loki (SimpleScalable)"]
        loki["Loki (log store)"]
    end

    subgraph minio_ns["MinIO"]
        minio["MinIO (S3 backend)"]
    end

    subgraph grafana_ns["Grafana"]
        grafana["Grafana"]
    end

    subgraph ingress["Ingress"]
        traefik["Traefik (NodePort:30000)"]
    end

    apps l1@-->|"OTLP (:4317/:4318)"| alloy_otel
    alloy_logs l2@-->|"pod logs (file tail)"| loki
    alloy_logs l3@-->|"k8s events"| loki
    alloy_otel l4@-->|"OTLP logs"| loki
    loki s1@-->|"S3 (chunks / ruler)"| minio
    grafana r1@-.->|"LogQL"| loki
    traefik i1@-->|"grafana.k8s.localhost"| grafana
    traefik i2@-->|"alloy.k8s.localhost"| alloy_otel
    traefik i3@-->|"minio.k8s.localhost"| minio

    classDef logs stroke:#27ae60
    classDef storage stroke:#95a5a6
    classDef read stroke:#c0392b
    classDef ingress stroke:#3498db

    class l1,l2,l3,l4 logs
    class s1 storage
    class r1 read
    class i1,i2,i3 ingress
```

### Traces

```mermaid
graph LR
    subgraph apps["アプリケーション"]
        httpbin["httpbin"]
        mythical["mythical"]
    end

    subgraph beyla_ns["Beyla"]
        beyla["Grafana Beyla (eBPF Auto-instrumentation)"]
    end

    subgraph alloy_ns["Alloy (DaemonSet + Clustering)"]
        alloy_otel["OTel Pipeline"]
    end

    subgraph tempo_ns["Tempo"]
        tempo["Tempo (trace store)"]
        tempo_mg["MetricsGenerator (span-metrics, service-graphs)"]
    end

    subgraph mimir_ns["Mimir"]
        mimir["Mimir (metrics store)"]
    end

    subgraph minio_ns["MinIO"]
        minio["MinIO (S3 backend)"]
    end

    subgraph grafana_ns["Grafana"]
        grafana["Grafana"]
    end

    subgraph ingress["Ingress"]
        traefik["Traefik (NodePort:30000)"]
    end

    beyla t1@-->|"OTLP (:4318)"| alloy_otel
    apps t2@-->|"OTLP (:4317/:4318)"| alloy_otel
    alloy_otel t3@-->|"OTLP traces"| tempo
    tempo_mg m1@-->|"remote_write (span metrics)"| mimir
    tempo s1@-->|"S3 (traces)"| minio
    grafana r1@-.->|"TraceQL"| tempo
    traefik i1@-->|"grafana.k8s.localhost"| grafana
    traefik i2@-->|"alloy.k8s.localhost"| alloy_otel

    classDef traces stroke:#9b59b6
    classDef metrics stroke:#f5a623
    classDef storage stroke:#95a5a6
    classDef read stroke:#c0392b
    classDef ingress stroke:#3498db

    class t1,t2,t3 traces
    class m1 metrics
    class s1 storage
    class r1 read
    class i1,i2 ingress
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
