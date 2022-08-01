# 4.0.0

## 重大变更

- 服务声明文件、基础设施信息文件、部署信息文件中不得存在值为空的项

## 功能特性

- 新增接口GetLatestConfigByEnvNum使外部应用能够获取某一环境号下最新生成的配置文件
- 配置模板新增预设函数GetIpByNet、GetNextTcpPort
- 为每个set生成两个额外的组内组播通道，分别命名为extra1、extra2
- 新增CHANGELOG.md文件

## 问题修复

无
