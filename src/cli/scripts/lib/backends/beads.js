'use strict';

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');
const { TaskBackend, BackendCapabilities, TASK_STATUS, TASK_TYPE } = require('./interface');
const { loadMainConfig, findProposalDir } = require('../config');

/**
 * Beads 任务存储后端。
 * 薄封装 bd CLI，beads 是任务的唯一真理源。
 * Agent 不操作文件，所有任务增删改查通过本后端走 bd 命令。
 */
class BeadsBackend extends TaskBackend {
  constructor(projectRoot) {
    super(projectRoot);
    this._epicCache = new Map();
  }

  // ── 能力 ──

  async checkAvailability() {
    try {
      execSync('bd context --json', {
        cwd: this._projectRoot,
        stdio: 'pipe',
        timeout: 5000
      });
      return { available: true };
    } catch (e) {
      return { available: false, reason: _extractError(e) };
    }
  }

  getCapabilities() {
    return new BackendCapabilities({
      atomicClaim: true,
      discoveredFrom: true,
      auditTrail: true,
      dependencySort: true,
    });
  }

  // ── 生命周期 ──

  async init(proposalId, title) {
    const result = this._bd(
      `create "${_escape(title)}" --type epic --labels proposal:${proposalId} --json`
    );
    const epic = JSON.parse(result);
    this._epicCache.set(proposalId, epic.id);
    return { epicId: epic.id };
  }

  async teardown(proposalId) {
    try {
      const epicId = await this._resolveEpic(proposalId);
      if (epicId) {
        this._bd(`close ${epicId} --reason "Proposal archived"`);
      }
    } catch (_) { /* 清理失败不阻塞归档 */ }
    this._epicCache.delete(proposalId);
  }

  // ── 任务创建 ──

  async addTask(proposalId, task, parentTaskId) {
    const epicId = await this._resolveEpic(proposalId);
    if (!epicId) {
      throw new Error(`No epic found for proposal ${proposalId}. Run 'flowforge task init --proposal ${proposalId} "<title>"' first.`);
    }
    const parentId = parentTaskId || epicId;
    const labels = this._buildLabels(task, proposalId);
    const result = this._bd(
      `create "${_escape(task.title)}" --type task --parent ${parentId} ` +
      `--labels ${labels} --description "${_escape(task.description || '')}" --json`
    );
    const issue = JSON.parse(result);
    const taskId = issue.id;

    if (task.dependencies && task.dependencies.length > 0) {
      for (const depId of task.dependencies) {
        try { this._bd(`link ${taskId} ${depId}`); } catch (_) {}
      }
    }

    return { taskId };
  }

  async addTasks(proposalId, tasks) {
    const epicId = await this._resolveEpic(proposalId);
    if (!epicId) {
      throw new Error(`No epic found for proposal ${proposalId}. Run 'flowforge task init' first.`);
    }
    const ids = [];

    for (const t of tasks) {
      const labels = this._buildLabels(t, proposalId);
      const result = this._bd(
        `create "${_escape(t.title)}" --type task --parent ${epicId} ` +
        `--labels ${labels} --description "${_escape(t.description || '')}" --json`
      );
      ids.push(JSON.parse(result).id);
    }

    for (let i = 0; i < tasks.length; i++) {
      const deps = tasks[i].dependencies || [];
      for (const depRef of deps) {
        const depIdx = typeof depRef === 'number' ? depRef : ids.findIndex(id => id === depRef);
        if (depIdx >= 0 && depIdx < ids.length && depIdx !== i) {
          try { this._bd(`link ${ids[i]} ${ids[depIdx]}`); } catch (_) {}
        }
      }
    }

    return { taskIds: ids };
  }

  async discoverTask(proposalId, parentTaskId, task) {
    const epicId = await this._resolveEpic(proposalId);
    if (!epicId) {
      throw new Error(`No epic found for proposal ${proposalId}. Run 'flowforge task init' first.`);
    }
    const labels = this._buildLabels(task, proposalId);
    const depFlag = `--dep discovered-from:${parentTaskId}`;
    const result = this._bd(
      `create "${_escape(task.title)}" --type task --parent ${epicId} ` +
      `${depFlag} --labels ${labels} --json`
    );
    return { taskId: JSON.parse(result).id };
  }

  async cancelTask(proposalId, taskId, reason) {
    const closeReason = reason ? ` --reason "${_escape(reason)}"` : ' --reason "Cancelled"';
    this._bd(`close ${taskId}${closeReason}`);
  }

  // ── 状态流转 ──

  async claimTask(proposalId, taskId) {
    try {
      this._bd(`update ${taskId} --claim`);
      return { claimed: true };
    } catch (e) {
      return { claimed: false, conflict: _extractError(e) };
    }
  }

