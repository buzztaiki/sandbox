# CronJob の Forbid の挙動を確認する

[CronJobを使用して自動化タスクを実行する | Kubernetes](https://kubernetes.io/ja/docs/tasks/job/automated-tasks-with-cron-jobs/#concurrency-policy)

- 毎秒実行するように schedule を設定
- Job は 90秒 sleep するだけ

## startingDeadlineSeconds の設定がない場合

[forbid-overlap.yaml](forbid-overlap.yaml) で確認した。
普通に前のジョブが終わった後に実行される。実装面倒で確認してないけど、実行されてなかった分はそのまま消えてると思われる。

```
% kubectl get job -ojson | jq -C '.items[] | {"name": .metadata.name, "startTime": .status.startTime, "completionTime": .status.completionTime}' && kubectl get cronjob forbid-overlap -ojson | jq -C '.status.lastScheduleTime'
{
  "name": "forbid-overlap-27583956",
  "startTime": "2022-06-12T12:36:00Z",
  "completionTime": "2022-06-12T12:37:34Z"
}
{
  "name": "forbid-overlap-27583957",
  "startTime": "2022-06-12T12:37:34Z",
  "completionTime": "2022-06-12T12:39:09Z"
}
{
  "name": "forbid-overlap-27583959",
  "startTime": "2022-06-12T12:39:09Z",
  "completionTime": "2022-06-12T12:40:43Z"
}
{
  "name": "forbid-overlap-27583960",
  "startTime": "2022-06-12T12:40:43Z",
  "completionTime": "2022-06-12T12:42:17Z"
}
{
  "name": "forbid-overlap-27583962",
  "startTime": "2022-06-12T12:42:17Z",
  "completionTime": null
}
"2022-06-12T12:42:00Z"
```

## startingDeadlineSeconds の設定がある場合

[forbid-overlap-with-deadline.yaml](forbid-overlap-with-deadline.yaml) で確認した。
この時間だけ経過すると実行されない。正しい。

```console
% kubectl get job -ojson | jq -C '.items[] | {"name": .metadata.name, "startTime": .status.startTime, "completionTime": .status.completionTime}' && kubectl get cronjob forbid-overlap-with-deadline -ojson | jq -C '.status.lastScheduleTime'
{
  "name": "forbid-overlap-with-deadline-27583969",
  "startTime": "2022-06-12T12:49:00Z",
  "completionTime": "2022-06-12T12:50:34Z"
}
{
  "name": "forbid-overlap-with-deadline-27583971",
  "startTime": "2022-06-12T12:51:00Z",
  "completionTime": "2022-06-12T12:52:34Z"
}
{
  "name": "forbid-overlap-with-deadline-27583973",
  "startTime": "2022-06-12T12:53:00Z",
  "completionTime": "2022-06-12T12:54:34Z"
}
{
  "name": "forbid-overlap-with-deadline-27583975",
  "startTime": "2022-06-12T12:55:00Z",
  "completionTime": "2022-06-12T12:56:34Z"
}
{
  "name": "forbid-overlap-with-deadline-27583977",
  "startTime": "2022-06-12T12:57:00Z",
  "completionTime": "2022-06-12T12:58:34Z"
}
{
  "name": "forbid-overlap-with-deadline-27583979",
  "startTime": "2022-06-12T12:59:00Z",
  "completionTime": null
}
"2022-06-12T12:59:00Z"
```
