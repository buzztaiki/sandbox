# k8s-monitoring chart + monitoring-mixin

kube-prometheus-stack を使わずに、k8s-monitoring と monitoring-mixin で監視環境を作ってみる
- scrape: k8s-monitoring chart (alloy)
- prometheus-operator: prometheus-operator-crds chart
  - controller が無くても alloy が集めてくれるから問題ない
- dashboard: monitoring-mixin
- alert-rule, recording-rule: monitoring-mixin
  - prometheus-operator のリソースとして作って、mimir.rules.kubernetes で同期する
- grafana: grafana chart


## k8s-monitoring の感想

便利だけど、設定した内容が alloy にどう反映されるかが直感的に分かりにくいのが難点。
一度 alloy を触ったから何とかなったけど、触ってなければだいぶ困ったのでは。
あと、job label が kube-prom と違うから、合わせる必要があるのがちょっと面倒。

## monitoring-mixin を自分でバンドルした感想

やり方に慣れさえすればそこまで面倒でもないし、気にいらない所は手を入れられるので結構いいかもしれない。