  async completeTask(proposalId, taskId, summary) {
    const reason = summary ? ` --reason "${_escape(summary)}"` : '';
    this._bd(`close ${taskId}${reason}`);
  }

  async blockTask(proposalId, taskId, reason) {
    try {
      this._bd(`update ${taskId} --status blocked`);
      if (reason) {
        this._bd(`update ${taskId} --reason "${_escape(reason)}"`);
      }
    } catch (_) { /* 阻塞失败不抛异常 */ }
  }

  async unclaimTask(proposalId, taskId) {
    try {
      this._bd(`update ${taskId} --status pending`);
    } catch (_) {}
  }

  async reopenTask(proposalId, taskId) {
    try {
      this._bd(`update ${taskId} --status pending`);
    } catch (_) {}
  }

  // ── 依赖管理 ──

  async addDependency(proposalId, taskId, dependsOnTaskId) {
    this._bd(`link ${taskId} ${dependsOnTaskId}`);
  }

  async removeDependency(proposalId, taskId, dependsOnTaskId) {
    try {
      this._bd(`unlink ${taskId} ${dependsOnTaskId}`);
    } catch (_) {}
  }

  // ── 查询 ──

  async getReadyTasks(proposalId) {
    try {
      const result = this._bd('ready --json');
      const all = JSON.parse(result);
      if (!Array.isArray(all)) return [];
      return all
        .filter(t => {
          const labels = t.labels || [];
          return labels.some(l => l === `proposal:${proposalId}`);
        })
        .map(t => this._toTask(t));
    } catch (_) {
      return this._getReadyFromList(proposalId);
    }
  }

  async getStatus(proposalId) {
    const tasks = await this._listTasks(proposalId);
    const counts = { done: 0, in_progress: 0, pending: 0, blocked: 0 };
    const byType = {};

    for (const t of tasks) {
      if (counts[t.status] !== undefined) counts[t.status]++;
      const typeKey = t.type || TASK_TYPE.IMPLEMENTATION;
      if (!byType[typeKey]) {
        byType[typeKey] = { total: 0, done: 0, inProgress: 0, pending: 0, blocked: 0 };
      }
      byType[typeKey].total++;
      const statusKey = t.status === 'in_progress' ? 'inProgress' : t.status;
      if (byType[typeKey][statusKey] !== undefined) byType[typeKey][statusKey]++;
    }

    return { total: tasks.length, byStatus: counts, byType, tasks };
  }

  async getBlockedTasks(proposalId) {
    const tasks = await this._listTasks(proposalId);
    return tasks.filter(t => t.status === TASK_STATUS.BLOCKED);
  }

  async isAllDone(proposalId) {
    const tasks = await this._listTasks(proposalId);
    if (tasks.length === 0) return false;
    return tasks.every(t =>
      t.status === TASK_STATUS.DONE || t.status === TASK_STATUS.CANCELLED
    );
  }

  // ── 标签管理 ──

  async addLabel(proposalId, taskId, label) {
    this._bd(`label add ${taskId} ${label}`);
  }

  async removeLabel(proposalId, taskId, label) {
    try { this._bd(`label remove ${taskId} ${label}`); } catch (_) {}
  }

  async listLabels(proposalId, taskId) {
    try {
      return this._bd(`label list ${taskId} --json`);
    } catch (_) {
      return JSON.stringify([]);
    }
  }

  // ── 快照导出 ──

  async exportSnapshot(proposalId) {
    const tasks = await this._listTasks(proposalId);
    const proposalDir = this._findProposalDir(proposalId);
    if (!proposalDir) throw new Error(`Proposal dir not found: ${proposalId}`);

    const md = this._formatSnapshot(tasks, proposalId);
    const filePath = path.join(proposalDir, 'tasks.snapshot.md');
    fs.writeFileSync(filePath, md, 'utf8');
    return filePath;
  }

  // ── 迁移支持 ──

