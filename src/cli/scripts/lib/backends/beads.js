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

  /**
   * 初始化 proposal 的任务空间，创建 4 层任务层级：
   *
   *   Main Epic: "CRID: Proposal Title"
   *   ├── Sub-Epic: "CRID: 分析"
   *   │   └── Task (analysis)
   *   ├── Sub-Epic: "CRID: 设计"
   *   │   └── Task (design)
   *   └── Sub-Epic: "CRID: 实施"
   *       ├── Parent Task (implementation)
   *       │   └── Child Task (implementation)
   *       └── Standalone Task (implementation)
   *
   * 层级说明见 guides/task-hierarchy.md
   */
  async init(proposalId, title) {
    const existingId = await this._resolveEpic(proposalId);
    if (existingId) {
      // 重建：关闭旧子 epic 和旧主 epic，清除缓存后重建
      for (const type of ['analysis', 'design', 'implementation']) {
        const subId = await this._resolveSubEpic(proposalId, type);
        if (subId) {
          try { this._bd(`close ${subId} --reason "Re-initialized"`); } catch (_) {}
        }
      }
      try { this._bd(`close ${existingId} --reason "Re-initialized"`); } catch (_) {}
      this._epicCache.delete(proposalId);
    }

    // 1. 创建主 epic（格式: "CRID: Proposal Title"）
    const epicTitle = `${proposalId}: ${title}`;
    const result = this._bd(
      `create "${_escape(epicTitle)}" --type epic ` +
      `--labels type:epic,type:main-epic,proposal:${proposalId} --json`
    );
    const epic = JSON.parse(result);
    const mainEpicId = epic.id;
    this._epicCache.set(proposalId, mainEpicId);

    // 2. 创建 3 个类型子 epic（挂在主 epic 下）
    const TYPE_SUB_EPICS = {
      analysis: '分析',
      design: '设计',
      implementation: '实施',
    };
    const subEpics = {};
    for (const [type, label] of Object.entries(TYPE_SUB_EPICS)) {
      const subTitle = `${proposalId}: ${label}`;
      const subResult = this._bd(
        `create "${_escape(subTitle)}" --type epic --parent ${mainEpicId} ` +
        `--labels type:epic,type:sub-epic,type:${type},proposal:${proposalId} --json`
      );
      subEpics[type] = JSON.parse(subResult).id;
    }

    return { epicId: mainEpicId, subEpics };
  }

  async hasTaskSpace(proposalId) {
    const epicId = await this._resolveEpic(proposalId);
    return epicId !== null;
  }

  async teardown(proposalId) {
    try {
      // 关闭所有子 epic
      for (const type of ['analysis', 'design', 'implementation']) {
        const subId = await this._resolveSubEpic(proposalId, type);
        if (subId) {
          try { this._bd(`close ${subId} --reason "Proposal archived"`); } catch (_) {}
        }
      }
      // 关闭主 epic
      const epicId = await this._resolveEpic(proposalId);
      if (epicId) {
        this._bd(`close ${epicId} --reason "Proposal archived"`);
      }
    } catch (_) { /* 清理失败不阻塞归档 */ }
    this._epicCache.delete(proposalId);
  }

  // ── 任务创建 ──

  async addTask(proposalId, task, parentTaskId) {
    const subEpicId = parentTaskId || await this._resolveSubEpic(proposalId, task.type);
    if (!subEpicId) {
      throw new Error(`No epic found for proposal ${proposalId}. Run 'flowforge task init --proposal ${proposalId} "<title>"' first.`);
    }
    const labels = this._buildLabels(task, proposalId);
    const result = this._bd(
      `create "${_escape(task.title)}" --type task --parent ${subEpicId} ` +
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
    const ids = [];

    for (const t of tasks) {
      const subEpicId = await this._resolveSubEpic(proposalId, t.type);
      if (!subEpicId) {
        throw new Error(`No epic found for proposal ${proposalId}. Run 'flowforge task init' first.`);
      }
      const labels = this._buildLabels(t, proposalId);
      const result = this._bd(
        `create "${_escape(t.title)}" --type task --parent ${subEpicId} ` +
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
    const subEpicId = await this._resolveSubEpic(proposalId, task.type);
    if (!subEpicId) {
      throw new Error(`No epic found for proposal ${proposalId}. Run 'flowforge task init' first.`);
    }
    const labels = this._buildLabels(task, proposalId);
    const depFlag = `--dep discovered-from:${parentTaskId}`;
    const result = this._bd(
      `create "${_escape(task.title)}" --type task --parent ${parentTaskId} ` +
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

  async getTask(proposalId, taskId) {
    const tasks = await this._listTasks(proposalId);
    return tasks.find(t => t.id === taskId) || null;
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

  async migrateFromYaml(proposalDir, proposalId) {
    const yaml = require('../../vendor/js-yaml');
    const taskMapPath = path.join(proposalDir, 'task-map.yaml');
    if (!fs.existsSync(taskMapPath)) {
      return { migrated: 0, skipped: 0, errors: ['task-map.yaml not found'] };
    }

    const data = yaml.load(fs.readFileSync(taskMapPath, 'utf8'));
    const tasks = data?.tasks || [];
    const pid = proposalId || data?.proposal_id || path.basename(proposalDir);

    let epicId;
    try {
      const result = this._bd(
        `create "Proposal: ${_escape(pid)}" --type epic --labels proposal:${pid} --json`
      );
      epicId = JSON.parse(result).id;
      this._epicCache.set(pid, epicId);
    } catch (e) {
      return { migrated: 0, skipped: 0, errors: [`Failed to create epic: ${_extractError(e)}`] };
    }

    let migrated = 0;
    let skipped = 0;
    const errors = [];
    const idMap = new Map();

    for (const task of tasks) {
      try {
        const existing = this._findExistingBead(pid, task.title);
        if (existing) {
          idMap.set(task.id, existing);
          skipped++;
          continue;
        }

        const labels = this._buildLabels(task, pid);
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
        const labels = b.labels || [];
        const isEpic = labels.some(l => l === 'type:epic' || l === 'type:main-epic' || l === 'type:sub-epic');
        if (isEpic) continue;
        if (b.status === 'closed' || b.status === 'done' || b.status === 'cancelled') continue;
        try {
          this._bd(`close ${b.id} --reason "Orphan cleanup"`);
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
      const result = this._bd(
        `list --label proposal:${proposalId} --label type:main-epic --json`
      );
      const epics = JSON.parse(result);
      if (Array.isArray(epics) && epics.length > 0) {
        this._epicCache.set(proposalId, epics[0].id);
        return epics[0].id;
      }
    } catch (_) {}
    return null;
  }

  async _resolveSubEpic(proposalId, type) {
    const typeLabel = `type:${type}`;
    try {
      const result = this._bd(
        `list --label proposal:${proposalId} --label type:sub-epic --label ${typeLabel} --json`
      );
      const issues = JSON.parse(result);
      if (Array.isArray(issues)) {
        const subEpic = issues.find(b => b.issue_type === 'epic');
        if (subEpic) return subEpic.id;
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
    // beads 会将父 epic 标签继承到子任务，需过滤出实际 task type
    const VALID_TASK_TYPES = new Set(['analysis', 'design', 'implementation']);
    const typeLabel = labels.find(l => {
      const t = l.replace('type:', '');
      return VALID_TASK_TYPES.has(t);
    })?.replace('type:', '');
    const proposalLabel = labels.find(l => l.startsWith('proposal:'));
    // beads 的父子关系通过 ID 前缀表达，depends_on 不包含 parent-child
    const idParts = beadIssue.id.split('.');
    const parentId = idParts.length > 1 ? idParts.slice(0, -1).join('.') : null;
    return {
      id: beadIssue.id,
      title: beadIssue.title || '',
      description: beadIssue.description || '',
      type: typeLabel || TASK_TYPE.IMPLEMENTATION,
      status: _mapBeadStatus(beadIssue.status),
      parentId,
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
    // 写操作自动加 --sandbox 避免 dolt auto-push 阻塞，Agent 无需感知
    const needsSandbox = /^(create|update|close|link|unlink|label)\b/.test(args);
    const sandboxFlag = needsSandbox ? ' --sandbox' : '';
    const result = execSync(`bd ${args}${sandboxFlag}`, {
      cwd: this._projectRoot,
      encoding: 'utf8',
      stdio: 'pipe',
      timeout: 30000,
      maxBuffer: 1024 * 1024
    }).trim();
    return result;
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

    // 按类型分组
    const typeOrder = ['analysis', 'design', 'implementation'];
    const typeLabels = { analysis: '分析', design: '设计', implementation: '实施' };

    for (const type of typeOrder) {
      const typeTasks = tasks.filter(t => t.type === type);
      if (typeTasks.length === 0) continue;

      const totalByStatus = {};
      for (const t of typeTasks) {
        totalByStatus[t.status] = (totalByStatus[t.status] || 0) + 1;
      }
      const doneCount = totalByStatus.done || 0;
      md += `## ${typeLabels[type]} (${typeTasks.length} tasks, ${doneCount} done)\n\n`;
      md += `| Status | ID | Title |\n`;
      md += `|--------|----|-------|\n`;

      // 构建父子树：仅当 parentId 是另一个 task 时才作为子节点
      const childrenMap = new Map();
      const roots = [];
      const taskIds = new Set(typeTasks.map(t => t.id));
      for (const t of typeTasks) {
        if (t.parentId && taskIds.has(t.parentId)) {
          const siblings = childrenMap.get(t.parentId) || [];
          siblings.push(t);
          childrenMap.set(t.parentId, siblings);
        } else {
          roots.push(t);
        }
        childrenMap.set(t.id, childrenMap.get(t.id) || []);
      }

      const printed = new Set();
      const renderTask = (task, depth) => {
        if (printed.has(task.id)) return;
        printed.add(task.id);
        const indent = '　'.repeat(depth);
        const icon = statusIcon[task.status] || '❓';
        md += `| ${icon} ${task.status} | ${indent}${task.id} | ${indent}${task.title} |\n`;
        const children = childrenMap.get(task.id) || [];
        const sorted = [...children].sort((a, b) => a.id.localeCompare(b.id));
        for (const child of sorted) {
          renderTask(child, depth + 1);
        }
      };

      // 根任务按 ID 排序
      const sortedRoots = [...roots].sort((a, b) => a.id.localeCompare(b.id));
      for (const root of sortedRoots) {
        renderTask(root, 0);
      }

      md += `\n`;
    }

    const blocked = tasks.filter(t => t.status === TASK_STATUS.BLOCKED);
    if (blocked.length > 0) {
      md += `---\n`;
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
