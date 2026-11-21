// Det lilla av Node som vite.config.ts använder, så att tsc kan kontrollera
// den utan hela @types/node.
declare module "node:fs" {
  export function readFileSync(path: string, encoding: "utf8"): string;
}
