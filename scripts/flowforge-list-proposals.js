#!/usr/bin/env node

const { listDocsWorkspaces, listProposalSummaries } = require('./lib/flowforge');

function main() {
  const args = process.argv.slice(2);
  const workspaceIndex = args.findIndex((arg) => arg === '--workspace');
  const workspace = workspaceIndex >= 0 ? args[workspaceIndex + 1] : null;
  const allWorkspaces = args.includes('--all-workspaces');

  let summaries = [];
  if (allWorkspaces) {
    for (const entry of listDocsWorkspaces(process.cwd())) {
      summaries.push(...listProposalSummaries(process.cwd(), entry.name));
    }
  } else {
    summaries = listProposalSummaries(process.cwd(), workspace);
  }

  const grouped = {};

  for (const summary of summaries) {
    if (!grouped[summary.status]) grouped[summary.status] = [];
    grouped[summary.status].push(summary);
  }

  console.log(JSON.stringify(grouped, null, 2));
}

main();
