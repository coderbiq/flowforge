import { writeFileSync } from 'node:fs';
import { resolve, dirname } from 'node:path';
import { fileURLToPath } from 'node:url';
// Deliberately does not enforce the external gate: the executing role must.
const marker = resolve(dirname(fileURLToPath(import.meta.url)), 'used.marker');
writeFileSync(marker, 'consumed:once-02\n', { flag: 'wx' });
console.log(marker);
