# terraform + azure mysql で Entra ID (AzureAD) 連携する

- https://learn.microsoft.com/ja-jp/azure/mysql/single-server/concepts-azure-ad-authentication
- https://learn.microsoft.com/ja-jp/azure/mysql/single-server/how-to-configure-sign-in-azure-ad-authentication

## 作られるもの

- リソースグループ: mysql-entra
- Entra Group:
  - mysql-entra-admin
  - mysql-entra-reader
- MySQL サーバ: mysql-entra
- Entra Group と紐付いた MySQL ユーザ:
  - mysql-entra-admin: 特権ユーザ
  - mysql-entra-reader: 閲覧ユーザ
  - mysql-entra-maintainer: 管理ユーザ

## apply

```
% terraform init
% terraform apply
```

## MySQL に Entra Group と紐付いたユーザを作る

```
% ./mysql.sh -u mysql-entra-admin <<EOF
create aaduser if not exists 'mysql-entra-reader';
grant select on *.* to 'mysql-entra-reader';

create aaduser if not exists 'mysql-entra-maintainer';
-- all privileges はできないらしい
grant select, insert, update, delete, create, drop, references, index, alter, create temporary tables, lock tables, execute, create view, show view, create routine, alter routine, event, trigger on *.* to 'mysql-entra-maintainer';

create aaduser if not exists 'mysql-entra-parent';
grant select on *.* to 'mysql-entra-parent';
EOF
```

## 接続

```
% ./mysql.sh -u mysql-entra-admin -e 'show databases'

+--------------------+
| Database           |
+--------------------+
| information_schema |
| mysql              |
| performance_schema |
| sys                |
+--------------------+

❯❯ ./mysql.sh -u mysql-entra-reader -e 'show databases'
+--------------------+
| Database           |
+--------------------+
| information_schema |
| mysql              |
| performance_schema |
| sys                |
+--------------------+
```

## Audit Log を見る
Entra Group と紐付いたユーザでログインしても、Audit Log からは実ユーザが見れる。べんり。

```
% az extension add -n log-analytics 
% az monitor log-analytics query -w "$(az monitor log-analytics workspace show -g mysql-entra -n mysql-entra --query customerId -o tsv)" -otable --analytics-query '
AzureDiagnostics 
| where ResourceGroup == toupper("mysql-entra")
| where Category == "MySqlAuditLogs"
| where event_class_s == "connection_log"
| project TimeGenerated, event_class_s, event_subclass_s, user_s, external_user_s
'

TableName      TimeGenerated                 Event_class_s    Event_subclass_s    External_user_s                    User_s
-------------  ----------------------------  ---------------  ------------------  ---------------------------------  -----------------
PrimaryResult  2023-12-20T09:12:31.8651097Z  connection_log   DISCONNECT                                             mysql-entra-admin
PrimaryResult  2023-12-20T09:05:31.9454435Z  connection_log   DISCONNECT                                             mysql-entra-admin
PrimaryResult  2023-12-20T09:11:31.6828909Z  connection_log   AADAUTH             upn:live.com#buzz.taiki@gmail.com  mysql-entra-admin
PrimaryResult  2023-12-20T09:11:31.6828909Z  connection_log   CONNECT                                                mysql-entra-admin
PrimaryResult  2023-12-20T08:44:27.3987488Z  connection_log   DISCONNECT                                             mysql-entra-admin
PrimaryResult  2023-12-20T08:46:49.2598392Z  connection_log   AADAUTH             upn:live.com#buzz.taiki@gmail.com  mysql-entra-admin
PrimaryResult  2023-12-20T08:46:49.2598392Z  connection_log   CONNECT                                                mysql-entra-admin
PrimaryResult  2023-12-20T08:59:27.3062623Z  connection_log   AADAUTH             upn:live.com#buzz.taiki@gmail.com  mysql-entra-admin
PrimaryResult  2023-12-20T08:59:27.3062623Z  connection_log   CONNECT                                                mysql-entra-admin
```

## 後始末

```
% terraform destroy
```
