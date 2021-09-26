Feature: handle Github event
  In order to manage quick actions on Github
  As Github application
  I need to handle correctly Github events

  Background:
    Given quick action "/hello_world" is registered for "issue_comment" events

  Scenario: should handle '/hello_world' command
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/hello_world" },
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 1 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/hello_world" for "issue_comment" event with no argument by sending these following requests
      | API request method | API request URL                                                           | API request payload |
      | POST               | https://api.github.com/repos/xunleii/github-quick-actions/issues/1/labels | ["hello_world"]     |

  Scenario: should handle multi '/hello_world' commands with arguments
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/hello_world\n/hello_world all \"all 'n' everyone\"\n/ignore_me" },
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 1 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/hello_world" for "issue_comment" event with no argument by sending these following requests
      | API request method | API request URL                                                           | API request payload |
      | POST               | https://api.github.com/repos/xunleii/github-quick-actions/issues/1/labels | ["hello_world"]     |
    And Github Quick Actions should handle command "/hello_world" for "issue_comment" event with arguments ["all","all 'n' everyone"] by sending these following requests
      | API request method | API request URL                                                           | API request payload                   |
      | POST               | https://api.github.com/repos/xunleii/github-quick-actions/issues/1/labels | ["hello_world@all","hello_world@all-n-everyone"] |

  Scenario: should handle '/hello_world' with error
    When Github replies to 'POST https://api.github.com/repos/xunleii/github-quick-actions/issues/1/labels' with '404 {"message": "Not Found", "documentation_url": "https://docs.github.com/en/rest/reference/issues#add-labels-to-an-issue"}'
    And Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/hello_world" },
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 1 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/hello_world" for "issue_comment" event with no argument but returns this error: 'POST https://api.github.com/repos/xunleii/github-quick-actions/issues/1/labels: 404 Not Found []'
