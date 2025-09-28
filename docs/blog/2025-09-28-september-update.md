---
slug: september-2025-update
title: Conveyor CI – September 2025 Progress Update
authors: jim-junior
tags: [conveyor-ci, update]
---

Hey, Here’s a quick summary of what’s been accomplished and what’s coming next.

<!-- truncate -->

## What’s New

- Added Embedded Log Storage
- Implemented Pipeline Engine Spec

## Next Steps

- Implement Comprehensive Security:
  - Ensure Conveyor CI passes [CNCF Day 0](https://github.com/cncf/toc/blob/main/toc_subprojects/project-reviews-subproject/general-technical-questions.md#day-0---planning-phase) security assessment.
  - Mandatory Access Control (MAC) Security via OS security frameworks like AppArmor, Seacomp and SE Linux. See [https://github.com/open-ug/conveyor/issues/43](https://github.com/open-ug/conveyor/issues/43)
  - End to End Security among System components i.e API Server and Driver Authentication and Authorisation.
- Declarative Configuration via `conveyor.yml` file. See [https://github.com/open-ug/conveyor/issues/43](https://github.com/open-ug/conveyor/issues/43)