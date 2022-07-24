- https://www.telepresence.io/docs/latest/quick-start/qs-node/

## Worklog

minikube

```console
 $  minikube start
...
🏄  完了しました！ kubectl が「"minikube"」クラスタと「"default"」ネームスペースを使用するよう構成されました
```

---

telepresence connect

```console
 $  telepresence connect
Launching Telepresence Root Daemon
Need root privileges to run: /usr/bin/telepresence daemon-foreground /home/taiki/.cache/telepresence/logs /home/taiki/.config/telepresence ''
[sudo] taiki のパスワード:
Launching Telepresence User Daemon

Connected to context minikube (https://192.168.49.2:8443)
```

---

https://kubernetes.default につながる

```console
 $  curl -ik https://kubernetes.default
HTTP/2 403
cache-control: no-cache, private
content-type: application/json
x-content-type-options: nosniff
x-kubernetes-pf-flowschema-uid: 917d15f1-55d8-47cb-84e9-1e680ea2a499
x-kubernetes-pf-prioritylevel-uid: 1913e095-e3ec-4afb-988f-557befda2404
content-length: 233
date: Fri, 06 Aug 2021 04:27:48 GMT

{
  "kind": "Status",
  "apiVersion": "v1",
  "metadata": {

  },
  "status": "Failure",
  "message": "forbidden: User \"system:anonymous\" cannot get path \"/\"",
  "reason": "Forbidden",
  "details": {

  },
  "code": 403
}
```

---

スタックをデプロイ。

```console
 $  kubectl apply -f https://raw.githubusercontent.com/datawire/edgey-corp-nodejs/main/k8s-config/edgey-corp-web-app-no-mapping.yaml
deployment.apps/dataprocessingservice created
service/dataprocessingservice created
deployment.apps/verylargejavaservice created
service/verylargejavaservice created
deployment.apps/verylargedatastore created
service/verylargedatastore created

 $  kubectl get all
NAME                                         READY   STATUS    RESTARTS   AGE
pod/dataprocessingservice-685cb9d6f6-99tkj   1/1     Running   0          7m2s
pod/verylargedatastore-98d78d474-64v44       1/1     Running   0          7m2s
pod/verylargejavaservice-689dbc854b-fpldx    1/1     Running   0          7m2s

NAME                            TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
service/dataprocessingservice   ClusterIP   10.100.123.231   <none>        3000/TCP   7m2s
service/verylargedatastore      ClusterIP   10.104.110.238   <none>        8080/TCP   7m2s
service/verylargejavaservice    ClusterIP   10.108.92.185    <none>        8080/TCP   7m2s

NAME                                    READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/dataprocessingservice   1/1     1            1           7m2s
deployment.apps/verylargedatastore      1/1     1            1           7m2s
deployment.apps/verylargejavaservice    1/1     1            1           7m2s

NAME                                               DESIRED   CURRENT   READY   AGE
replicaset.apps/dataprocessingservice-685cb9d6f6   1         1         1       7m2s
replicaset.apps/verylargedatastore-98d78d474       1         1         1       7m2s
replicaset.apps/verylargejavaservice-689dbc854b    1         1         1       7m2s
 ~ 
```

---

`https://<service-name>.<namespace>:<port>` で開く事ができる。

```console
 $  xdg-open http://verylargejavaservice.sandbox:8080/
```

やっぱ名前解決のデーモンが `telepresence connect` で動いてるな？


---

デモ用のソースコードを取ってくる

```console
 $  git clone https://github.com/datawire/edgey-corp-nodejs.git
...
 $  cd edgey-corp-nodejs/DataProcessingService/
 $  npm install && npm start
...
Welcome to the DataProcessingNodeService!
{ _: [] }
Server running on port 3000
```

---

接続

```console
 $  curl http://localhost:3000/color
"blue"
```

telepresence 側につなぐと green がかえる

```console
 $  curl http://dataprocessingservice.sandbox:3000/color
"green" 
```

---

`telepresence intercept` でクラスタ内のトラフィックを intercept してローカルのポートにフォワードできる

```console
 $  telepresence intercept dataprocessingservice --port 3000
Using Deployment dataprocessingservice
intercepted
    Intercept name    : dataprocessingservice
    State             : ACTIVE
    Workload kind     : Deployment
    Destination       : 127.0.0.1:3000
    Volume Mount Error: sshfs is not installed on your local machine
    Intercepting      : all TCP connections
Intercepting all traffic to your service. To route a subset of the traffic instead, use a personal intercept. You can enable personal intercepts by authenticating to the Ambassador Developer Control Plane with "telepresence login".
```

当然、ローカルに変更を加えるとそれが反映される。

---

ambassador.io にログインして intercept すると Preview URL が作れる (ngrok みたいなかんじ)。
このとき、どの service を公開するかを指定する。


```console
 $  telepresence login
Launching browser authentication flow...
Login successful.
```

```console
 $  telepresence intercept dataprocessingservice --port 3000

To create a preview URL, telepresence needs to know how cluster
ingress works for this service.  Please Select the ingress to use.

1/4: What's your ingress' layer 3 (IP) address?
     You may use an IP address or a DNS name (this is usually a
     "service.namespace" DNS name).

       [no default]: verylargejavaservice.sandbox

2/4: What's your ingress' layer 4 address (TCP port number)?

       [no default]: 8080

3/4: Does that TCP port on your ingress use TLS (as opposed to cleartext)?

       [default: n]:

4/4: If required by your ingress, specify a different layer 5 hostname
     (TLS-SNI, HTTP "Host" header) to access this service.

       [default: verylargejavaservice.sandbox]:

Using Deployment dataprocessingservice
intercepted
    Intercept name    : dataprocessingservice
    State             : ACTIVE
    Workload kind     : Deployment
    Destination       : 127.0.0.1:3000
    Volume Mount Error: sshfs is not installed on your local machine
    Intercepting      : HTTP requests that match all headers:
      'x-telepresence-intercept-id: 9ff81e20-fa6b-4049-9863-a634818ab29d:dataprocessingservice'
    Preview URL       : https://objective-brahmagupta-4775.preview.edgestack.me
    Layer 5 Hostname  : verylargejavaservice.sandbox
```

---

Preview URL にアクセスするとインターネットから見る事ができる。

```console
 $  xdg-open https://objective-brahmagupta-4775.preview.edgestack.me
```


## telepresence 用の 名前以外からアクセスした場合は intercept されずに 動くのか？

intercept される。内部ネットワークが変えられてるって事。

```console
 $  curl -s http://verylargejavaservice.sandbox:8080 | grep color
<h1 style="color:orange">Welcome to the EdgyCorp WebApp</h1>
```

```
 $  kubectl run -it --rm busybox --image=busybox
If you don't see a command prompt, try pressing enter.
```
/ # wget -O- http://verylargejavaservice:8080 | grep color 2>/dev/null

```console
 $  kubectl run -it --rm busybox --image=busybox
If you don't see a command prompt, try pressing enter.
/ # wget -O- http://verylargejavaservice:8080 | grep color 2>/dev/null

/ # wget -q -O- http://verylargejavaservice:8080 | grep color 2>/dev/null
<h1 style="color:orange">Welcome to the EdgyCorp WebApp</h1>
```
