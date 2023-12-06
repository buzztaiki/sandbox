# SP を使って terraform から AzureAD (Entra ID) を操作する

## ディレクトリ構成
- prepare: app と sp を作る。これは管理者で実行する。
- main: 作成した sp (app) を使って plan や apply できるか確認する。

## 準備
```
% cd prepare
% az login
% terraform apply
```

## plan sp で実行
```
% cd main
% terraform plan --var-file=plan.tfvars
% terraform apply --var-file=plan.tfvars
エラーになる
```

## apply sp で実行
```
% cd main
% terraform plan --var-file=apply.tfvars
% terraform apply --var-file=apply.tfvars
```

## メモ
- AzureAD を操作するには MicrosoftGraph App の操作を許可する必要がある
  - https://learn.microsoft.com/en-us/graph/auth-v2-service
  - https://jpazureid.github.io/blog/azure-active-directory/oauth2-application-resource-and-api-permissions/
- MicrosoftGraph App の Client ID は `00000003-0000-0000-c000-000000000000`
  - 普通に Enterprise Application の中に登録されてる
- アプリが利用できるロールの許可を得て、クライアント (SP) がその範囲内でロールを要求する形になる
- 画面からロールを設定する場合は以下から設定する
  - **[Microsoft Entra ID]** > **[App registrations]** > **[All Applications]** > **[\<app name\>]** > **[API Permissions]**
  - https://portal.azure.com/#view/Microsoft_AAD_RegisteredApps/ApplicationMenuBlade/~/CallAnAPI/appId/:appid/isMSAApp~/false
- 利用できるロールは以下に記述されている
  - https://learn.microsoft.com/en-us/graph/permissions-reference

