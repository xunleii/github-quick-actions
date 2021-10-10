# Github quick actions
![GitHub deployments](https://img.shields.io/github/deployments/xunleii/github-quick-actions/AWS%20Lambda?label=Published%20on%20AWS%20Lambda)
[![Total alerts](https://img.shields.io/lgtm/alerts/g/xunleii/github-quick-actions.svg?logo=lgtm&logoWidth=18)](https://lgtm.com/projects/g/xunleii/github-quick-actions/alerts/)
[![GoReportCard example](https://goreportcard.com/badge/github.com/nanomsg/mangos)](https://goreportcard.com/report/github.com/xunleii/github-quick-actions)
[![CodeFactor](https://www.codefactor.io/repository/github/xunleii/github-quick-actions/badge/main)](https://www.codefactor.io/repository/github/xunleii/github-quick-actions/overview/main)
[![codecov](https://codecov.io/gh/xunleii/github-quick-actions/branch/main/graph/badge.svg?token=N69O0F7FGJ)](https://codecov.io/gh/xunleii/github-quick-actions)
[![GitHub license](https://img.shields.io/github/license/xunleii/github-quick-actions.svg)](https://github.com/xunleii/github-quick-actions/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/xunleii/github-quick-actions.svg)](https://GitHub.com/xunleii/github-quick-actions/releases/)

This workflow allows everyone to use [Gitlab quick actions](https://docs.gitlab.com/ee/user/project/quick_actions.html)
on their repository.

> NOTE: this documentation is mainly based on the Gitlab one, available
> on [Gitlab](https://gitlab.com/gitlab-org/gitlab/blob/master/doc/user/project/quick_actions.md)

## What are _quick actions_

Quick actions are text-based shortcuts for common actions that are usually done
by selecting buttons or dropdowns in the GitLab user interface. You can enter
these commands in the descriptions or comments of issues, epics, merge requests,
and commits.

Be sure to enter each quick action on a separate line to allow GitLab to
properly detect and execute the commands.

## Parameters

Many quick actions require a parameter. For example, the `/assign` quick action
requires a username.

If you manually enter a parameter, it must be enclosed in double quotation marks
(`"`), unless it contains only these characters:

- ASCII letters
- Numbers (0-9)
- Underscore (`_`), hyphen (`-`), question mark (`?`), dot (`.`), or ampersand (`&`)

Parameters are case-sensitive.

<!-- ::include quick_actions_table -->

## Contributing

See the [contributing guide](CONTRIBUTING.md) for detailed instructions of how to get
started with our project.

We accept different types of contributions, including some that don't require you to write
a single line of code. If you're looking for a way to contribute, you can scan through our
existing issues for something to work on.
