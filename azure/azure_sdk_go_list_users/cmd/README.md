# Azure SDK for Go で Entra (AzureAD) User を一覧してみる

- [no_paging](no_paging)
  - ページングしない
  - 簡単
  - 件数多いとだめ
- [paging_iterator](paging_iterator)
  - `msgraphcore.PageIterator` を使う
  - https://github.com/microsoftgraph/msgraph-sdk-go/blob/main/README.md#41-get-all-the-users-in-an-environment を参考にした
  - 一番抽象化できるけど、iterator を作るときに型の制約がないのがちょっとつらい
  - あまり go っぽくない
  - `CreateUserCollectionResponseFromDiscriminatorValue` がつらい
- [paging_no_iterator](paging_no_iterator)
  - `GetOdataNextLink` を使って自分で追う
  - https://learn.microsoft.com/ja-jp/graph/sdks/paging?tabs=go#manually-requesting-subsequent-pages を参考にした
  - 一番 go っぽい気がする
  - go の generics だと抽象化は多分できない
    - interface に型パラメータが付けられないのがつらい
