'use strict';

const { YamlAdapter } = require('./yaml');
const { BeadsAdapter } = require('./beads');

/**
 * 创建任务存储适配器。
 *
 * 根据 config.yaml 中 taskBackend.adapter 选择适配器：
 * - 'yaml'（默认）：纯 YAML 文件存储
 * - 'beads'：Beads 增强存储（task-map.yaml + beads 双写）
 *
 * @param {object} config - loadMainConfig() 返回的配置对象
 * @param {string} projectRoot - 项目根目录绝对路径
 * @returns {TaskAdapter}
 */
function createAdapter(config, projectRoot) {
  const adapterType = config?.taskBackend?.adapter || 'yaml';

  switch (adapterType) {
    case 'beads':
      return new BeadsAdapter(projectRoot);
    case 'yaml':
    default:
      return new YamlAdapter(projectRoot);
  }
}

module.exports = { createAdapter };
