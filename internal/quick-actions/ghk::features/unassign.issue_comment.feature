Feature: unassign someone with /unassign @user [@user...] on issue comment event

  Background:
    Given quick action "/unassign" is registered for "issue_comment" events

  @unassign
  Scenario: /unassign @mojombo
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/unassign @mojombo", "user": { "login":"xunleii" }},
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 0 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/unassign" for "issue_comment" event with arguments ["@mojombo"] by sending these following requests
      | API request method | API request URL                                                              | API request payload       |
      | DELETE             | https://api.github.com/repos/xunleii/github-quick-actions/issues/0/assignees | {"assignees":["mojombo"]} |

  @unassign
  Scenario: /unassign @mojombo @defunkt
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/unassign @mojombo @defunkt", "user": { "login":"xunleii" }},
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 0 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/unassign" for "issue_comment" event with arguments ["@mojombo","@defunkt"] by sending these following requests
      | API request method | API request URL                                                              | API request payload                  |
      | DELETE             | https://api.github.com/repos/xunleii/github-quick-actions/issues/0/assignees | {"assignees":["mojombo", "defunkt"]} |

  @unassign
  Scenario: /unassign me
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/unassign me", "user": { "login":"xunleii" }},
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 0 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/unassign" for "issue_comment" event with arguments ["me"] by sending these following requests
      | API request method | API request URL                                                              | API request payload       |
      | DELETE             | https://api.github.com/repos/xunleii/github-quick-actions/issues/0/assignees | {"assignees":["xunleii"]} |

  @unassign
  Scenario: /unassign @mojombo me
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/unassign @mojombo me", "user": { "login":"xunleii" }},
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 0 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/unassign" for "issue_comment" event with arguments ["@mojombo","me"] by sending these following requests
      | API request method | API request URL                                                              | API request payload                  |
      | DELETE             | https://api.github.com/repos/xunleii/github-quick-actions/issues/0/assignees | {"assignees":["mojombo", "xunleii"]} |

  @unassign
  Scenario: /unassign @mojombo @mojombo
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/unassign @mojombo @mojombo", "user": { "login":"xunleii" }},
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 0 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/unassign" for "issue_comment" event with arguments ["@mojombo","@mojombo"] by sending these following requests
      | API request method | API request URL                                                              | API request payload       |
      | DELETE             | https://api.github.com/repos/xunleii/github-quick-actions/issues/0/assignees | {"assignees":["mojombo"]} |

  @unassign
  Scenario: invalid /unassign mojombo
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/unassign mojombo", "user": { "login":"xunleii" }},
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 0 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/unassign" for "issue_comment" event with arguments ["mojombo"] without sending anything

  @unassign
  Scenario: /unassign me on an invalid repository
    Given Github replies to 'DELETE https://api.github.com/repos/xunleii/github-quick-actions/issues/0/assignees' with '404 {"message": "Not Found", "documentation_url": "https://docs.github.com/en/rest/reference/issues#add-labels-to-an-issue"}'
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/unassign me", "user": { "login":"xunleii" }},
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 0 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/unassign" for "issue_comment" event with arguments ["me"] but returns this error: 'DELETE https://api.github.com/repos/xunleii/github-quick-actions/issues/0/assignees: 404 Not Found []'
