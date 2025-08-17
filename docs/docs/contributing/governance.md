---
sidebar_position: 5
---

# Governance

**Conveyor CI** is dedicated to simplify development of CI/CD platforms in cloud-native environments. This governance document explains how the project is run.

- [Governance](#governance)
  - [Values](#values)
  - [Roles](#roles)
    - [Contributor](#contributor)
    - [Maintainer](#maintainer)
      - [Becoming a Maintainer](#becoming-a-maintainer)
      - [Removing a Maintainer](#removing-a-maintainer)
    - [Admin](#admin)
  - [Code of Conduct](#code-of-conduct)

## Values

**Conveyor CI** and its leadership embrace the following values:

* Openness: Communication and decision-making happens in the open and is
  discoverable for future reference. As much as possible, all discussions and
  work take place in public forums and open repositories.

* Fairness: All stakeholders have the opportunity to provide feedback and
  submit contributions, which will be considered on their merits.

* Community over Product or Company: Sustaining and growing our community takes
  priority over shipping code or sponsors' organizational goals.  Each
  contributor participates in the project as an individual.

* Inclusivity: We innovate through different perspectives and skill sets, which
  can only be accomplished in a welcoming and respectful environment.

* Participation: Responsibilities within the project are earned through
  participation, and there is a clear path up the contributor ladder into
  leadership positions.

## Roles

There are several roles relevant to **Conveyor CI**'s governance:

### Contributor

A Contributor to **Conveyor CI** is someone who has contributed to the project (e.g. code,
docs, CI) within the last 12 months. Contributors have read only access to the
Conveyor CI repositories on GitHub.

### Maintainer

**Conveyor CI** Maintainers (as defined by the [Conveyor's Maintainers team](./mantainers.md) ) have write access to the **Conveyor CI** repository on GitHub, which gives the ability to approve / merge / close PRs, trigger the CI and manage / classify project issues. All pull requests require review by a maintainer other than the one submitting it. "Large" changes are encouraged to gather consensus from multiple maintainers and interested community members. Maintainers are active Contributors and participants in the project, collectively managing the project's resources and contributors.

This privilege is granted with some expectation of responsibility: maintainers
are people who care about **Conveyor CI** and want to help it grow and improve. A
maintainer is not just someone who can make changes, but someone who has
demonstrated their ability to collaborate with the team, get the most
knowledgeable people to review code and docs, contribute high-quality code, and
follow through to fix issues (in code or tests).

A maintainer is a contributor to the project's success and a citizen helping
the project succeed.

The collective team of all Maintainers is known as the Maintainer Council, which
is the governing body for the project.

#### Becoming a Maintainer

To become a Maintainer you need to demonstrate the following:

  * commitment to the project:
    * actively participate in meetings, discussions, contributions, code and
      documentation reviews for at least 6 weeks,
    * contribute 3 non-trivial pull requests and have them merged,
  * ability to write quality code and/or documentation,
  * ability to collaborate with the team,
  * understanding of how the team works (policies, processes for testing and
    code review, etc),
  * understanding of the project's code base and coding / documentation
    style.

A new Maintainer must be proposed by an established Contributor and/or an
existing maintainer. A simple majority vote of existing Maintainers
approves the application. Maintainer nominations will be evaluated without
prejudice to employer or demographics.

Maintainers who are selected will be granted the necessary GitHub rights.

#### Removing a Maintainer

Maintainers may resign at any time if they feel that they will not be able to
continue fulfilling their project duties.

Maintainers may also be removed after being inactive, failure to fulfill their
Maintainer responsibilities, violating the Code of Conduct, or other reasons.
Inactivity is defined as a period of very low or no activity in the project for
6 months or more, with no definite schedule to return to full Maintainer
activity.

Depending on the reason for removal, a Maintainer may be converted to Emeritus
status. Emeritus Maintainers will still be consulted on some project matters,
and can be rapidly returned to Maintainer status if their availability changes.

### Admin

**Conveyor CI** Admins have admin access to the **Conveyor CI** repo, allowing them to do actions like, change the branch protection
rules for repositories, delete a repository and manage the access of others.
The Admin group is intentionally kept small, however, individuals can
be granted temporary admin access to carry out tasks, like creating a secret
that is used in a particular CI infrastructure.
The Admin list is reviewed and updated twice a year and typically contains:

- A subset of the maintainer team
- Optionally, some specific people that the Maintainers agree on adding for a
  specific purpose (e.g. to manage the CI)

## Code of Conduct

[Code of Conduct](https://github.com/open-ug/conveyor/blob/main/CODE_OF_CONDUCT.md) violations by community members will be
discussed and resolved by the maintainers privately. If a Maintainer is
directly involved in the report, the Maintainers will instead designate two
Maintainers to work with the CNCF Code of Conduct Committee in resolving it.
