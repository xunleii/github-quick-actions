@pull_request_review_comment
Feature: remove label with /unlabel ~label [~label...] on pull request review comment

  Background:
    Given quick action "/unlabel" is registered for "pull_request_review_comment" events

  @unlabel
  Scenario: /unlabel ~feature
    When Github sends an event "pull_request_review_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/unlabel ~feature" },
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "pull_request": { "number": 1 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/unlabel" for "pull_request_review_comment" event with arguments ["~feature"] by sending these following requests
      | API request method | API request URL                                                                   | API request payload |
      | DELETE             | https://api.github.com/repos/xunleii/github-quick-actions/issues/1/labels/feature |                     |

  @unlabel
  Scenario: /unlabel ~feature ~bug:critical
    When Github sends an event "pull_request_review_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/unlabel ~feature ~bug:critical" },
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "pull_request": { "number": 1 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/unlabel" for "pull_request_review_comment" event with arguments ["~feature","~bug:critical"] by sending these following requests
      | API request method | API request URL                                                                        | API request payload |
      | DELETE             | https://api.github.com/repos/xunleii/github-quick-actions/issues/1/labels/feature      |                     |
      | DELETE             | https://api.github.com/repos/xunleii/github-quick-actions/issues/1/labels/bug:critical |                     |

  @unlabel
  Scenario: /unlabel ~feature feature
    When Github sends an event "pull_request_review_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/unlabel ~feature feature", "user": { "login":"xunleii" }},
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "pull_request": { "number": 1 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/unlabel" for "pull_request_review_comment" event with arguments ["~feature","feature"] by sending these following requests
      | API request method | API request URL                                                                   | API request payload |
      | DELETE             | https://api.github.com/repos/xunleii/github-quick-actions/issues/1/labels/feature |                     |

  @unlabel
  Scenario: /unlabel all labels
    When Github sends an event "pull_request_review_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/unlabel", "user": { "login":"xunleii" }},
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "pull_request": { "number": 1 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/unlabel" for "pull_request_review_comment" event without argument by sending these following requests
      | API request method | API request URL                                                           | API request payload |
      | DELETE             | https://api.github.com/repos/xunleii/github-quick-actions/issues/1/labels |                     |

  @unlabel @error
  Scenario: error handling on /unlabel
    Given Github replies to 'DELETE https://api.github.com/repos/xunleii/github-quick-actions/issues/1/labels/feature' with '404 {"message": "Not Found", "documentation_url": "https://docs.github.com/en/rest/reference/pull_requests#add-labels-to-an-pull_request"}'
    When Github sends an event "pull_request_review_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/unlabel ~feature" },
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "pull_request": { "number": 1 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/unlabel" for "pull_request_review_comment" event with arguments ["~feature"] but returns this error: 'DELETE https://api.github.com/repos/xunleii/github-quick-actions/issues/1/labels/feature: 404 Not Found []'
