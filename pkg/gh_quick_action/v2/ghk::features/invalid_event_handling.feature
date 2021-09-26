Feature: handle Github event
  In order to manage quick actions on Github
  As Github application
  I need to handle correctly Github events

  Background:
    Given quick action "/hello_world" is registered for "issue_comment" events

  Scenario: should handle invalid event type with an appropriate error
    When Github sends an event "unknown_event" with
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
    Then Github Quick Actions should return these errors
    """
    event type 'unknown_event' not managed
    """

  Scenario: should handle invalid event JSON with an appropriate error
    When Github sends an event "issue_comment" with
      """
      ...
      """
    Then Github Quick Actions should return these errors
    """
    failed to extract data from JSON for event 'issue_comment': invalid character '.' looking for beginning of value
    """

  Scenario: should do nothing if no valid command found
    When Github sends an event "issue_comment" with
      """
      {
        "action": "created",
        "comment": { "body": "/ignore_me" },
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 1 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions shouldn't do anything

  Scenario: should do nothing if the event action is not "created"
    When Github sends an event "issue_comment" with
      """
      {
        "action": "modified",
        "comment": { "body": "/hello_world" },
        "repository": {
          "owner": { "login": "xunleii" },
          "name": "github-quick-actions"
        },
        "issue": { "number": 1 },
        "installation": { "id": 123456789 }
      }
      """
    Then Github Quick Actions shouldn't do anything
