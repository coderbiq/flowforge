'use strict';

/**
 * TaskBackend — 任务存储后端接口契约
 *
 * 所有后端实现以此接口为准。BeadsBackend 是默认实现，通过 bd CLI 操作 beads。
 * 可替换为 Jira/Linear/自定义后端，只需实现相同接口。
 *
 * 每个方法首参为 proposalId（字符串），后端自行定位任务数据。
 * 不再暴露文件路径（proposalDir）或 YAML 操作（_readTaskMap / _writeTaskMap）。
 *
 * @interface
 */

// ============================================================
// 能力声明
// ============================================================
class BackendCapabilities {
  constructor(opts = {}) {
    /** 认领是否原子（多 Agent 安全） */
    this.atomicClaim = opts.atomicClaim ?? false;
    /** 是否支持 discovered-from 因果链 */
    this.discoveredFrom = opts.discoveredFrom ?? false;
    /** 是否支持审计日志 */
    this.auditTrail = opts.auditTrail ?? false;
    /** 是否支持依赖拓扑排序（就绪检测） */
    this.dependencySort = opts.dependencySort ?? false;
  }
}

// ============================================================
// 任务状态常量
// ============================================================
const TASK_STATUS = {
  PENDING: 'pending',
  IN_PROGRESS: 'in_progress',
  DONE: 'done',
  BLOCKED: 'blocked',
  CANCELLED: 'cancelled',
};

const TASK_TYPE = {
  ANALYSIS: 'analysis',
  DESIGN: 'design',
  IMPLEMENTATION: 'implementation',
};

// ============================================================
// TaskBackend 基类
// ============================================================
class TaskBackend {
  /**
   * @param {string} projectRoot — 项目根目录绝对路径
   */
  constructor(projectRoot) {
    this._projectRoot = projectRoot;
  }

  // ── 能力 ──

  /**
   * 返回后端能力声明。
   * @returns {BackendCapabilities}
   */
  getCapabilities() {
    return new BackendCapabilities();
  }

  /**
   * 检查后端在当前环境是否可用。
   * @returns {Promise<{available: boolean, reason?: string}>}
   */
  async checkAvailability() {
    return { available: true };
  }

  // ── 生命周期 ──

  /**
   * 为 proposal 初始化任务空间（创建 epic 等容器）。
   * 在 Design SKILL 阶段 5.2 首次创建任务时调用。
   *
   * @param {string} proposalId — CR-id，如 CR26060101
   * @param {string} title — proposal 标题
   * @returns {Promise<{epicId: string}>}
   */
  async init(proposalId, title) {
    throw new Error('Not implemented: init');
  }

  /**
   * 归档/清理 proposal 的任务空间（关闭 epic、清理关联 issue）。
   * 在 Archive SKILL 阶段 4 调用。
   *
   * @param {string} proposalId
   * @returns {Promise<void>}
   */
  async teardown(proposalId) {
    throw new Error('Not implemented: teardown');
  }

  // ── 任务创建 ──

  /**
   * 新增单个任务。
   *
   * @param {string} proposalId
   * @param {{title: string, description?: string, type?: string, dependencies?: string[], labels?: string[]}} task — 任务定义
   * @param {string} [parentTaskId] — 父任务 ID（子任务场景）
   * @returns {Promise<{taskId: string}>}
   */
  async addTask(proposalId, task, parentTaskId) {
    throw new Error('Not implemented: addTask');
  }

  /**
   * 批量新增任务。
   * 用于 Design SKILL 阶段 5.2 首次创建多个 analysis 任务。
   *
   * @param {string} proposalId
   * @param {Array<{title: string, description?: string, type?: string, dependencies?: string[], sourceTasks?: string[], epic?: string[]}>} tasks
   * @returns {Promise<{taskIds: string[]}>}
   */
  async addTasks(proposalId, tasks) {
    const ids = [];
    for (const task of tasks) {
      const result = await this.addTask(proposalId, task);
      ids.push(result.taskId);
    }
    return { taskIds: ids };
  }

  /**
   * 实施中发现的新任务，链接到父任务（discovered-from 因果链）。
   *
   * @param {string} proposalId
   * @param {string} parentTaskId — 发现此任务的父任务 ID
   * @param {{title: string, description?: string, type?: string}} task
   * @returns {Promise<{taskId: string}>}
   */
  async discoverTask(proposalId, parentTaskId, task) {
    throw new Error('Not implemented: discoverTask');
  }

  /**
   * 废弃任务（设计回退场景，不再执行）。
   * 状态变为 cancelled，不会出现在就绪列表中。
   *
   * @param {string} proposalId
   * @param {string} taskId
   * @param {string} [reason] — 废弃原因
   * @returns {Promise<void>}
   */
  async cancelTask(proposalId, taskId, reason) {
    throw new Error('Not implemented: cancelTask');
  }

  // ── 状态流转 ──

  /**
   * 认领任务（标记为 in_progress）。
   * 支持 atomicClaim 的后端保证原子性，多 Agent 安全。
   *
   * @param {string} proposalId
   * @param {string} taskId
   * @returns {Promise<{claimed: boolean, conflict?: string}>}
   *   claimed=false 时 conflict 说明原因（已被他人认领 / 状态不是 pending 等）
   */
  async claimTask(proposalId, taskId) {
    throw new Error('Not implemented: claimTask');
  }

