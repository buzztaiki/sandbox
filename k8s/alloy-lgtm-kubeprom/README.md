# alloy + LGTM stack + kube-prom-stack

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

    subgraph kps["kube-prometheus-stack"]
        kube_monitors["ServiceMonitors / PodMonitors"]
        kube_rules["PrometheusRules"]
        exporters["Node/KSM/etc Exporters"]
    end

    subgraph alloy_ns["Alloy (DaemonSet + Clustering)"]
        alloy_metrics["Metrics Collection"]
        alloy_rules["mimir.rules.kubernetes"]
    end

    subgraph mimir_ns["Mimir"]
        mimir["Mimir (metrics store)"]
        mimir_am["Alertmanager"]
        mimir_ruler["Ruler"]
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

    alloy_metrics m1@-.->|"scrape"| kube_monitors
    alloy_metrics m2@-.->|"scrape"| exporters
    alloy_rules m3@-.->|"PrometheusRules"| kube_rules
    alloy_metrics m4@-->|"remote_write"| mimir
    alloy_rules m5@-->|"rules sync"| mimir_ruler
    mimir_ruler x1@-.->|"query"| mimir
    mimir_ruler x2@-->|"recording rules"| mimir
    mimir_ruler x3@-->|"alerts"| mimir_am
    mimir s1@-->|"S3 (tsdb / ruler / alertmanager)"| minio
    grafana r1@-.->|"PromQL"| mimir
    grafana r2@-.->|"Alertmanager"| mimir_am
    traefik i1@-->|"grafana.k8s.localhost"| grafana
    traefik i2@-->|"httpbin.k8s.localhost"| httpbin

    classDef metrics stroke:#f5a623
    classDef storage stroke:#95a5a6
    classDef read stroke:#c0392b
    classDef ingress stroke:#3498db

    class m1,m2,m3,m4,m5,x1,x2,x3 metrics
    class s1 storage
    class r1,r2 read
    class i1,i2 ingress
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
    traefik i2@-->|"httpbin.k8s.localhost"| httpbin

    classDef logs stroke:#27ae60
    classDef storage stroke:#95a5a6
    classDef read stroke:#c0392b
    classDef ingress stroke:#3498db

    class l1,l2,l3,l4 logs
    class s1 storage
    class r1 read
    class i1,i2 ingress
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
    tempo t4@-->|"spans"| tempo_mg
    alloy_otel m2@-->|"OTLP metrics"| mimir
    tempo_mg m1@-->|"remote_write (span metrics)"| mimir
    tempo s1@-->|"S3 (traces)"| minio
    grafana r1@-.->|"TraceQL"| tempo
    traefik i1@-->|"grafana.k8s.localhost"| grafana
    traefik i2@-->|"httpbin.k8s.localhost"| httpbin

    classDef traces stroke:#9b59b6
    classDef metrics stroke:#f5a623
    classDef storage stroke:#95a5a6
    classDef read stroke:#c0392b
    classDef ingress stroke:#3498db

    class t1,t2,t3,t4 traces
    class m1,m2 metrics
    class s1 storage
    class r1 read
    class i1,i2 ingress
```

### Alloy Pipelines

```mermaid
graph TD
    subgraph sources["Sources"]
        k8s_pods["k8s pods"]
        k8s_events["k8s events"]
        servicemonitors["ServiceMonitors"]
        podmonitors["PodMonitors"]
        self["alloy self"]
        prometheusrules["PrometheusRules"]
        apps["Apps / Beyla (OTLP)"]
    end

    subgraph alloy_ns["Alloy"]
        subgraph metrics_pipeline["Metrics Pipeline"]
            disc_k8s["discovery.kubernetes"]
            disc_relabel["discovery.relabel"]
            op_sm["prometheus.operator.servicemonitors"]
            op_pm["prometheus.operator.podmonitors"]
            prom_self["prometheus.exporter.self"]
            prom_scrape_self["prometheus.scrape / self"]
            prom_rw["prometheus.remote_write"]
            mimir_rules["mimir.rules.kubernetes"]
        end

        subgraph logs_pipeline["Logs Pipeline"]
            disc_relabel_logs["discovery.relabel / pod_logs"]
            file_match["local.file_match"]
            loki_file["loki.source.file"]
            loki_events["loki.source.kubernetes_events"]
            loki_write["loki.write"]
        end

        subgraph otel_pipeline["OTel Pipeline"]
            otel_recv["otelcol.receiver.otlp"]
            otel_mem["otelcol.processor.memory_limiter"]
            otel_sampler["otelcol.processor.probabilistic_sampler"]
            otel_resdet["otelcol.processor.resourcedetection"]
            otel_k8sattr["otelcol.processor.k8sattributes"]
            otel_batch["otelcol.processor.batch"]
            otel_transform["otelcol.processor.transform / external_labels"]
            otel_exp_prom["otelcol.exporter.prometheus"]
            otel_exp_loki["otelcol.exporter.loki"]
            otel_exp_otlp["otelcol.exporter.otlphttp"]
        end
    end

    subgraph backends["Backends"]
        mimir["Mimir"]
        loki["Loki"]
        tempo["Tempo"]
    end

    %% Metrics pipeline
    k8s_pods m1@-->| | disc_k8s
    disc_k8s m2@-->| | disc_relabel
    disc_relabel m3@-->| | disc_relabel_logs
    servicemonitors m4@-.->| | op_sm
    podmonitors m5@-.->| | op_pm
    prometheusrules m6@-.->| | mimir_rules
    op_sm m7@-->| | prom_rw
    op_pm m8@-->| | prom_rw
    self m9@-->| | prom_self
    prom_self m10@-->| | prom_scrape_self
    prom_scrape_self m11@-->| | prom_rw
    otel_exp_prom m12@-->| | prom_rw
    prom_rw m13@-->|"remote_write"| mimir
    mimir_rules m14@-->|"rules sync"| mimir

    %% Logs pipeline
    disc_relabel_logs l1@-->| | file_match
    file_match l2@-->| | loki_file
    loki_file l3@-->| | loki_write
    k8s_events l4@-->| | loki_events
    loki_events l5@-->| | loki_write
    otel_exp_loki l6@-->| | loki_write
    loki_write l7@-->|"push"| loki

    %% OTel pipeline
    apps o1@-->|"OTLP gRPC/HTTP"| otel_recv
    otel_recv o2@-->| | otel_mem
    otel_mem o3@-->|"traces"| otel_sampler
    otel_mem o4@-->|"metrics/logs"| otel_resdet
    otel_sampler o5@-->|"traces"| otel_resdet
    otel_resdet o6@-->| | otel_k8sattr
    otel_k8sattr o7@-->|"metrics/logs"| otel_batch
    otel_k8sattr o8@-->|"traces"| otel_transform
    otel_batch o9@-->|"metrics"| otel_exp_prom
    otel_batch o10@-->|"logs"| otel_exp_loki
    otel_transform o11@-->|"traces"| otel_exp_otlp
    otel_exp_otlp o12@-->|"OTLP"| tempo

    classDef metrics stroke:#f5a623
    classDef logs stroke:#27ae60
    classDef traces stroke:#9b59b6

    class m1,m2,m3,m4,m5,m6,m7,m8,m9,m10,m11,m12,m13,m14 metrics
    class l1,l2,l3,l4,l5,l6,l7 logs
    class o1,o2,o3,o4,o5,o6,o7,o8,o9,o10,o11,o12 traces
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
