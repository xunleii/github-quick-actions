Feature: add label with /label ~label [~label...] on issue comment event

  Background:
    Given quick action "/label" is registered for "issue_comment" events

  @label
  Scenario: /label ~feature
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/label ~feature" },
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 0 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/label" for "issue_comment" event with arguments ["~feature"] by sending these following requests
      | API request method | API request URL                                                           | API request payload |
      | POST               | https://api.github.com/repos/xunleii/github-quick-actions/issues/0/labels | ["feature"]}        |

  @label
  Scenario: /label ~feature ~bug:critical
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/label ~feature ~bug:critical" },
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 0 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/label" for "issue_comment" event with arguments ["~feature","~bug:critical"] by sending these following requests
      | API request method | API request URL                                                           | API request payload         |
      | POST               | https://api.github.com/repos/xunleii/github-quick-actions/issues/0/labels | ["feature","bug:critical"]} |

  @label
  Scenario: /label ~feature feature
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/label ~feature feature", "user": { "login":"xunleii" }},
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 0 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/label" for "issue_comment" event with arguments ["~feature","feature"] by sending these following requests
      | API request method | API request URL                                                           | API request payload |
      | POST               | https://api.github.com/repos/xunleii/github-quick-actions/issues/0/labels | ["feature"]}        |

  @label
  Scenario: /label without argument
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/label", "user": { "login":"xunleii" }},
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 0 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/label" for "issue_comment" event with no argument without sending anything

  @label
  Scenario: /label ~feature on an invalid repository
    Given Github replies to 'POST https://api.github.com/repos/xunleii/github-quick-actions/issues/0/labels' with '404 {"message": "Not Found", "documentation_url": "https://docs.github.com/en/rest/reference/issues#add-labels-to-an-issue"}'
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/label ~feature" },
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 0 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/label" for "issue_comment" event with arguments ["~feature"] but returns this error: 'POST https://api.github.com/repos/xunleii/github-quick-actions/issues/0/labels: 404 Not Found []'
