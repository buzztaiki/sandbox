# Deployment と CronJob の両方が selector にマッチする場合

当然 Service が両方を見ようとするから、job に転送しようとした場合にエラーになる。

```
❯❯ kubectl get pods
NAME                      READY   STATUS      RESTARTS   AGE
test-54c987d97c-8bqmw     1/1     Running     0          11m
test-job-28224172-vjc4c   0/1     Completed   0          4m4s
test-job-28224173-jljqh   0/1     Completed   0          3m4s
test-job-28224174-lskql   1/1     Running     0          2m4s
test-job-28224175-r76td   1/1     Running     0          64s
test-job-28224176-9g85t   1/1     Running     0          4s
```

```
❯❯ kubectl describe endpoints
Name:         test
Namespace:    test
Labels:       app=test
Annotations:  <none>
Subsets:
  Addresses:          10.244.0.17,10.244.0.24,10.244.0.25,10.244.0.26
  NotReadyAddresses:  <none>
  Ports:
    Name     Port  Protocol
    ----     ----  --------
    <unset>  80    TCP

Events:  <none>
```

```
❯❯ kubectl run busybox --image busybox -it --rm
If you don't see a command prompt, try pressing enter.

/ # wget -q -O- http://test >/dev/null && echo ok
ok
/ # wget -q -O- http://test >/dev/null && echo ok
wget: can't connect to remote host (10.96.248.234): Connection refused
/ # wget -q -O- http://test >/dev/null && echo ok
ok
/ # wget -q -O- http://test >/dev/null && echo ok
ok
/ # wget -q -O- http://test >/dev/null && echo ok
ok
/ # wget -q -O- http://test >/dev/null && echo ok
wget: can't connect to remote host (10.96.248.234): Connection refused
/ # wget -q -O- http://test >/dev/null && echo ok
wget: can't connect to remote host (10.96.248.234): Connection refused
```
