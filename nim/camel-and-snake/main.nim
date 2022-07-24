proc nimnimnim(): string =
  "ok"

# camelCase と snake_case が同一される。というか、case insensitive で `_` が無視される。
# see https://nim-lang.org/docs/manual.html#lexical-analysis-identifier-equality
echo nimnimnim()
echo nimNimNim()
echo nim_nim_nim()
echo nim_nimnim()
echo niMniMnim()
