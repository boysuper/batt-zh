# batt v0.10.0 · 社区汉化版（中/英双语）

> 基于 [@charlie0129](https://github.com/charlie0129) 的 [batt](https://github.com/charlie0129/batt) v0.10.0 源码构建，非官方汉化，遵循 GPL-2.0。
>
> 发布页：https://github.com/boysuper/batt-zh/releases/tag/v0.10.0-zh.1 · 上游告知 Issue：[#179](https://github.com/charlie0129/batt/issues/179)

## 本次内容

- 新增中/英运行时语言切换：**高级设置 → 语言 → 自动（跟随系统）/ 中文 / English**
- 切换语言**就地重建菜单**，进程不退出（修复旧方案「退出再 open」造成的闪退重启观感）
- 全量汉化：菜单项、状态行、确认弹窗、18 条悬停长提示、电源三行标签
- 图标为黑色底 + 底部红色「双语版」横幅（10 个尺寸全做）
- 未改动：充电控制、SMC 逻辑、daemon 行为、版本号 0.10.0、bundle id `cc.chlc.batt`

## 下载

GitHub 会在上传时去掉附件名中的中文，因此两个附件使用英文名。

| 文件 | 说明 |
| --- | --- |
| `batt-0.10.0-bilingual.dmg` | 双语安装包，卷内两个文件夹：**batt-中文版（汉化 i18n）** 与 **batt-原版（官方 0.10.0）**，各含一个 `batt.app`，互不冲突 |
| `batt-0.10.0-zh-source.zip` | 完整改造源码（可直接 `make app VERSION=v0.10.0` 构建），含 `pkg/gui/i18n.go` 词典 |

## 安装

1. 打开 DMG，进入 `batt-中文版（汉化 i18n）`，把 `batt.app` 拖进「应用程序」
2. 首次打开若被 Gatekeeper 拦下：右键 → 打开；或 `xattr -dr com.apple.quarantine /Applications/batt.app`
3. 汉化版为 ad-hoc 临时签名，未做 Apple 公证；要免提示需 Developer ID 签名 + 公证

## 已知说明

- 与上游同版本号、同 bundle id：`brew upgrade batt` 或安装官方新版会覆盖汉化
- 上游发新版后，把本仓库对 `pkg/gui/` 的改动重新应用到新源码即可

## 致谢

- 原作者：[@charlie0129](https://github.com/charlie0129)（batt）
- 图标与译文改动：本仓库维护者；欢迎提交更好的译法

---

# batt v0.10.0 · Unofficial Chinese localization

Built from [batt](https://github.com/charlie0129/batt) v0.10.0 by @charlie0129. GPL-2.0, not affiliated with upstream.

- Runtime **Auto / Chinese / English** switch under **Advanced → Language**
- Menu rebuilt **in place**; the process never exits
- Full localization incl. alerts, all 18 tooltips, power labels
- Icon: black tile with a red "双语版" banner along the bottom
- Charging logic, daemon, version 0.10.0 and bundle id unchanged

Assets: `batt-0.10.0-bilingual.dmg` (two folders: localized / official) and `batt-0.10.0-zh-source.zip`.
