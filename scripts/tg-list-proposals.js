#!/usr/bin/env node

const { listProposalSummaries } = require('./lib/tg-workflow');

function main() {
  const summaries = listProposalSummaries(process.cwd());
  const grouped = {};

  for (const summary of summaries) {
    if (!grouped[summary.status]) grouped[summary.status] = [];
    grouped[summary.status].push(summary);
  }

  console.log(JSON.stringify(grouped, null, 2));
}

main();
