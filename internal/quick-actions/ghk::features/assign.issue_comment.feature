@issue_comment
Feature: assign someone with /assign @user [@user...] on issue comment

  Background:
    Given quick action "/assign" is registered for "issue_comment" events

  @assign
  Scenario: /assign @mojombo
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/assign @mojombo", "user": { "login":"xunleii" }},
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 1 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/assign" for "issue_comment" event with arguments ["@mojombo"] by sending these following requests
      | API request method | API request URL                                                              | API request payload       |
      | POST               | https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees | {"assignees":["mojombo"]} |

  @assign
  Scenario: /assign @mojombo @defunkt
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/assign @mojombo @defunkt", "user": { "login":"xunleii" }},
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 1 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/assign" for "issue_comment" event with arguments ["@mojombo","@defunkt"] by sending these following requests
      | API request method | API request URL                                                              | API request payload                 |
      | POST               | https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees | {"assignees":["mojombo","defunkt"]} |

  @assign
  Scenario: /assign me
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/assign me", "user": { "login":"xunleii" }},
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 1 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/assign" for "issue_comment" event with arguments ["me"] by sending these following requests
      | API request method | API request URL                                                              | API request payload       |
      | POST               | https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees | {"assignees":["xunleii"]} |

  @assign
  Scenario: /assign @mojombo me
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/assign @mojombo me", "user": { "login":"xunleii" }},
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 1 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/assign" for "issue_comment" event with arguments ["@mojombo","me"] by sending these following requests
      | API request method | API request URL                                                              | API request payload                 |
      | POST               | https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees | {"assignees":["mojombo","xunleii"]} |

  @assign
  Scenario: /assign @mojombo @mojombo
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/assign @mojombo @mojombo", "user": { "login":"xunleii" }},
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 1 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/assign" for "issue_comment" event with arguments ["@mojombo","@mojombo"] by sending these following requests
      | API request method | API request URL                                                              | API request payload       |
      | POST               | https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees | {"assignees":["mojombo"]} |

  @assign
  Scenario: invalid /assign mojombo
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/assign mojombo", "user": { "login":"xunleii" }},
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 1 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/assign" for "issue_comment" event with arguments ["mojombo"] without sending anything

  @assign
  Scenario: /assign without argument
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/assign", "user": { "login":"xunleii" }},
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 1 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/assign" for "issue_comment" event without argument without sending anything

  @assign
  Scenario: /assign me on an invalid repository
    Given Github replies to 'POST https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees' with '404 {"message": "Not Found", "documentation_url": "https://docs.github.com/en/rest/reference/issues#add-labels-to-an-issue"}'
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/assign me", "user": { "login":"xunleii" }},
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 1 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/assign" for "issue_comment" event with arguments ["me"] but returns this error: 'POST https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees: 404 Not Found []'
