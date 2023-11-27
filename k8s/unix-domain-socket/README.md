# Pod で unix domain socket 使えるかテスト

以下の二つのコンテナを pod で動かして、unix domain socket で通信できるかテストしてみた。
- server: socat で /app/socket/test.sock を listen する
- client: 1秒毎に /app/socket/test.sock に今の時間を書き込む

## デプロイ

```
❯❯ just up
kind create cluster -n unix-domain-socket
Creating cluster "unix-domain-socket" ...
 ✓ Ensuring node image (kindest/node:v1.27.3) 🖼
 ✓ Preparing nodes 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
Set kubectl context to "kind-unix-domain-socket"
You can now use your cluster with:

kubectl cluster-info --context kind-unix-domain-socket

Not sure what to do next? 😅  Check out https://kind.sigs.k8s.io/docs/user/quick-start/
docker build -t local/unix-domain-socket:local .
[+] Building 0.8s (6/6) FINISHED                                                                                                                                                                                                                                            docker:default
 => [internal] load build definition from Dockerfile                                                                                                                                                                                                                                  0.1s
 => => transferring dockerfile: 148B                                                                                                                                                                                                                                                  0.0s
 => [internal] load .dockerignore                                                                                                                                                                                                                                                     0.2s
 => => transferring context: 2B                                                                                                                                                                                                                                                       0.0s
 => [internal] load metadata for docker.io/library/debian:bookworm-slim                                                                                                                                                                                                               0.0s
 => [1/2] FROM docker.io/library/debian:bookworm-slim                                                                                                                                                                                                                                 0.0s
 => CACHED [2/2] RUN apt-get update && apt-get install -y socat     && rm -rf /var/lib/apt/lists/*                                                                                                                                                                                    0.0s
 => exporting to image                                                                                                                                                                                                                                                                0.0s
 => => exporting layers                                                                                                                                                                                                                                                               0.0s
 => => writing image sha256:209915297f12d55e3267ae3aa3cd6a2496c549a67a927ed2b8f555f8b2c94994                                                                                                                                                                                          0.0s
 => => naming to docker.io/local/unix-domain-socket:local                                                                                                                                                                                                                             0.0s
kind load docker-image -n unix-domain-socket local/unix-domain-socket:local
Image: "local/unix-domain-socket:local" with ID "sha256:209915297f12d55e3267ae3aa3cd6a2496c549a67a927ed2b8f555f8b2c94994" not yet present on node "unix-domain-socket-control-plane", loading...
kubectl apply --context kind-unix-domain-socket -f all.yaml --wait=true
namespace/unix-domain-socket created
deployment.apps/test created
kubectl --context kind-unix-domain-socket -n unix-domain-socket get all
NAME                   READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/test   0/1     0            0           0s

❯❯ just get_all
kubectl --context kind-unix-domain-socket -n unix-domain-socket get all
NAME                        READY   STATUS    RESTARTS   AGE
pod/test-695765d468-wvcwf   2/2     Running   0          10s

NAME                   READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/test   1/1     1            1           14s

NAME                              DESIRED   CURRENT   READY   AGE
replicaset.apps/test-695765d468   1         1         1       11s
```

## ログ見る

```
❯❯ just tail
kubectl stern --context kind-unix-domain-socket -n unix-domain-socket test
+ test-695765d468-wvcwf › server
+ test-695765d468-wvcwf › client
test-695765d468-wvcwf server Mon Nov 27 02:43:34 UTC 2023
test-695765d468-wvcwf server Mon Nov 27 02:43:35 UTC 2023
test-695765d468-wvcwf server Mon Nov 27 02:43:36 UTC 2023
```


## 後片付け

```
❯❯ just down
kind delete cluster -n unix-domain-socket
Deleting cluster "unix-domain-socket" ...
Deleted nodes: ["unix-domain-socket-control-plane"]
```