  /**
   * 完成任务。
   *
   * @param {string} proposalId
   * @param {string} taskId
   * @param {string} [summary] — 完成摘要
   * @returns {Promise<void>}
   */
  async completeTask(proposalId, taskId, summary) {
    throw new Error('Not implemented: completeTask');
  }

  /**
   * 阻塞任务。
   *
   * @param {string} proposalId
   * @param {string} taskId
   * @param {string} reason — 阻塞原因
   * @returns {Promise<void>}
   */
  async blockTask(proposalId, taskId, reason) {
    throw new Error('Not implemented: blockTask');
  }

  /**
   * 释放认领（退回 pending）。
   *
   * @param {string} proposalId
   * @param {string} taskId
   * @returns {Promise<void>}
   */
  async unclaimTask(proposalId, taskId) {
    throw new Error('Not implemented: unclaimTask');
  }

  /**
   * 重开已完成的任务（design 回退场景）。
   *
   * @param {string} proposalId
   * @param {string} taskId
   * @returns {Promise<void>}
   */
  async reopenTask(proposalId, taskId) {
    throw new Error('Not implemented: reopenTask');
  }

  // ── 依赖管理 ──

  /**
   * 为任务添加依赖。
   *
   * @param {string} proposalId
   * @param {string} taskId — 被依赖方（需要等待的任务）
   * @param {string} dependsOnTaskId — 依赖方（前置任务）
   * @returns {Promise<void>}
   */
  async addDependency(proposalId, taskId, dependsOnTaskId) {
    throw new Error('Not implemented: addDependency');
  }

  /**
   * 移除依赖。
   *
   * @param {string} proposalId
   * @param {string} taskId
   * @param {string} dependsOnTaskId
   * @returns {Promise<void>}
   */
  async removeDependency(proposalId, taskId, dependsOnTaskId) {
    throw new Error('Not implemented: removeDependency');
  }

  // ── 查询 ──

  /**
   * 获取所有依赖已满足的 pending 任务。
   *
   * @param {string} proposalId
   * @returns {Promise<Array<{id: string, title: string, description: string, type: string, status: string, dependencies: string[], sourceTasks: string[], epic: string[]}>>}
   */
  async getReadyTasks(proposalId) {
    throw new Error('Not implemented: getReadyTasks');
  }

  /**
   * 获取完整状态概览。
   *
   * @param {string} proposalId
   * @returns {Promise<{total: number, byStatus: Record<string, number>, byType: Record<string, {total: number, done: number, inProgress: number, pending: number, blocked: number}>, tasks: Array}>}
   */
  async getStatus(proposalId) {
    throw new Error('Not implemented: getStatus');
  }

  /**
   * 获取所有阻塞任务。
   *
   * @param {string} proposalId
   * @returns {Promise<Array>}
   */
  async getBlockedTasks(proposalId) {
    const status = await this.getStatus(proposalId);
    return status.tasks.filter(t => t.status === TASK_STATUS.BLOCKED);
  }

  /**
   * 获取单任务详情。
   *
   * @param {string} proposalId
   * @param {string} taskId
   * @returns {Promise<object|null>}
   */
  async getTask(proposalId, taskId) {
    const status = await this.getStatus(proposalId);
    return status.tasks.find(t => t.id === taskId) || null;
  }

  /**
   * 检查 proposal 下是否所有任务已完成（归档前置条件）。
   *
   * @param {string} proposalId
   * @returns {Promise<boolean>}
   */
  async isAllDone(proposalId) {
    const status = await this.getStatus(proposalId);
    if (status.total === 0) return false;
    return status.tasks.every(t =>
      t.status === TASK_STATUS.DONE || t.status === TASK_STATUS.CANCELLED
    );
  }

  // ── 快照导出 ──

  /**
   * 导出人类可读快照（Markdown），写入 proposal 目录供 git 追踪和人工浏览。
   * 此文件由后端自动维护（可配合 hook 自动刷新），Agent 绝不操作它。
   *
   * @param {string} proposalId
   * @returns {Promise<string>} — 快照文件路径
   */
  async exportSnapshot(proposalId) {
    throw new Error('Not implemented: exportSnapshot');
  }

  // ── 迁移支持 ──

  /**
   * 从 task-map.yaml 迁移任务到后端。
   * 仅在升级场景使用（v0.8 → v0.9 升级）。
   *
   * @param {string} proposalDir — proposal 目录绝对路径
   * @returns {Promise<{migrated: number, skipped: number, errors: string[]}>}
   */
  async migrateFromYaml(proposalDir, proposalId) {
    throw new Error('Not implemented: migrateFromYaml');
  }

  /**
   * 清理后端中的孤儿 issue（proposal 标签相关但 task-map.yaml 中无对应）。
   * 仅在升级场景使用。
   *
   * @param {string} proposalId
   * @returns {Promise<{cleaned: number}>}
   */
  async cleanupOrphans(proposalId) {
    throw new Error('Not implemented: cleanupOrphans');
  }
}

module.exports = {
  TaskBackend,
  BackendCapabilities,
  TASK_STATUS,
  TASK_TYPE,
};
