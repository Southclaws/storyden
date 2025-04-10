---
title: What is Storyden?
description: It's basically run-your-own-reddit, privately and securely.
---

Storyden is a modern forum, wiki and community hub. Like the internet forums of the past, but with a fresh coat of paint and modern security, deployment and intelligence features.

## What it isn't

Storyden isn't a help desk, customer support platform or headless CMS. There are some great products out there for these cases such as Discourse, Zendesk and Sanity. Storyden focuses on providing tooling for real internet communities, creators and curators. Those who find Discord too chaotic, Notion too locked-down and mediawiki too outdated. There are AI features too, which are purely optional.

## Philosophy

Storyden aims to be scalable from a tiny deployment with sane defaults to a large community with a lot to share.

Simplicity, privacy and security are at the core of Storyden's values.

### Simplicity

Sane defaults and zero dependencies. You can run a full production deployment using SQLite and the filesystem with no other dependencies forced on your setup. No Redis, PostgreSQL, S3, email servers, OAuth2 or other providers. While these can be enabled based on your needs, they aren't necessary for operating a production-grade installation of Storyden.

### Privacy

Storyden's frontend includes no cookie warning because it doesn't need one. Email-based login is opt-in if you want to allow membership without custodianship of people's personal information. Passwords are Argon2 hashed and no data ever leaves your installation.

### Security

Thanks to the above two values, security is made easy. All database queries run via ORM, all user-input is sanitised and functionality is rigorously end-to-end tested.

## Why Storyden

You've outgrown a WhatsApp group chat. You're sick of vulnerabilities in oldschool PHP forums. You're tired of sharing a spreadsheet for the car club. Storyden provides discussion, community knowledgebase and social bookmarking for anyone looking to upgrade their fan club, gaming group, clothing curation directory or whatever else you love to do with your people.

- Discuss topics in an oldschool-but-fresh forum interface.
- Curate and organise content in Notion-style databases with submission/review queues.
- Integrate with your favourite chat software so the conversation doesn't fragment.

## Who's behind it?

[Me!](https://barney.is/) I grew up on forums in the 2000s and Storyden is a love letter to the internet of my childhood.
