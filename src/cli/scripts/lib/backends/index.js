'use strict';

const { BeadsBackend } = require('./beads');

function createBackend(config, projectRoot) {
  const adapterType = config?.taskBackend?.adapter || 'beads';
  switch (adapterType) {
    case 'beads':
    default:
      return new BeadsBackend(projectRoot);
  }
}

module.exports = { createBackend };
