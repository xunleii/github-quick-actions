@issue
Feature: assign someone with /duplicate #issue [#issue...] on issue description

  Background:
    Given quick action "/duplicate" is registered for "issue" events

  @duplicate
  Scenario: /duplicate #1
    # Comment ID required to remove the comment
    Given Github replies to 'POST https://api.github.com/repos/xunleii/github-quick-actions/issues/2/comments' with '200 {"id": 1234}'

    When Github sends an event "issue" with
      """
      {
        "action": "created",
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": {
          "body": "/duplicate #1",
          "number": 2,
          "user": { "login":"xunleii" }
        },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/duplicate" for "issue" event with arguments ["#1"] by sending these following requests
      | API request method | API request URL                                                                | API request payload        |
      | GET                | https://api.github.com/repos/xunleii/github-quick-actions/issues/1             |                            |
      | POST               | https://api.github.com/repos/xunleii/github-quick-actions/issues/2/comments    | {"body":"Duplicate of #1"} |
      | DELETE             | https://api.github.com/repos/xunleii/github-quick-actions/issues/comments/1234 |                            |

  @duplicate
  Scenario: /duplicate #1 #2
    # Comment ID required to remove the comment
    Given Github replies to 'POST https://api.github.com/repos/xunleii/github-quick-actions/issues/3/comments' with '200 {"id": 1234}'
    Given Github replies to 'POST https://api.github.com/repos/xunleii/github-quick-actions/issues/3/comments' with '200 {"id": 1235}'

    When Github sends an event "issue" with
      """
      {
        "action": "created",
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": {
          "body": "/duplicate #1 #2",
          "number": 3,
          "user": { "login":"xunleii" }
        },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/duplicate" for "issue" event with arguments ["#1","#2"] by sending these following requests
      | API request method | API request URL                                                                | API request payload        |
      | GET                | https://api.github.com/repos/xunleii/github-quick-actions/issues/1             |                            |
      | GET                | https://api.github.com/repos/xunleii/github-quick-actions/issues/2             |                            |
      | POST               | https://api.github.com/repos/xunleii/github-quick-actions/issues/3/comments    | {"body":"Duplicate of #1"} |
      | DELETE             | https://api.github.com/repos/xunleii/github-quick-actions/issues/comments/1234 |                            |
      | POST               | https://api.github.com/repos/xunleii/github-quick-actions/issues/3/comments    | {"body":"Duplicate of #2"} |
      | DELETE             | https://api.github.com/repos/xunleii/github-quick-actions/issues/comments/1235 |                            |

  @duplicate
  Scenario: /duplicate #1 #1
    # Comment ID required to remove the comment
    Given Github replies to 'POST https://api.github.com/repos/xunleii/github-quick-actions/issues/2/comments' with '200 {"id": 1234}'

    When Github sends an event "issue" with
      """
      {
        "action": "created",
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": {
          "body": "/duplicate #1 #1",
          "number": 2,
          "user": { "login":"xunleii" }
        },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/duplicate" for "issue" event with arguments ["#1","#1"] by sending these following requests
      | API request method | API request URL                                                                | API request payload        |
      | GET                | https://api.github.com/repos/xunleii/github-quick-actions/issues/1             |                            |
      | POST               | https://api.github.com/repos/xunleii/github-quick-actions/issues/2/comments    | {"body":"Duplicate of #1"} |
      | DELETE             | https://api.github.com/repos/xunleii/github-quick-actions/issues/comments/1234 |                            |

  @duplicate
  Scenario: invalid /duplicate 1
    When Github sends an event "issue" with
      """
      {
        "action": "created",
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": {
          "body": "/duplicate 1",
          "number": 2,
          "user": { "login":"xunleii" }
        },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/duplicate" for "issue" event with arguments ["1"] without sending anything

  @duplicate
  Scenario: invalid /duplicate wrong #issue
    When Github sends an event "issue" with
      """
      {
        "action": "created",
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": {
          "body": "/duplicate wrong #issue",
          "number": 2,
          "user": { "login":"xunleii" }
        },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/duplicate" for "issue" event with arguments ["wrong","#issue"] without sending anything

  @duplicate @error
  Scenario: /duplicate without arguments
    When Github sends an event "issue" with
      """
      {
        "action": "created",
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": {
          "body": "/duplicate",
          "number": 2,
          "user": { "login":"xunleii" }
        },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/duplicate" for "issue" event without argument without sending anything

  @duplicate @error
  Scenario: /duplicate itself
    Given Github replies to 'GET https://api.github.com/repos/xunleii/github-quick-actions/issues/1' with '404 {"message": "Not Found", "documentation_url": ""}'

    When Github sends an event "issue" with
      """
      {
        "action": "created",
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": {
          "body": "/duplicate #1",
          "number": 1,
          "user": { "login":"xunleii" }
        },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/duplicate" for "issue" event with arguments ["#1"] without sending anything

  @duplicate @error
  Scenario: /duplicate #1 when issue #1 doesn't exist
    Given Github replies to 'GET https://api.github.com/repos/xunleii/github-quick-actions/issues/1' with '404 {"message": "Not Found", "documentation_url": ""}'

    When Github sends an event "issue" with
      """
      {
        "action": "created",
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": {
          "body": "/duplicate #1",
          "number": 2,
          "user": { "login":"xunleii" }
        },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/duplicate" for "issue" event with arguments ["#1"] by sending these following requests
      | API request method | API request URL                                                                | API request payload        |
      | GET                | https://api.github.com/repos/xunleii/github-quick-actions/issues/1             |                            |

  @duplicate @error
  Scenario: error handling on /duplicate
    Given Github replies to 'POST https://api.github.com/repos/xunleii/github-quick-actions/issues/2/comments' with '404 {"message": "Not Found", "documentation_url": "https://docs.github.com/en/rest/reference/issues#add-labels-to-an-issue"}'
    When Github sends an event "issue" with
      """
      {
        "action": "created",
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": {
          "body": "/duplicate #1",
          "number": 2,
          "user": { "login":"xunleii" }
        },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/duplicate" for "issue" event with arguments ["#1"] but returns this error: 'POST https://api.github.com/repos/xunleii/github-quick-actions/issues/2/comments: 404 Not Found []'
