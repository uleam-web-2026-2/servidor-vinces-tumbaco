import { spawn } from 'child_process';

console.log('--- Verificando arranque del servidor Go ---');

const server = spawn('go', ['run', 'main.go']);

server.stdout.on('data', (data) => {
  console.log(`[Servidor Output]: ${data}`);
});

server.stderr.on('data', (data) => {
  console.error(`[Servidor Error]: ${data}`);
});

setTimeout(() => {
  console.log('Comprobación completada. Cerrando proceso...');
  server.kill();
  process.exit(0);
}, 3000);