'use strict';

/**
 * 任务存储适配器接口契约
 *
 * 所有适配器方法以 proposalDir 为第一参数，自行决定如何定位和操作任务数据。
 * 可选增强方法均有默认降级实现，子类可按需覆盖。
 *
 * adapter 字段在 config.yaml 的 taskBackend.adapter 中配置。
 *
 * @interface
 */
class TaskAdapter {
  /**
   * 检查适配器在当前环境是否可用
   * @param {string} projectRoot - 项目根目录
   * @returns {Promise<{available: boolean, reason?: string}>}
   */
  async checkAvailability(projectRoot) {
    return { available: true };
  }

  // ========== 核心操作（所有适配器必须实现） ==========

  /**
   * 从 task-map.yaml 创建任务到后端。
   * 读取 proposalDir 下的 task-map.yaml，在后端中创建对应任务。
   * @param {string} proposalDir - proposal 目录绝对路径
   * @returns {Promise<{created: number, epicId?: string}>}
   */
  async createFromTaskMap(proposalDir) {
    throw new Error('Not implemented: createFromTaskMap');
  }

  /**
   * 获取就绪任务列表（依赖已满足且状态为 pending）。
   * @param {string} proposalDir
   * @returns {Promise<{tasks: Array<{id, title, description, deliverable}>, source: string}>}
   *   source: 'adapter' — 适配器原生能力（beads 拓扑排序）；'yaml' — 文本解析
   */
  async getReadyTasks(proposalDir) {
    throw new Error('Not implemented: getReadyTasks');
  }

  /**
   * 认领任务，标记为进行中。
   * @param {string} proposalDir
   * @param {string} taskId
   * @returns {Promise<{claimed: boolean, conflict?: string}>}
   *   claimed=false 时 conflict 说明原因（已被他人认领 / 状态不是 pending 等）
   */
  async claimTask(proposalDir, taskId) {
    throw new Error('Not implemented: claimTask');
  }

  /**
   * 完成任务。
   * @param {string} proposalDir
   * @param {string} taskId
   * @param {string} summary - 完成摘要
   * @returns {Promise<{done: boolean}>}
   */
  async completeTask(proposalDir, taskId, summary) {
    throw new Error('Not implemented: completeTask');
  }

  /**
   * 将任务标记为阻塞。
   * @param {string} proposalDir
   * @param {string} taskId
   * @param {string} reason - 阻塞原因
   * @returns {Promise<{blocked: boolean}>}
   */
  async blockTask(proposalDir, taskId, reason) {
    throw new Error('Not implemented: blockTask');
  }

  /**
   * 获取所有任务状态概览。
   * @param {string} proposalDir
   * @returns {Promise<{total: number, done: number, in_progress: number, pending: number, blocked: number, tasks: Array}>}
   */
  async getStatus(proposalDir) {
    throw new Error('Not implemented: getStatus');
  }

  /**
   * 归档前清理：校验任务完整性、关闭史诗等。
   * @param {string} proposalDir
   * @returns {Promise<{clean: boolean, issues?: string[]}>}
   *   clean=false 时 issues 列出未完成的任务/需处理的问题
   */
  async cleanup(proposalDir) {
    // 默认：检查 task-map.yaml 中是否有非 done 状态的任务
    const status = await this.getStatus(proposalDir);
    const open = status.tasks.filter(t => t.status !== 'done');
    if (open.length > 0) {
      return {
        clean: false,
        issues: open.map(t => `[${t.id}] ${t.title} (${t.status})`)
      };
    }
    return { clean: true };
  }

  /**
   * 向已有 task-map 增量添加任务（设计回退场景）。
   * 自动分配 ID，追加到任务列表末尾。
   * @param {string} proposalDir
   * @param {{title: string, description: string, deliverable?: string, dependencies?: string[]}} task
   * @returns {Promise<{added: boolean, taskId?: string}>}
   */
  async addTask(proposalDir, task) {
    const data = this._readTaskMap(proposalDir);
    if (!data) return { added: false };

    const newId = String(data.tasks.length + 1);
    const entry = {
      id: newId,
      title: task.title,
      description: task.description || '',
      deliverable: task.deliverable || '',
      status: 'pending',
      dependencies: task.dependencies || [],
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    };
    data.tasks.push(entry);
    this._writeTaskMap(proposalDir, data.tasks, data.proposal_id);
    return { added: true, taskId: newId };
  }

  /**
   * 废弃任务（设计回退场景，不再执行）。
   * 将任务状态设为 cancelled，不会出现在就绪列表中。
   * @param {string} proposalDir
   * @param {string} taskId
   * @param {string} reason - 废弃原因
   * @returns {Promise<{cancelled: boolean}>}
   */
  async cancelTask(proposalDir, taskId, reason) {
    const data = this._readTaskMap(proposalDir);
    if (!data) return { cancelled: false };

    const task = data.tasks.find(t => t.id === taskId);
    if (!task) return { cancelled: false };

    task.status = 'cancelled';
    task.updated_at = new Date().toISOString();
    if (reason) task.cancel_reason = reason;
    this._writeTaskMap(proposalDir, data.tasks, data.proposal_id);
    return { cancelled: true };
  }

  /**
   * 数据对账：修复 task-map.yaml 与后端之间的不一致。
   *
   * @param {string} proposalDir
   * @param {string} direction - 'yaml-to-backend' | 'backend-to-yaml' | 'merge'
   *   yaml-to-backend: 以 yaml 为源，覆盖后端（重新创建缺失任务、关闭多余任务）
   *   backend-to-yaml: 以后端为源，更新 yaml 状态
   *   merge（默认）: yaml 管定义（title/desc/deps），后端管状态（status/claim）
   * @returns {Promise<{synced: boolean, summary: {created: number, updated: number, closed: number, skipped: number}}>}
   */
  async sync(proposalDir, direction = 'merge') {
    return { synced: true, summary: { created: 0, updated: 0, closed: 0, skipped: 0 } };
  }

