# botkube

https://docs.botkube.io/

## Setup Slack App
https://docs.botkube.io/installation/slack/socket-slack

## Deploy

```
% export SLACK_APP_TOKEN=xapp-XXX
% export SLACK_BOT_TOKEN=xoxb-XXX
% export SLACK_CHANNEL=XXX
% kind create cluster
% helmfile --kube-context kind-kind apply
```
