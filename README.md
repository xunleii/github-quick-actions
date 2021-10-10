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

## Available quick actions

The following quick actions are already released and available on the Github application.

|                               Command                               | Applicable on                |                                     Description                                     |
| :-----------------------------------------------------------------: | :--------------------------- | :---------------------------------------------------------------------------------: |
|                     `/assign @user [@user...]`                      | **&#10003;** `issue_comment` |           Assign one or more users.<br>_Use `me` to assign yourself._<br>           |
|                    `/unassign @user [@user...]`                     | **&#10003;** `issue_comment` |         Remove one or more assignees.<br>_Use `me` to remove yourself._<br>         |
|                     `/label ~label [~label...]`                     | **&#10003;** `issue_comment` | Add one or more labels.<br>_Label names can also start without a tilde (`~`)._<br>  |
|                    `/unlabel`<br>`/remove_label`                    | **&#10003;** `issue_comment` | Remove specified labels.<br>_Label names can also start without a tilde (`~`)._<br> |
| `/unlabel ~label [~label...]`<br>`/remove_label ~label [~label...]` | **&#10003;** `issue_comment` |                                 Remove all labels.                                  |

## Quick actions to be developed

The following quick actions will be available in the future (must need times to develop them).

|                               Command                               | Applicable on                                                                                          |                                                                                         Description                                                                                         |
| :-----------------------------------------------------------------: | :----------------------------------------------------------------------------------------------------- | :-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: |
|                     `/assign @user [@user...]`                      | **&#9676;** `issue`<br>**&#9676;** `pull_request`<br>**&#9676;** `pull_request_review_comment`         |                                                               Assign one or more users.<br>_Use `me` to assign yourself._<br>                                                               |
|                    `/reassign @user [@user...]`                     | **&#9676;** `issue_comment`<br>**&#9676;** `pull_request`<br>**&#9676;** `pull_request_review_comment` |                                                    Replace current assignees with those specified.<br>_Use `me` to assign yourself._<br>                                                    |
|                    `/unassign @user [@user...]`                     | **&#9676;** `pull_request`<br>**&#9676;** `pull_request_review_comment`                                |                                                             Remove one or more assignees.<br>_Use `me` to remove yourself._<br>                                                             |
|                 `/unassign`<br>`/remove_assignees`                  | **&#9676;** `pull_request`<br>**&#9676;** `pull_request_review_comment`                                |                                                                                    Remove all assignees.                                                                                    |
|                     `/label ~label [~label...]`                     | **&#9676;** `issue`<br>**&#9676;** `pull_request`<br>**&#9676;** `pull_request_review_comment`         |                                                     Add one or more labels.<br>_Label names can also start without a tilde (`~`)._<br>                                                      |
|                    `/relabel ~label [~label...]`                    | **&#9676;** `issue_comment`<br>**&#9676;** `pull_request`<br>**&#9676;** `pull_request_review_comment` |                                           Replace current labels with those specified.<br>_Label names can also start without a tilde (`~`)._<br>                                           |
| `/unlabel ~label [~label...]`<br>`/remove_label ~label [~label...]` | **&#9676;** `pull_request`<br>**&#9676;** `pull_request_review_comment`                                |                                                                                     Remove all labels.                                                                                      |
|                    `/unlabel`<br>`/remove_label`                    | **&#9676;** `pull_request`<br>**&#9676;** `pull_request_review_comment`                                |                                                     Remove specified labels.<br>_Label names can also start without a tilde (`~`)._<br>                                                     |
|                `/assign_reviewer @user [@user ...]`                 | **&#9676;** `issue_comment`<br>**&#9676;** `pull_request`<br>**&#9676;** `pull_request_review_comment` |                                                        Assign one or more users as reviewers.<br>_Use `me` to assign yourself._<br>                                                         |
|               `/reassign_reviewer @user [@user ...]`                | **&#9676;** `issue_comment`<br>**&#9676;** `pull_request`<br>**&#9676;** `pull_request_review_comment` |                                                    Replace current reviewers with those specified.<br>_Use `me` to assign yourself._<br>                                                    |
|               `/unassign_reviewer @user [@user ...]`                | **&#9676;** `issue_comment`<br>**&#9676;** `pull_request`<br>**&#9676;** `pull_request_review_comment` |                                                              Remove specified reviewers.<br>_Use `me` to remove yourself._<br>                                                              |
|             `/unassign_reviewer`<br>`/remove_reviewer`              | **&#9676;** `issue_comment`<br>**&#9676;** `pull_request`<br>**&#9676;** `pull_request_review_comment` |                                                                                    Remove all reviewers.                                                                                    |
|                              `/draft`                               | **&#9676;** `issue_comment`<br>**&#9676;** `pull_request`<br>**&#9676;** `pull_request_review_comment` |                                                                                  Toggle the draft status.                                                                                   |
|                              `/reopen`                              | **&#9676;** `issue_comment`                                                                            |                                                                          Reopen the current issue or pull request.                                                                          |
|                              `/close`                               | **&#9676;** `issue_comment`                                                                            |                                                                          Close the current issue or pull request.                                                                           |
|                              `/merge`                               | **&#9676;** `issue_comment`                                                                            |                                                                               Merge the current pull request.                                                                               |
|              `/copy_metadata #issue field [field...]`               | **&#9676;** `issue`<br>**&#9676;** `issue_comment`<br>**&#9676;** `pull_request`                       | Copy specified metadata from another issue or pull request.<br>_Available metadata are: assignees, reviewers, labels,<br>project, milestones, related_issues and related_pull_requests_<br> |
|                       `/copy_metadata #issue`                       | **&#9676;** `issue`<br>**&#9676;** `issue_comment`<br>**&#9676;** `pull_request`                       |                                                                    Copy all metadata from another issue or pull request.                                                                    |
|                 `/create_pull_request branch_name`                  | **&#9676;** `issue_comment`                                                                            |                              Create a new merge request starting from the current issue.<br>_It will automatically link the current issue with the new PR_<br>                              |
|                         `/duplicate #issue`                         | **&#9676;** `issue_comment`<br>**&#9676;** `pull_request`                                              |                                                                 Close this issue and mark as a duplicate of another issue.                                                                  |
|                       `/milestone %milestone`                       | **&#9676;** `issue`<br>**&#9676;** `issue_comment`<br>**&#9676;** `pull_request`                       |                                                                                       Set milestone.                                                                                        |
|                    `/relate #issue [#issue...]`                     | **&#9676;** `issue`<br>**&#9676;** `issue_comment`<br>**&#9676;** `pull_request`                       |                                                                                   Mark issues as related.                                                                                   |
|                    `/target_branch branch_name`                     | **&#9676;** `issue_comment`                                                                            |                                                                                     Set target branch.                                                                                      |
|                         `/title new_title`                          | **&#9676;** `issue_comment`<br>**&#9676;** `pull_request`                                              |                                                                                        Change title.                                                                                        |
|                  `/submit_review @user [@user...]`                  | **&#9676;** `issue_comment`                                                                            |                                                                       Submit a pending review to specified reviewers.                                                                       |
|                          `/submit_review`                           | **&#9676;** `issue_comment`                                                                            |                                                                          Submit a pending review to all reviewers.                                                                          |

