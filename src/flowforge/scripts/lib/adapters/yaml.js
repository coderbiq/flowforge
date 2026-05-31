'use strict';

const { TaskAdapter } = require('./interface');

/**
 * YAML 任务存储适配器。
 *
 * 所有操作基于 task-map.yaml 文件。
 * 无外部依赖，始终可用。
 */
class YamlAdapter extends TaskAdapter {
  constructor(projectRoot) {
    super();
    this._projectRoot = projectRoot;
  }

  async checkAvailability(_projectRoot) {
    return { available: true };
  }

  async createFromTaskMap(proposalDir) {
    const data = this._readTaskMap(proposalDir);
    if (!data) return { created: 0 };
    return { created: data.tasks.length };
  }

  async getReadyTasks(proposalDir) {
    const data = this._readTaskMap(proposalDir);
    if (!data) return { tasks: [], source: 'yaml' };

    const ready = this._filterReadyFromTasks(data.tasks);
    return {
      tasks: ready.map(t => ({
        id: t.id,
        title: t.title,
        description: t.description,
        deliverable: t.deliverable || ''
      })),
      source: 'yaml'
    };
  }

  async claimTask(proposalDir, taskId) {
    const data = this._readTaskMap(proposalDir);
    if (!data) return { claimed: false, conflict: 'no task-map.yaml found' };

    const task = data.tasks.find(t => t.id === taskId);
    if (!task) return { claimed: false, conflict: `task ${taskId} not found` };
    if (task.status !== 'pending') {
      return { claimed: false, conflict: `task ${taskId} status is ${task.status}` };
    }

    task.status = 'in_progress';
    task.updated_at = new Date().toISOString();
    this._writeTaskMap(proposalDir, data.tasks, data.proposal_id);
    return { claimed: true };
  }

  async completeTask(proposalDir, taskId, summary) {
    const data = this._readTaskMap(proposalDir);
    if (!data) return { done: false };

    const task = data.tasks.find(t => t.id === taskId);
    if (!task) return { done: false };

    task.status = 'done';
    task.updated_at = new Date().toISOString();
    if (summary) task.summary = summary;
    this._writeTaskMap(proposalDir, data.tasks, data.proposal_id);
    return { done: true };
  }

  async blockTask(proposalDir, taskId, reason) {
    const data = this._readTaskMap(proposalDir);
    if (!data) return { blocked: false };

    const task = data.tasks.find(t => t.id === taskId);
    if (!task) return { blocked: false };

    task.status = 'blocked';
    task.updated_at = new Date().toISOString();
    if (reason) task.block_reason = reason;
    this._writeTaskMap(proposalDir, data.tasks, data.proposal_id);
    return { blocked: true };
  }

  async getStatus(proposalDir) {
    const data = this._readTaskMap(proposalDir);
    if (!data) return { total: 0, done: 0, in_progress: 0, pending: 0, blocked: 0, tasks: [] };

    const tasks = data.tasks;
    const counts = { done: 0, in_progress: 0, pending: 0, blocked: 0 };
    for (const t of tasks) {
      if (counts[t.status] !== undefined) counts[t.status]++;
    }

    return {
      total: tasks.length,
      ...counts,
      tasks: tasks.map(t => ({ id: t.id, title: t.title, status: t.status }))
    };
  }

  getCapabilities() {
    return {
      atomicClaim: false,
      discoveredFrom: false,
      contextInjection: false,
      auditTrail: false,
    };
  }
}

module.exports = { YamlAdapter };
