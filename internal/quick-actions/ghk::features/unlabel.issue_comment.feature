Feature: remove label with /unlabel ~label [~label...] on issue comment event

  Background:
    Given quick action "/unlabel" is registered for "issue_comment" events

  @unlabel
  Scenario: /unlabel ~feature
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/unlabel ~feature" },
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 0 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/unlabel" for "issue_comment" event with arguments ["~feature"] by sending these following requests
      | API request method | API request URL                                                                   | API request payload |
      | DELETE             | https://api.github.com/repos/xunleii/github-quick-actions/issues/0/labels/feature |                     |

  @unlabel
  Scenario: /unlabel ~feature ~bug:critical
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/unlabel ~feature ~bug:critical" },
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 0 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/unlabel" for "issue_comment" event with arguments ["~feature","~bug:critical"] by sending these following requests
      | API request method | API request URL                                                                        | API request payload |
      | DELETE             | https://api.github.com/repos/xunleii/github-quick-actions/issues/0/labels/feature      |                     |
      | DELETE             | https://api.github.com/repos/xunleii/github-quick-actions/issues/0/labels/bug:critical |                     |

  @unlabel
  Scenario: /unlabel ~feature feature
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/unlabel ~feature feature", "user": { "login":"xunleii" }},
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 0 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/unlabel" for "issue_comment" event with arguments ["~feature","feature"] by sending these following requests
      | API request method | API request URL                                                                   | API request payload |
      | DELETE             | https://api.github.com/repos/xunleii/github-quick-actions/issues/0/labels/feature |                     |

  @unlabel
  Scenario: /unlabel all labels
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/unlabel", "user": { "login":"xunleii" }},
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 0 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/unlabel" for "issue_comment" event without argument by sending these following requests
      | API request method | API request URL                                                           | API request payload |
      | DELETE             | https://api.github.com/repos/xunleii/github-quick-actions/issues/0/labels |                     |

  @unlabel
  Scenario: /unlabel ~feature on an invalid repository
    Given Github replies to 'DELETE https://api.github.com/repos/xunleii/github-quick-actions/issues/0/labels/feature' with '404 {"message": "Not Found", "documentation_url": "https://docs.github.com/en/rest/reference/issues#add-labels-to-an-issue"}'
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/unlabel ~feature" },
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 0 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/unlabel" for "issue_comment" event with arguments ["~feature"] but returns this error: 'DELETE https://api.github.com/repos/xunleii/github-quick-actions/issues/0/labels/feature: 404 Not Found []'
