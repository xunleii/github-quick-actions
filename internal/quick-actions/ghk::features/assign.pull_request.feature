@pull_request
Feature: assign someone with /assign @user [@user...] on pull request description

  Background:
    Given quick action "/assign" is registered for "pull_request" events

  @assign
  Scenario: /assign @mojombo
    When Github sends an event "pull_request" with
      """
      {
        "action": "created",
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "pull_request": {
          "body": "/assign @mojombo",
          "number": 1,
          "user": { "login":"xunleii" }
        },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/assign" for "pull_request" event with arguments ["@mojombo"] by sending these following requests
      | API request method | API request URL                                                              | API request payload       |
      | POST               | https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees | {"assignees":["mojombo"]} |

  @assign
  Scenario: /assign @mojombo @defunkt
    When Github sends an event "pull_request" with
      """
      {
        "action": "created",
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "pull_request": {
          "body": "/assign @mojombo @defunkt",
          "number": 1,
          "user": { "login":"xunleii" }
        },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/assign" for "pull_request" event with arguments ["@mojombo","@defunkt"] by sending these following requests
      | API request method | API request URL                                                              | API request payload                 |
      | POST               | https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees | {"assignees":["mojombo","defunkt"]} |

  @assign
  Scenario: /assign me
    When Github sends an event "pull_request" with
      """
      {
        "action": "created",
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "pull_request": {
          "body": "/assign me",
          "number": 1,
          "user": { "login":"xunleii" }
        },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/assign" for "pull_request" event with arguments ["me"] by sending these following requests
      | API request method | API request URL                                                              | API request payload       |
      | POST               | https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees | {"assignees":["xunleii"]} |

  @assign
  Scenario: /assign @mojombo me
    When Github sends an event "pull_request" with
      """
      {
        "action": "created",
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "pull_request": {
          "body": "/assign @mojombo me",
          "number": 1,
          "user": { "login":"xunleii" }
        },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/assign" for "pull_request" event with arguments ["@mojombo","me"] by sending these following requests
      | API request method | API request URL                                                              | API request payload                 |
      | POST               | https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees | {"assignees":["mojombo","xunleii"]} |

  @assign
  Scenario: /assign @mojombo @mojombo
    When Github sends an event "pull_request" with
      """
      {
        "action": "created",
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "pull_request": {
          "body": "/assign @mojombo @mojombo",
          "number": 1,
          "user": { "login":"xunleii" }
        },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/assign" for "pull_request" event with arguments ["@mojombo","@mojombo"] by sending these following requests
      | API request method | API request URL                                                              | API request payload       |
      | POST               | https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees | {"assignees":["mojombo"]} |

  @assign @error
  Scenario: /assign mojombo
    When Github sends an event "pull_request" with
      """
      {
        "action": "created",
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "pull_request": {
          "body": "/assign mojombo",
          "number": 1,
          "user": { "login":"xunleii" }
        },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/assign" for "pull_request" event with arguments ["mojombo"] without sending anything

  @assign @error
  Scenario: /assign without argument
    When Github sends an event "pull_request" with
      """
      {
        "action": "created",
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "pull_request": {
          "body": "/assign",
          "number": 1,
          "user": { "login":"xunleii" }
        },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/assign" for "pull_request" event without argument without sending anything

  @assign @error
  Scenario: error handling on /assign
    Given Github replies to 'POST https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees' with '404 {"message": "Not Found", "documentation_url": "https://docs.github.com/en/rest/reference/issues#add-labels-to-an-issue"}'
    When Github sends an event "pull_request" with
      """
      {
        "action": "created",
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "pull_request": {
          "body": "/assign me",
          "number": 1,
          "user": { "login":"xunleii" }
        },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions should handle command "/assign" for "pull_request" event with arguments ["me"] but returns this error: 'POST https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees: 404 Not Found []'
