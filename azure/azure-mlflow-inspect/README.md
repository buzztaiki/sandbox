# azure-mlflow-inspect

mlflow を使って AzureML の情報を取ってくるテスト

## 準備

AzureML の config.json をダウンロードして、ここに置いておく。


## 使い方

```
# run の一覧を取ってくる
% uv run list_all_runs.py

# run を inspect する
% uv run inspect_run.py <run-name>
```
