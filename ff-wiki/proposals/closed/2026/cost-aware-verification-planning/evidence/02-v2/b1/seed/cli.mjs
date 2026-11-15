import { quote } from './quote.mjs';
import { help } from './help.mjs';
if (process.argv[2] === '--help') {
  console.log(help);
} else {
  try {
    const cents = quote(JSON.parse(process.argv[2]));
    console.log(JSON.stringify({ cents }));
  } catch (error) {
    console.error(error.message);
    process.exitCode = 2;
  }
}
