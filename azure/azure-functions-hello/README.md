# azure-functions-hello

## プロジェクトを作る

https://learn.microsoft.com/en-us/azure/azure-functions/create-first-function-cli-typescript
```
% npx azure-functions-core-tools@4 init --typescript
```

## func cli

プロジェクトが作られてれば、devDependencies に入ってるから以下で ok
```
% npx func --help
```

## template list

```
% npx func templates list -l typescript
```

## http trigger

```
% npx func new --name HttpExample --template "HTTP trigger" --authlevel "anonymous"
```

```
% npm start
...
[2025-06-13T16:28:23.553Z] Worker process started and initialized.

Functions:

        HttpExample: [GET,POST] http://localhost:7071/api/HttpExample

```

```
% curl http://localhost:7071/api/HttpExample
Hello, world!
```

