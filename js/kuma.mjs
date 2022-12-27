import * as process from "process"

function kuma(size) {
  const kumas = ['ʕ•̫͡•', 'ʔ•̫͡•'];
  const indexes = '001101';
  let kuma = ""
  for (let i = 0; i < size; i++) {
    kuma += kumas[indexes[i % 6]];
  }

  return kuma + 'ʔ';
}

console.log(kuma(Number(process.argv[2] || 10)));