## Quick actions that will not be developed

The following quick actions will not be developed for specific reasons.

|    Command     | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                               |                                                                      Reasons                                                                      |
| :------------: | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | :-----------------------------------------------------------------------------------------------------------------------------------------------: |
|  `/subscribe`  | Subscribe to notifications.                                                                                                                                                                                                                                                                                                                                                                                                                                               |                                     Github App needs access to a user information and ability to modify them                                      |
| `/unsubscribe` | Unsubscribe from notifications.                                                                                                                                                                                                                                                                                                                                                                                                                                           |                                     Github App needs access to a user information and ability to modify them                                      |
|   `/approve`   | Approve the merge request or the review.                                                                                                                                                                                                                                                                                                                                                                                                                                  |                                        Github App needs access to a user information and impersonate them                                         |
|   `/rebase`    | Rebase source branch.<br>This schedules a background task that attempts to rebase the changes in the source branch on the latest commit of the target branch.<br>If `/rebase` is used, `/merge` is ignored to avoid a race condition where the source branch is merged or deleted before it is rebased.<br>If there are merge conflicts, GitLab displays a message that a rebase cannot be scheduled.<br>Rebase failures are displayed with the merge request status.<br> | Cost too much to implement and to execute; alternative exists like using GithubAction with specific labels<br>Could have impact on the PR content |

## Contributing

See the [contributing guide](CONTRIBUTING.md) for detailed instructions of how to get
started with our project.

We accept different types of contributions, including some that don't require you to write
a single line of code. If you're looking for a way to contribute, you can scan through our
existing issues for something to work on.
