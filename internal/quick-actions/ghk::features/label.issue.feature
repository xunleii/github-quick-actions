@issue
Feature: add label with /label ~label [~label...] on issue description

  Background:
    Given quick action "/label" is registered for "issue" events

  @label
  Scenario: /label ~feature
    When Github sends an event "issue" with
      """
      {
        "action": "created",
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": {
          "body": "/label ~feature",
          "number": 1,
          "user": { "login":"xunleii" }
         },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/label" for "issue" event with arguments ["~feature"] by sending these following requests
      | API request method | API request URL                                                           | API request payload |
      | POST               | https://api.github.com/repos/xunleii/github-quick-actions/issues/1/labels | ["feature"]         |

  @label
  Scenario: /label ~feature ~bug:critical
    When Github sends an event "issue" with
      """
      {
        "action": "created",
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": {
          "body": "/label ~feature ~bug:critical",
          "number": 1,
          "user": { "login":"xunleii" }
         },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/label" for "issue" event with arguments ["~feature","~bug:critical"] by sending these following requests
      | API request method | API request URL                                                           | API request payload        |
      | POST               | https://api.github.com/repos/xunleii/github-quick-actions/issues/1/labels | ["feature","bug:critical"] |

  @label
  Scenario: /label ~feature feature
    When Github sends an event "issue" with
      """
      {
        "action": "created",
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": {
          "body": "/label ~feature feature",
          "number": 1,
          "user": { "login":"xunleii" }
         },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/label" for "issue" event with arguments ["~feature","feature"] by sending these following requests
      | API request method | API request URL                                                           | API request payload |
      | POST               | https://api.github.com/repos/xunleii/github-quick-actions/issues/1/labels | ["feature"]         |

  @label
  Scenario: /label without argument
    When Github sends an event "issue" with
      """
      {
        "action": "created",
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": {
          "body": "/label",
          "number": 1,
          "user": { "login":"xunleii" }
         },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/label" for "issue" event without argument without sending anything

  @label
  Scenario: /label ~feature on an invalid repository
    Given Github replies to 'POST https://api.github.com/repos/xunleii/github-quick-actions/issues/1/labels' with '404 {"message": "Not Found", "documentation_url": "https://docs.github.com/en/rest/reference/issues#add-labels-to-an-issue"}'
    When Github sends an event "issue" with
      """
      {
        "action": "created",
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": {
          "body": "/label ~feature",
          "number": 1,
          "user": { "login":"xunleii" }
         },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/label" for "issue" event with arguments ["~feature"] but returns this error: 'POST https://api.github.com/repos/xunleii/github-quick-actions/issues/1/labels: 404 Not Found []'