  async migrateFromYaml(proposalDir) {
    const yaml = require('../../vendor/js-yaml');
    const taskMapPath = path.join(proposalDir, 'task-map.yaml');
    if (!fs.existsSync(taskMapPath)) {
      return { migrated: 0, skipped: 0, errors: ['task-map.yaml not found'] };
    }

    const data = yaml.load(fs.readFileSync(taskMapPath, 'utf8'));
    const tasks = data?.tasks || [];
    const proposalId = data?.proposal_id || path.basename(proposalDir);

    let epicId;
    try {
      const result = this._bd(
        `create "Proposal: ${_escape(proposalId)}" --type epic --labels proposal:${proposalId} --json`
      );
      epicId = JSON.parse(result).id;
      this._epicCache.set(proposalId, epicId);
    } catch (e) {
      return { migrated: 0, skipped: 0, errors: [`Failed to create epic: ${_extractError(e)}`] };
    }

    let migrated = 0;
    let skipped = 0;
    const errors = [];
    const idMap = new Map();

    for (const task of tasks) {
      try {
        const existing = this._findExistingBead(proposalId, task.title);
        if (existing) {
          idMap.set(task.id, existing);
          skipped++;
          continue;
        }

        const labels = this._buildLabels(task, proposalId);
        const result = this._bd(
          `create "${_escape(task.title)}" --type task --parent ${epicId} ` +
          `--labels ${labels} --description "${_escape(task.description || '')}" --json`
        );
        const beadId = JSON.parse(result).id;
        idMap.set(task.id, beadId);

        if (task.status === TASK_STATUS.IN_PROGRESS) {
          try { this._bd(`update ${beadId} --claim`); } catch (_) {}
        } else if (task.status === TASK_STATUS.DONE) {
          try { this._bd(`close ${beadId} --reason "Migrated from v0.8"`); } catch (_) {}
        } else if (task.status === TASK_STATUS.BLOCKED) {
          try { this._bd(`update ${beadId} --status blocked`); } catch (_) {}
        }

        migrated++;
      } catch (e) {
        errors.push(`Task ${task.id}: ${_extractError(e)}`);
      }
    }

    for (const task of tasks) {
      const deps = task.dependencies || [];
      for (const depId of deps) {
        const taskBeadId = idMap.get(task.id);
        const depBeadId = idMap.get(depId);
        if (taskBeadId && depBeadId) {
          try { this._bd(`link ${taskBeadId} ${depBeadId}`); } catch (_) {}
        }
      }
    }

    return { migrated, skipped, errors };
  }

  async cleanupOrphans(proposalId) {
    try {
      const result = this._bd(`list --label proposal:${proposalId} --all --json`);
      const beads = JSON.parse(result);
      if (!Array.isArray(beads)) return { cleaned: 0 };

      let cleaned = 0;
      for (const b of beads) {
        const isEpic = (b.labels || []).some(l => l === 'type:epic');
        if (isEpic) continue;
        if (b.status === 'closed' || b.status === 'done' || b.status === 'cancelled') continue;
        try {
          this._bd(`close ${b.id} --reason "Orphan cleanup during v0.8→v0.9 migration"`);
          cleaned++;
        } catch (_) {}
      }
      return { cleaned };
    } catch (e) {
      return { cleaned: 0 };
    }
  }

  // ── 内部方法 ──

  _resolveEpic(proposalId) {
    if (this._epicCache.has(proposalId)) {
      return this._epicCache.get(proposalId);
    }
    try {
      const result = this._bd(`list --label proposal:${proposalId} --label type:epic --json`);
      const epics = JSON.parse(result);
      if (Array.isArray(epics) && epics.length > 0) {
        this._epicCache.set(proposalId, epics[0].id);
        return epics[0].id;
      }
    } catch (_) {}
    return null;
  }

  async _listTasks(proposalId) {
    try {
      const result = this._bd(`list --label proposal:${proposalId} --all --json`);
      const beads = JSON.parse(result);
      if (!Array.isArray(beads)) return [];
      return beads
        .filter(b => b.issue_type === 'task')
        .map(b => this._toTask(b));
    } catch (_) {
      return [];
    }
  }

  _getReadyFromList(proposalId) {
    try {
      const result = this._bd(`list --label proposal:${proposalId} --all --json`);
      const beads = JSON.parse(result);
      if (!Array.isArray(beads)) return [];
      const tasks = beads
        .filter(b => b.issue_type === 'task')
        .map(b => this._toTask(b));

      return tasks.filter(t => {
        if (t.status !== TASK_STATUS.PENDING) return false;
        const deps = t.dependencies || [];
        if (deps.length === 0) return true;
        return deps.every(depId => {
          const dep = tasks.find(dt => dt.id === depId);
          return dep && (dep.status === TASK_STATUS.DONE || dep.status === TASK_STATUS.CANCELLED);
        });
      });
    } catch (_) {
      return [];
    }
  }

  _findExistingBead(proposalId, title) {
    try {
      const result = this._bd(`list --label proposal:${proposalId} --all --json`);
      const beads = JSON.parse(result);
      if (!Array.isArray(beads)) return null;
      const match = beads.find(b => b.title === title);
      return match ? match.id : null;
    } catch (_) {
      return null;
    }
  }

  _findProposalDir(proposalId) {
    const config = loadMainConfig(this._projectRoot);
    if (!config) return null;
    return findProposalDir(this._projectRoot, config, proposalId);
  }

