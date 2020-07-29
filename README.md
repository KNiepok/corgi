# Corgi

Slackbot for reminding team members to update their timesheets.

## Setup

### Prepare Slack application
You need to configure a slack application.
Create a new slack app, then head to Features / Slash commands section and add commands:
`/subscribe` ,`/unsubscribe`.
These need to be paired with /subscribe and /unsubscribe API endpoints of this project.

Then, you need to acquire two tokens:
1. Go to `Install App` settings - this is your token for `SLACK_TOKEN` env.
2. Go to `Basic information` settings and find Verification Token. This is your `SLACK_VERIFICATION_TOKEN` variable.

### Run the app

Use `make docker` to build a docker image with the app. Then you can run the image and bind ports.

## Local Development

### Communication with Slack app
To develop app locally, I recommend creating your own Slack namespace to communicate with slash commands
and to use Ngrok to expose locally run app. Run `ngrok http <port>` to port-forward the app, then use the generated URL
to bind slash commands on api.slack,.com
 
### Testing
Run `make lint` and `make test` to check code.

## Usage
After app is connected to slack, call Corgi with
`/subscribe daily @ 17:15` to get daily notifications.
You can also susbscribe for weekly notifications:
`/subscribe weekly @ SAT 9`



## Todo, ideas
* Swap slash commands for mentions 
* Commands for admin-tier team members to notify members immediately
* PostgreSQL instead of sqlite (or both)
* Endpoint for lookup of what interval I entered
* Support for monthly intervals (maybe such a big interval misses the point of reminding)
* Add ability for admin-tier team members to create notifications for members (with side effect of notifying them)
* Low priority: Keep history of changes of subscriptions
* Move representation out of tempo package (also rename message generator then)
* Cache existing (jira accountID, email) pairs