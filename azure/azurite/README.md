# azurite

Azure Storege Account のエミュレータ。

https://learn.microsoft.com/en-us/azure/storage/common/storage-use-azurite

- デフォルトのアカウント名とキーは
  - `devstoreaccount1`
  - `Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==`
- 環境変数 `AZURITE_ACCOUNTS` でアカウント名とキーが指定できる
  - キーは base64 encode する必要がある事に注意
  - base64 じゃなくても起動するけど、azure-cli で繋ぐときにこけたりする

## azure-cli で接続する

```
% az storage container list --connection-string 'DefaultEndpointsProtocol=http;AccountName=azurite;AccountKey=YXp1cml0ZQ==;BlobEndpoint=http://azurite.localhost:10000;'
```

または

```
% export AZURE_STORAGE_CONNECTION_STRING='DefaultEndpointsProtocol=http;AccountName=azurite;AccountKey=YXp1cml0ZQ==;BlobEndpoint=http://127.0.0.1:10000/azurite;'
% az storage container list
```