  _buildLabels(task, proposalId) {
    const labels = [
      `type:${task.type || TASK_TYPE.IMPLEMENTATION}`,
      `proposal:${proposalId}`,
    ];
    if (task.labels && task.labels.length > 0) {
      labels.push(...task.labels);
    }
    return labels.join(',');
  }

  _toTask(beadIssue) {
    const labels = beadIssue.labels || [];
    const typeLabel = labels.find(l => l.startsWith('type:'))?.replace('type:', '');
    const proposalLabel = labels.find(l => l.startsWith('proposal:'));
    return {
      id: beadIssue.id,
      title: beadIssue.title || '',
      description: beadIssue.description || '',
      type: typeLabel || TASK_TYPE.IMPLEMENTATION,
      status: _mapBeadStatus(beadIssue.status),
      dependencies: this._parseDeps(beadIssue),
      labels: labels.filter(l => !l.startsWith('type:') && !l.startsWith('proposal:')),
      claimedBy: beadIssue.assignee || '',
      blockReason: beadIssue.reason || '',
      summary: beadIssue.reason || '',
      createdAt: beadIssue.created_at || '',
      updatedAt: beadIssue.updated_at || '',
    };
  }

  _parseDeps(beadIssue) {
    const deps = [];
    const links = beadIssue.depends_on || [];
    for (const link of links) {
      if (typeof link === 'string' && !link.startsWith('discovered-from:')) {
        deps.push(link);
      }
    }
    return deps;
  }

  _bd(args) {
    return execSync(`bd ${args}`, {
      cwd: this._projectRoot,
      encoding: 'utf8',
      stdio: 'pipe',
      timeout: 30000,
      maxBuffer: 1024 * 1024
    }).trim();
  }

  _formatSnapshot(tasks, proposalId) {
    const now = new Date().toISOString().replace('T', ' ').slice(0, 19);
    const statusIcon = {
      pending: '⏳',
      in_progress: '🔄',
      done: '✅',
      blocked: '🚫',
      cancelled: '❌',
    };

    let md = `# Tasks — ${proposalId}\n`;
    md += `> Auto-generated at ${now}. Do not edit manually.\n\n`;
    md += `| Status | ID | Title | Type |\n`;
    md += `|--------|----|-------|------|\n`;

    const sorted = [...tasks].sort((a, b) => {
      const aParts = a.id.split('.');
      const bParts = b.id.split('.');
      for (let i = 0; i < Math.max(aParts.length, bParts.length); i++) {
        const aSeg = aParts[i] || '';
        const bSeg = bParts[i] || '';
        if (aSeg === bSeg) continue;
        const aNum = parseInt(aSeg, 10), bNum = parseInt(bSeg, 10);
        if (!isNaN(aNum) && !isNaN(bNum)) return aNum - bNum;
        if (!isNaN(aNum)) return -1;
        if (!isNaN(bNum)) return 1;
        return aSeg.localeCompare(bSeg);
      }
      return 0;
    });

    for (const t of sorted) {
      const depth = t.id.split('.').length - 1;
      const indent = '　'.repeat(depth);
      const icon = statusIcon[t.status] || '❓';
      md += `| ${icon} ${t.status} | ${indent}${t.id} | ${t.title} | ${t.type} |\n`;
    }

    const blocked = tasks.filter(t => t.status === TASK_STATUS.BLOCKED);
    if (blocked.length > 0) {
      md += `\n---\n`;
      for (const t of blocked) {
        md += `\n_Blocked: **${t.id}** — ${t.blockReason || 'No reason provided'}_\n`;
      }
    }

    return md;
  }
}

// ── 私有工具函数 ──

function _escape(str) {
  return String(str).replace(/"/g, '\\"').replace(/`/g, '\\`').replace(/\$/g, '\\$');
}

function _extractError(e) {
  if (e.stderr) return e.stderr.toString().trim().split('\n')[0];
  if (e.message) return e.message.split('\n')[0];
  return String(e);
}

function _mapBeadStatus(beadStatus) {
  const mapping = {
    'open': TASK_STATUS.PENDING,
    'pending': TASK_STATUS.PENDING,
    'in_progress': TASK_STATUS.IN_PROGRESS,
    'in progress': TASK_STATUS.IN_PROGRESS,
    'done': TASK_STATUS.DONE,
    'closed': TASK_STATUS.DONE,
    'completed': TASK_STATUS.DONE,
    'blocked': TASK_STATUS.BLOCKED,
    'cancelled': TASK_STATUS.CANCELLED,
  };
  return mapping[beadStatus] || TASK_STATUS.PENDING;
}

module.exports = { BeadsBackend };
