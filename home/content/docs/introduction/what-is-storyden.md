---
title: What is Storyden?
description: It's basically run-your-own-reddit, privately and securely.
---

Storyden is a modern forum, wiki and community hub. Like the internet forums of the past, but with a fresh coat of paint and modern security, deployment and intelligence features.

Storyden 是一个现代论坛、wiki 和社区枢纽。它保留了旧式互联网论坛的好东西，再加上更现代的安全、部署和智能能力。

## What it isn't

## 它不是什么

Storyden isn't a help desk, customer support platform or headless CMS. There are some great products out there for these cases such as Discourse, Zendesk and Sanity. Storyden focuses on providing tooling for real internet communities, creators and curators. Those who find Discord too chaotic, Notion too locked-down and mediawiki too outdated. There are AI features too, which are purely optional.

Storyden 不是帮助台、客服系统，也不是无头 CMS。这些场景已经有很优秀的产品，比如 Discourse、Zendesk 和 Sanity。Storyden 更关注真实互联网社区、创作者和内容策展者需要的工具：那些觉得 Discord 太吵、Notion 太封闭、MediaWiki 太老派的人。AI 功能也有，但完全是可选项。

## Philosophy

## 理念

Storyden aims to be scalable from a tiny deployment with sane defaults to a large community with a lot to share.

Storyden 希望从一个默认配置合理的小型部署，一路扩展到内容很多的大型社区。

Simplicity, privacy and security are at the core of Storyden's values.

简单、隐私和安全是 Storyden 的核心价值。

### Simplicity

### 简单

Sane defaults and zero dependencies. You can run a full production deployment using SQLite and the filesystem with no other dependencies forced on your setup. No Redis, PostgreSQL, S3, email servers, OAuth2 or other providers. While these can be enabled based on your needs, they aren't necessary for operating a production-grade installation of Storyden.

合理默认值，零强制依赖。你可以只用 SQLite 和文件系统跑起一套完整的生产部署，不会被迫安装 Redis、PostgreSQL、S3、邮件服务、OAuth2 或其他供应商服务。它们都可以按需启用，但不是运行 Storyden 的必需品。

### Privacy

### 隐私

Storyden's frontend includes no cookie warning because it doesn't need one. Email-based login is opt-in if you want to allow membership without custodianship of people's personal information. Passwords are Argon2 hashed and no data ever leaves your installation.

Storyden 的前端没有 cookie 弹窗，因为不需要。邮箱登录是可选的，如果你想让成员加入社区但不想托管他们的个人信息，可以用用户名优先的方式。密码使用 Argon2 哈希，数据也不会离开你的安装环境。

### Security

### 安全

Thanks to the above two values, security is made easy. All database queries run via ORM, all user-input is sanitised and functionality is rigorously end-to-end tested.

得益于上面的两个价值，安全也更容易做好。所有数据库查询都通过 ORM，所有用户输入都会被清洗，核心功能也经过端到端测试。

## Why Storyden

## 为什么选择 Storyden

You've outgrown a WhatsApp group chat. You're sick of vulnerabilities in oldschool PHP forums. You're tired of sharing a spreadsheet for the car club. Storyden provides discussion, community knowledgebase and social bookmarking for anyone looking to upgrade their fan club, gaming group, clothing curation directory or whatever else you love to do with your people.

你可能已经厌倦了 WhatsApp 群聊，受够了老式 PHP 论坛的安全问题，也不想再用共享表格管理车友会。Storyden 为想升级粉丝社群、游戏小组、服装资料目录，或任何你和伙伴们共同热爱的社区的人，提供讨论、社区知识库和社交书签能力。

- Discuss topics in an oldschool-but-fresh forum interface.
- Curate and organise content in Notion-style databases with submission/review queues.
- Integrate with your favourite chat software so the conversation doesn't fragment.

- 用老派但清爽的论坛界面讨论话题。
- 用类似 Notion 的数据库整理内容，并支持投稿与审核队列。
- 和你喜欢的聊天软件集成，让讨论不会四处碎掉。

## Who's behind it?

## 谁在做它？

[Me!](https://barney.is/) I grew up on forums in the 2000s and Storyden is a love letter to the internet of my childhood.

[我！](https://barney.is/) 我在 2000 年代的论坛里长大，Storyden 是写给童年互联网的一封情书。