  // ========== 增强操作（可选实现，均有默认降级） ==========

  /**
   * 【增强】执行中发现新任务，链接因果关系。
   *
   * beads 适配器：使用 discovered-from 依赖类型。
   * yaml 适配器（默认降级）：追加到 task-map.yaml 末尾，无因果链。
   *
   * @param {string} proposalDir
   * @param {string} parentTaskId - 发现此任务的父任务 ID
   * @param {{title: string, description: string, deliverable?: string}} task
   * @returns {Promise<{created: boolean, taskId?: string}>}
   */
  async discoverTask(proposalDir, parentTaskId, task) {
    const yaml = require('../../vendor/js-yaml');
    const fs = require('fs');
    const path = require('path');
    const taskMapPath = path.join(proposalDir, 'task-map.yaml');
    if (!fs.existsSync(taskMapPath)) return { created: false };

    const data = yaml.load(fs.readFileSync(taskMapPath, 'utf8'));
    const tasks = data.tasks || [];
    const newId = String(tasks.length + 1);
    tasks.push({
      id: newId,
      title: task.title,
      description: task.description || '',
      deliverable: task.deliverable || '',
      status: 'pending',
      dependencies: [parentTaskId],
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    });
    data.tasks = tasks;
    fs.writeFileSync(taskMapPath, yaml.dump(data, { lineWidth: -1, noRefs: true }), 'utf8');
    return { created: true, taskId: newId };
  }

  /**
   * 【增强】获取工作流上下文，供 Agent 跨 session 恢复。
   *
   * beads 适配器：bd prime 输出。
   * yaml 适配器（默认降级）：返回 null，Agent 通过 implement-context 获取。
   *
   * @param {string} proposalDir
   * @returns {Promise<string|null>} 可注入 Agent 上下文的文本，null 表示无增强
   */
  async getContext(proposalDir) {
    return null;
  }

  /**
   * 【增强】释放认领，将任务退回 pending。
   * @param {string} proposalDir
   * @param {string} taskId
   * @returns {Promise<{released: boolean}>}
   */
  async unclaimTask(proposalDir, taskId) {
    const yaml = require('../../vendor/js-yaml');
    const fs = require('fs');
    const path = require('path');
    const taskMapPath = path.join(proposalDir, 'task-map.yaml');
    if (!fs.existsSync(taskMapPath)) return { released: false };

    const data = yaml.load(fs.readFileSync(taskMapPath, 'utf8'));
    const tasks = data.tasks || [];
    const task = tasks.find(t => t.id === taskId);
    if (!task || task.status !== 'in_progress') return { released: false };

    task.status = 'pending';
    task.updated_at = new Date().toISOString();
    fs.writeFileSync(taskMapPath, yaml.dump(data, { lineWidth: -1, noRefs: true }), 'utf8');
    return { released: true };
  }

  /**
   * 【增强】适配器能力声明。
   *
   * 供 SKILL / context 脚本根据适配器能力调整 Agent 行为：
   * - atomicClaim: 认领是否原子（多 Agent 安全）
   * - discoveredFrom: 是否支持发现式任务 + 因果链
   * - contextInjection: 是否支持 getContext 注入
   * - auditTrail: 是否有结构化审计日志
   *
   * @returns {{atomicClaim: boolean, discoveredFrom: boolean, contextInjection: boolean, auditTrail: boolean}}
   */
  getCapabilities() {
    return {
      atomicClaim: false,
      discoveredFrom: false,
      contextInjection: false,
      auditTrail: false,
    };
  }

  // ========== 共享工具方法（子类可用） ==========

  /**
   * 读取 task-map.yaml。
   * @param {string} proposalDir
   * @returns {{proposal_id?: string, tasks: Array}|null}
   */
  _readTaskMap(proposalDir) {
    const yaml = require('../../vendor/js-yaml');
    const fs = require('fs');
    const path = require('path');
    const taskMapPath = path.join(proposalDir, 'task-map.yaml');
    if (!fs.existsSync(taskMapPath)) return null;
    const data = yaml.load(fs.readFileSync(taskMapPath, 'utf8'));
    return { tasks: data?.tasks || [], proposal_id: data?.proposal_id };
  }

  /**
   * 写入 task-map.yaml。
   * @param {string} proposalDir
   * @param {Array} tasks
   * @param {string} [proposalId]
   */
  _writeTaskMap(proposalDir, tasks, proposalId) {
    const yaml = require('../../vendor/js-yaml');
    const fs = require('fs');
    const path = require('path');
    const taskMapPath = path.join(proposalDir, 'task-map.yaml');
    const data = { tasks };
    if (proposalId) data.proposal_id = proposalId;
    fs.writeFileSync(taskMapPath, yaml.dump(data, { lineWidth: -1, noRefs: true }), 'utf8');
  }

  /**
   * 从 task-map.yaml 中筛选就绪任务（依赖满足的 pending 任务）。
   * 作为不支持拓扑排序查询的后端的降级方案。
   * @param {Array} tasks
   * @returns {Array}
   */
  _filterReadyFromTasks(tasks) {
    return tasks.filter(t => {
      if (t.status !== 'pending') return false;
      const deps = t.dependencies || [];
      if (deps.length === 0) return true;
      return deps.every(depId => {
        const dep = tasks.find(dt => dt.id === depId);
        return dep && (dep.status === 'done' || dep.status === 'cancelled');
      });
    });
  }
}

module.exports = { TaskAdapter };
