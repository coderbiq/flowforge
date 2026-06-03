'use strict';

const { execSync } = require('child_process');
const { TaskAdapter } = require('./interface');

/**
 * Beads 任务存储适配器。
 *
 * task-map.yaml 始终作为真理源。
 * beads 操作镜像到 YAML，bead_id 作为内部字段维护。
 * 查询优先使用 beads 能力（拓扑排序、就绪检测），YAML 作为降级。
 */
class BeadsAdapter extends TaskAdapter {
  constructor(projectRoot) {
    super();
    this._projectRoot = projectRoot;
  }

  // ========== 可用性检查 ==========

  async checkAvailability(_projectRoot) {
    try {
      execSync('bd context --json', {
        cwd: this._projectRoot,
        stdio: 'pipe',
        timeout: 5000
      });
      return { available: true };
    } catch (e) {
      return { available: false, reason: `bd not available: ${_extractError(e)}` };
    }
  }

  // ========== 核心操作 ==========

  async createFromTaskMap(proposalDir) {
    const data = this._readTaskMap(proposalDir);
    if (!data || data.tasks.length === 0) return { created: 0 };

    // 读取 proposal meta 获取标题
    const meta = this._readMeta(proposalDir);
    const proposalId = data.proposal_id || (meta ? meta.id : 'unknown');
    const title = meta ? meta.title : proposalId;

    // 创建 epic
    let epicId;
    try {
      const result = _bd(
        `create "${_escape(title)}" --type epic --labels proposal:${proposalId} --json`,
        this._projectRoot
      );
      const parsed = JSON.parse(result);
      epicId = parsed.id;
    } catch (e) {
      return { created: 0, error: `Failed to create epic: ${_extractError(e)}` };
    }

    // 为每个任务创建 beads issue
    let created = 0;
    for (const task of data.tasks) {
      try {
        const result = _bd(
          `create "${_escape(task.title)}" --type task --parent ${epicId} ` +
          `--labels task:${task.id},proposal:${proposalId} ` +
          `--description "${_escape(task.description || '')}" --json`,
          this._projectRoot
        );
        const parsed = JSON.parse(result);
        task._beadId = parsed.id;
        created++;
      } catch (e) {
        // 单个任务创建失败不阻塞其他任务
      }
    }

    // 创建依赖链接
    for (const task of data.tasks) {
      const deps = task.dependencies || [];
      for (const depId of deps) {
        const depTask = data.tasks.find(t => t.id === depId);
        if (depTask && depTask._beadId && task._beadId) {
          try {
            _bd(`link ${task._beadId} ${depTask._beadId}`, this._projectRoot);
          } catch (_) { /* 链接失败不阻塞 */ }
        }
      }
    }

    // 将 bead_id 写回 task-map.yaml
    this._writeTaskMap(proposalDir, data.tasks, data.proposal_id);

    return { created, epicId };
  }

