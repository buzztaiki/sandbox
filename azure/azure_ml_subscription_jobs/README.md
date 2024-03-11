# azure-ml-subscription-jobs

自分がログインしてるサブスクリプションのAzure MLジョブを全部まとめて取る。

## Usage

```
% rye sync
% rye run app
```

## TODO
- pipeline step (`type=base`) の場合に開始日時が取れない
  - Machine Learning Studio からだと取れるんだけど、SDK から取る口がなさそう。
