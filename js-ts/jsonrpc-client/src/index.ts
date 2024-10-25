import * as cp from 'child_process';
import * as rpc from 'vscode-jsonrpc/node';

type Req = {
  woo: string
}

async function main() {
  let childProcess = cp.spawn(process.env["TAPLO_LSP"] ?? "taplo-lsp", ["run"]);

  let connection = rpc.createMessageConnection(
    new rpc.StreamMessageReader(childProcess.stdout),
    new rpc.StreamMessageWriter(childProcess.stdin));
  
  connection.listen();
  let request = new rpc.RequestType<Req, any, any>("initialize");
  console.log(await connection.sendRequest(request, { woo: "moo"}));
}


(async function() {
  await main();
})()