  async getReadyTasks(proposalDir) {
    const data = this._readTaskMap(proposalDir);
    if (!data) return { tasks: [], source: 'yaml' };

    // 优先尝试 beads 查询
    try {
      const proposalId = data.proposal_id;
      if (proposalId) {
        const result = _bd(`ready --json`, this._projectRoot);
        const beadTasks = JSON.parse(result);
        if (Array.isArray(beadTasks) && beadTasks.length > 0) {
          // 用 beads 结果交叉匹配 YAML 任务
          const tasks = beadTasks.map(bt => {
            const yt = data.tasks.find(t => t._beadId === bt.id);
            if (yt) {
              return {
                id: yt.id,
                title: yt.title,
                description: yt.description,
                deliverable: yt.deliverable || ''
              };
            }
            // beads 中有但 YAML 中找不到的，作为一个引用返回
            return {
              id: bt.id,
              title: bt.title || '',
              description: bt.description || '',
              deliverable: '',
              _beadOnly: true
            };
          }).filter(t => t.title);
          return { tasks, source: 'beads' };
        }
      }
    } catch (_) { /* beads 查询失败，降级到 YAML */ }

    // 降级：YAML 文本解析
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

    // 尝试 beads 原子认领
    if (task._beadId) {
      try {
        _bd(`update ${task._beadId} --claim`, this._projectRoot);
      } catch (e) {
        return { claimed: false, conflict: `beads claim failed: ${_extractError(e)}` };
      }
    }

    // 更新 YAML
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

    // beads 关闭
    if (task._beadId) {
      try {
        const reason = summary ? ` --reason "${_escape(summary)}"` : '';
        _bd(`close ${task._beadId}${reason}`, this._projectRoot);
      } catch (_) { /* beads 关闭失败不阻塞 */ }
    }

    // 更新 YAML
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

    // beads 阻塞
    if (task._beadId) {
      try {
        _bd(`update ${task._beadId} --status blocked`, this._projectRoot);
      } catch (_) { /* beads 更新失败不阻塞 */ }
    }

    // 更新 YAML
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

  async cleanup(proposalDir) {
    const data = this._readTaskMap(proposalDir);
    if (!data) return { clean: true };

    // 通过 label 查询 beads 中是否有未关闭的任务
    try {
      const proposalId = data.proposal_id;
      if (proposalId) {
        const result = _bd(`list --label proposal:${proposalId} --all --json`, this._projectRoot);
        const beadTasks = JSON.parse(result);
        const openTasks = Array.isArray(beadTasks)
          ? beadTasks.filter(t => !['closed', 'done', 'completed'].includes(t.status))
          : [];

        if (openTasks.length > 0) {
          return {
            clean: false,
            issues: openTasks.map(t => `[${t.id}] ${t.title} (${t.status})`)
          };
        }
      }
    } catch (_) { /* beads 查询失败，降级到 YAML 检查 */ }

    // 降级：检查 YAML
    return super.cleanup(proposalDir);
  }

  async addTask(proposalDir, task) {
    const data = this._readTaskMap(proposalDir);
    if (!data) return { added: false };

    const newId = String(data.tasks.length + 1);
    const proposalId = data.proposal_id;

    let beadId = null;
    try {
      const meta = this._readMeta(proposalDir);
      const result = _bd(
        `create "${_escape(task.title)}" --type task ` +
        `--labels task:${newId},proposal:${proposalId || ''} ` +
        `--description "${_escape(task.description || '')}" --json`,
        this._projectRoot
      );
      beadId = JSON.parse(result).id;

      if (task.dependencies && task.dependencies.length > 0) {
        for (const depId of task.dependencies) {
          const depTask = data.tasks.find(t => t.id === depId);
          if (depTask && depTask._beadId) {
            try { _bd(`link ${beadId} ${depTask._beadId}`, this._projectRoot); } catch (_) {}
          }
        }
      }
    } catch (_) {}

    const entry = {
      id: newId,
      title: task.title,
      description: task.description || '',
      deliverable: task.deliverable || '',
      status: 'pending',
      dependencies: task.dependencies || [],
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      _beadId: beadId
    };
    data.tasks.push(entry);
    this._writeTaskMap(proposalDir, data.tasks, data.proposal_id);
    return { added: true, taskId: newId };
  }

  async cancelTask(proposalDir, taskId, reason) {
    const data = this._readTaskMap(proposalDir);
    if (!data) return { cancelled: false };

    const task = data.tasks.find(t => t.id === taskId);
    if (!task) return { cancelled: false };

    if (task._beadId) {
      try {
        const closeReason = reason ? ` --reason "${_escape(reason)}"` : '';
        _bd(`close ${task._beadId}${closeReason}`, this._projectRoot);
      } catch (_) {}
    }

    task.status = 'cancelled';
    task.updated_at = new Date().toISOString();
    if (reason) task.cancel_reason = reason;
    this._writeTaskMap(proposalDir, data.tasks, data.proposal_id);
    return { cancelled: true };
  }

  async sync(proposalDir, direction = 'merge') {
    const data = this._readTaskMap(proposalDir);
    if (!data) return { synced: false, summary: { created: 0, updated: 0, closed: 0, skipped: 0 } };

    const proposalId = data.proposal_id;
    const yamlTasks = data.tasks;
    const summary = { created: 0, updated: 0, closed: 0, skipped: 0 };

    let beadTasks = [];
    try {
      if (proposalId) {
        const result = _bd(`list --label proposal:${proposalId} --all --json`, this._projectRoot);
        beadTasks = JSON.parse(result);
        if (!Array.isArray(beadTasks)) beadTasks = [];
      }
    } catch (_) {
      beadTasks = [];
    }

    const beadById = {};
    for (const bt of beadTasks) beadById[bt.id] = bt;

    if (direction === 'yaml-to-beads') {
      for (const yt of yamlTasks) {
        if (yt._beadId && beadById[yt._beadId]) {
          summary.skipped++;
          continue;
        }
        try {
          const result = _bd(
            `create "${_escape(yt.title)}" --type task ` +
            `--labels task:${yt.id},proposal:${proposalId || ''} ` +
            `--description "${_escape(yt.description || '')}" --json`,
            this._projectRoot
          );
          yt._beadId = JSON.parse(result).id;
          summary.created++;
        } catch (_) {
          summary.skipped++;
        }
      }

      for (const bt of beadTasks) {
        const match = yamlTasks.find(yt => yt._beadId === bt.id);
        if (!match) {
          try {
            _bd(`close ${bt.id} --reason "Not in task-map.yaml"`, this._projectRoot);
            summary.closed++;
          } catch (_) {}
        }
      }
      this._writeTaskMap(proposalDir, yamlTasks, proposalId);

    } else if (direction === 'beads-to-yaml') {
      for (const yt of yamlTasks) {
        if (!yt._beadId) { summary.skipped++; continue; }
        const bt = beadById[yt._beadId];
        if (!bt) { summary.skipped++; continue; }

        const beadStatus = _mapBeadStatus(bt.status);
        if (beadStatus && yt.status !== beadStatus) {
          yt.status = beadStatus;
          yt.updated_at = new Date().toISOString();
          summary.updated++;
        } else {
          summary.skipped++;
        }
      }
      this._writeTaskMap(proposalDir, yamlTasks, proposalId);

    } else {
      for (const yt of yamlTasks) {
        if (!yt._beadId) {
          try {
            const result = _bd(
              `create "${_escape(yt.title)}" --type task ` +
              `--labels task:${yt.id},proposal:${proposalId || ''} ` +
              `--description "${_escape(yt.description || '')}" --json`,
              this._projectRoot
            );
            yt._beadId = JSON.parse(result).id;
            summary.created++;
          } catch (_) {
            summary.skipped++;
          }
          continue;
        }

        const bt = beadById[yt._beadId];
        if (!bt) {
          try {
            const result = _bd(
              `create "${_escape(yt.title)}" --type task ` +
              `--labels task:${yt.id},proposal:${proposalId || ''} ` +
              `--description "${_escape(yt.description || '')}" --json`,
              this._projectRoot
            );
            yt._beadId = JSON.parse(result).id;
            summary.created++;
          } catch (_) {
            summary.skipped++;
          }
          continue;
        }

        const beadStatus = _mapBeadStatus(bt.status);
        if (beadStatus && yt.status !== beadStatus) {
          yt.status = beadStatus;
          yt.updated_at = new Date().toISOString();
          summary.updated++;
        } else {
          summary.skipped++;
        }
      }
      this._writeTaskMap(proposalDir, yamlTasks, proposalId);
    }

    return { synced: true, summary };
  }

  // ========== 增强操作 ==========

  async discoverTask(proposalDir, parentTaskId, task) {
    const data = this._readTaskMap(proposalDir);
    if (!data) return { created: false };

    const parent = data.tasks.find(t => t.id === parentTaskId);
    const newId = String(data.tasks.length + 1);

    // beads: 使用 discovered-from 依赖
    let beadId = null;
    if (parent && parent._beadId) {
      try {
        const result = _bd(
          `create "${_escape(task.title)}" --type task ` +
          `--dep discovered-from:${parent._beadId} --json`,
          this._projectRoot
        );
        beadId = JSON.parse(result).id;
      } catch (_) { /* beads 创建失败，仍然追加到 YAML */ }
    }

    // 追加到 YAML
    data.tasks.push({
      id: newId,
      title: task.title,
      description: task.description || '',
      deliverable: task.deliverable || '',
      status: 'pending',
      dependencies: [parentTaskId],
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      _beadId: beadId
    });
    this._writeTaskMap(proposalDir, data.tasks, data.proposal_id);
    return { created: true, taskId: newId };
  }

  async getContext(_proposalDir) {
    try {
      return _bd('prime', this._projectRoot);
    } catch (_) {
      return null;
    }
  }

  async unclaimTask(proposalDir, taskId) {
    const data = this._readTaskMap(proposalDir);
    if (!data) return { released: false };

    const task = data.tasks.find(t => t.id === taskId);
    if (!task || task.status !== 'in_progress') return { released: false };

    // beads 释放
    if (task._beadId) {
      try {
        _bd(`update ${task._beadId} --status pending`, this._projectRoot);
      } catch (_) { /* beads 更新失败不阻塞 */ }
    }

    task.status = 'pending';
    task.updated_at = new Date().toISOString();
    this._writeTaskMap(proposalDir, data.tasks, data.proposal_id);
    return { released: true };
  }

  getCapabilities() {
    return {
      atomicClaim: true,
      discoveredFrom: true,
      contextInjection: true,
      auditTrail: true,
    };
  }

  // ========== 内部方法 ==========

  _readMeta(proposalDir) {
    const yaml = require('../vendor/js-yaml');
    const fs = require('fs');
    const path = require('path');
    const metaPath = path.join(proposalDir, 'meta.yaml');
    if (!fs.existsSync(metaPath)) return null;
    return yaml.load(fs.readFileSync(metaPath, 'utf8'));
  }
}

// ========== 私有工具函数 ==========

function _bd(args, cwd) {
  return execSync(`bd ${args}`, {
    cwd,
    encoding: 'utf8',
    stdio: 'pipe',
    timeout: 10000,
    maxBuffer: 1024 * 1024
  }).trim();
}

function _escape(str) {
  return str.replace(/"/g, '\\"').replace(/`/g, '\\`').replace(/\$/g, '\\$');
}

function _extractError(e) {
  if (e.stderr) return e.stderr.toString().trim().split('\n')[0];
  if (e.message) return e.message.split('\n')[0];
  return String(e);
}

function _mapBeadStatus(beadStatus) {
  const mapping = {
    'open': 'pending',
    'pending': 'pending',
    'in_progress': 'in_progress',
    'in progress': 'in_progress',
    'done': 'done',
    'closed': 'done',
    'completed': 'done',
    'blocked': 'blocked',
    'cancelled': 'cancelled',
  };
  return mapping[beadStatus] || null;
}

module.exports = { BeadsAdapter };
