# botkube

https://docs.botkube.io/

## Setup Slack App
https://docs.botkube.io/installation/slack/socket-slack

## Deploy

```
% cat <<EOF > envs/default.yaml
slackAppToken: xapp-XXX
slackBotToken: xoxb-XXX
slackChannel: XXX
EOF
% kind create cluster
% helmfile apply
```

## Multi Cluster

クラスターごとに分ける必要がある。
https://docs.botkube.io/installation/slack/socket-slack



### 同じ Slack App を使う (NG)

```
% cat <<EOF > envs/default.yaml
slackAppToken: xapp-XXX
slackBotToken: xoxb-XXX
EOF
% echo 'slackChannel: XXX' > envs/botkube1.yaml
% echo 'slackChannel: YYY' > envs/botkube2.yaml
% kind create cluster -n botkube1
% kind create cluster -n botkube2
% helmfile -f helmfile.multi-cluster.yaml -e botkube1 apply
% helmfile -f helmfile.multi-cluster.yaml -e botkube2 apply
```

### 別の Slack App を使う (OK)
```
% cat <<EOF > envs/botkube1.yaml
slackAppToken: xapp-XXX
slackBotToken: xoxb-XXX
slackChannel: XXX
EOF
% cat <<EOF > envs/botkube2.yaml
slackAppToken: xapp-YYY
slackBotToken: xoxb-YYY
slackChannel: YYY
EOF
% kind create cluster -n botkube1
% kind create cluster -n botkube2
% helmfile -f helmfile.multi-cluster.yaml -e botkube1 apply
% helmfile -f helmfile.multi-cluster.yaml -e botkube2 apply
```
