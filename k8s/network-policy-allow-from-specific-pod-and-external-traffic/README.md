# 特定の Pod と外部からのトラフィックだけを許可するネットワークポリシー

https://kubernetes.io/ja/docs/concepts/services-networking/network-policies/

どうするのがいいんだろう的にやったもの。

## まとめ

- `ipBlock` に VPC のネットワークアドレス的なものを指定すれば良さそうではある
  - CNI、LoadBalancer、Ingress Controller 次第な気がしなくもない
- `0.0.0.0/0` を指定した場合
   - Calico は内部の疎通も全て通る
   - Cilium は内部の疎通は通らない
   - へぇぇぇってなった。
- 外部からアクセス可能な Pod は空のブロックを指定して、内部からも全て疎通可能とするのが良い気もする

## 動作確認

### マニフェストの適用と確認

```console
% kubectl apply -f pods.yaml
pod/allow created
pod/deny created
pod/nginx created
service/nginx created
service/nginx-external created

% kubectl get -f pods.yaml -owide
NAME        READY   STATUS    RESTARTS   AGE     IP              NODE       NOMINATED NODE   READINESS GATES
pod/allow   1/1     Running   0          7m39s   10.244.120.68   minikube   <none>           <none>
pod/deny    1/1     Running   0          7m39s   10.244.120.67   minikube   <none>           <none>
pod/nginx   1/1     Running   0          7m39s   10.244.120.69   minikube   <none>           <none>

NAME                     TYPE           CLUSTER-IP      EXTERNAL-IP     PORT(S)          AGE     SELECTOR
service/nginx            ClusterIP      10.103.122.95   <none>          8000/TCP         7m39s   name=nginx
service/nginx-external   LoadBalancer   10.101.129.11   10.101.129.11   8001:32297/TCP   7m39s   name=nginx
```

### 外部を全て許可した場合の疎通確認 (Calico)

```console
% kubectl apply -f allow-external-all.yaml
networkpolicy.networking.k8s.io/ingress-deny-all created
networkpolicy.networking.k8s.io/nginx-ingress created

% kubectl describe networkpolicies nginx-ingress | awk '/^Spec/ {ok=1} ok'
Spec:
  PodSelector:     name=nginx
  Allowing ingress traffic:
    To Port: 80/TCP
    From:
      PodSelector: name=allow
    From:
      IPBlock:
        CIDR: 0.0.0.0/0
        Except:
  Not affecting egress traffic
  Policy Types: Ingress

% ./curl.sh allow
+ kubectl exec allow -- curl -s -I --connect-timeout 2 http://nginx:8000
HTTP/1.1 200 OK
...

% ./curl.sh deny
+ kubectl exec deny -- curl -s -I --connect-timeout 2 http://nginx:8000
HTTP/1.1 200 OK
...

% curl -s -I --connect-timeout 2 http://10.101.129.11:8001
HTTP/1.1 200 OK
...
```

### 外部を一部許可した場合の疎通確認 (Calico)

```console
% kubectl apply -f allow-external-vpc.yaml
networkpolicy.networking.k8s.io/ingress-deny-all unchanged
networkpolicy.networking.k8s.io/nginx-ingress configured

% ./curl.sh allow
+ kubectl exec allow -- curl -s -I --connect-timeout 2 http://nginx:8000
HTTP/1.1 200 OK
..

% ./curl.sh deny
+ kubectl exec deny -- curl -s -I --connect-timeout 2 http://nginx:8000
command terminated with exit code 28
% curl -s -I --connect-timeout 2 http://10.101.129.11:8001
HTTP/1.1 200 OK
...
```

### 外部を全て許可した場合の疎通確認 (Cilium)

```console
% kubectl apply -f allow-external-all.yaml
networkpolicy.networking.k8s.io/ingress-deny-all created
networkpolicy.networking.k8s.io/nginx-ingress created

% kubectl describe networkpolicies nginx-ingress | awk '/^Spec/ {ok=1} ok'
Spec:
  PodSelector:     name=nginx
  Allowing ingress traffic:
    To Port: 80/TCP
    From:
      PodSelector: name=allow
    From:
      IPBlock:
        CIDR: 0.0.0.0/0
        Except:
  Not affecting egress traffic
  Policy Types: Ingress

% ./curl.sh allow
+ kubectl exec allow -- curl -s -I --connect-timeout 2 http://nginx:8000
HTTP/1.1 200 OK
...

% ./curl.sh deny
+ kubectl exec deny -- curl -s -I --connect-timeout 2 http://nginx:8000
command terminated with exit code 28

% curl -s -I --connect-timeout 2 http://10.103.175.229:8001
HTTP/1.1 200 OK
...
```

### 外部を一部許可した場合の疎通確認 (Cilium)

% kubectl apply -f allow-external-vpc.yaml
networkpolicy.networking.k8s.io/ingress-deny-all unchanged
networkpolicy.networking.k8s.io/nginx-ingress configured

% kubectl describe networkpolicies nginx-ingress | awk '/^Spec/ {ok=1} ok'
Spec:
  PodSelector:     name=nginx
  Allowing ingress traffic:
    To Port: 80/TCP
    From:
      PodSelector: name=allow
    From:
      IPBlock:
        CIDR: 192.168.49.0/24
        Except:
  Not affecting egress traffic
  Policy Types: Ingress

% ./curl.sh allow
+ kubectl exec allow -- curl -s -I --connect-timeout 2 http://nginx:8000
HTTP/1.1 200 OK
...

% ./curl.sh deny
+ kubectl exec deny -- curl -s -I --connect-timeout 2 http://nginx:8000
command terminated with exit code 28

% curl -s -I --connect-timeout 2 http://10.103.175.229:8001
HTTP/1.1 200 OK
...

## CNI プラグインの有効化

Calico:

```console
% minikube delete
% minikube start --cni calico
```

Cilium

```console
% minikube delete
% minikube start --cni cilium
```
