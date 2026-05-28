#!/usr/bin/env node

const { listDocsWorkspaces, listProposalSummaries, listExplorationSummaries } = require('./lib/flowforge');

function main() {
  const args = process.argv.slice(2);
  const kindIndex = args.findIndex((arg) => arg === '--kind');
  const kind = kindIndex >= 0 ? args[kindIndex + 1] : 'all';
  const workspaceIndex = args.findIndex((arg) => arg === '--workspace');
  const workspace = workspaceIndex >= 0 ? args[workspaceIndex + 1] : null;
  const allWorkspaces = args.includes('--all-workspaces');
  const configuredWorkspaces = listDocsWorkspaces(process.cwd());
  const listAll = allWorkspaces || (!workspace && configuredWorkspaces.length > 1);

  let summaries = [];
  const includeProposals = kind === 'all' || kind === 'proposals';
  const includeExplorations = kind === 'all' || kind === 'explorations';

  if (listAll) {
    for (const entry of configuredWorkspaces) {
      if (includeProposals) {
        summaries.push(...listProposalSummaries(process.cwd(), entry.name));
      }
      if (includeExplorations) {
        summaries.push(...listExplorationSummaries(process.cwd(), entry.name));
      }
    }
  } else {
    if (includeProposals) {
      summaries.push(...listProposalSummaries(process.cwd(), workspace));
    }
    if (includeExplorations) {
      summaries.push(...listExplorationSummaries(process.cwd(), workspace));
    }
  }

  const grouped = {};

  for (const summary of summaries) {
    if (!grouped[summary.status]) grouped[summary.status] = [];
    grouped[summary.status].push(summary);
  }

  console.log(JSON.stringify(grouped, null, 2));
}

main();
