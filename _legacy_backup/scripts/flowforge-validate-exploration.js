#!/usr/bin/env node

const { loadExplorationContext, validateExplorationContext } = require('./lib/flowforge');

function main() {
  const target = process.argv[2];
  if (!target) {
    console.error('Usage: scripts/flowforge-validate-exploration.js <exploration-slug|exploration-dir>');
    process.exit(2);
  }

  try {
    const context = loadExplorationContext(target, process.cwd());
    const result = validateExplorationContext(context, process.cwd());

    if (result.errors.length === 0 && result.warnings.length === 0) {
      console.log(`OK ${context.explorationDir}`);
      return;
    }

    for (const warning of result.warnings) {
      console.log(`WARN ${warning}`);
    }
    for (const error of result.errors) {
      console.error(`ERROR ${error}`);
    }

    if (result.errors.length > 0) {
      process.exit(1);
    }
  } catch (error) {
    console.error(`ERROR ${error.message}`);
    process.exit(1);
  }
}

main();
