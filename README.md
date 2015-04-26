## Victor

[![Build Status](https://travis-ci.org/brettbuddin/victor.png?branch=master)](https://travis-ci.org/brettbuddin/victor)

**Victor** is a library for creating your own chat bot.

We use Victor as the backbone of our bot, Virbot, within our team Slack chat at Virb (http://virb.com). We use him for all sorts things like:

- Deploying code
- Preparing/initiating builds of our projects
- Viewing information about our infrastructure
- CDN operations
- Jokes and laughs

### Supported Services

I currently have adapters written for [Slack](https://slack.com/) and [Campfire](https://campfirenow.com/), however more are to come. Writing an adapter for your favorite service is a good way to contribute to the project :wink:.

*   **Slack**
    To use victor with Slack, you need to configure two slack integrations first: An **[incoming webhook integration](https://my.slack.com/services/new/incoming-webhook)** and an **[outgoing webhook integration](https://my.slack.com/services/new/outgoing-webhook)**. The *incoming* webhook integration gives you an URI which is used by victor to post messages into your slack channel. Pass this URI to victor via the environment variable `SLACK_INCOMING_WEBHOOK_URI`. In the *outgoing* webhook integration, you configure an URI that Slack uses to notify victor about Slack messages. The *path* of the configured URI needs to be set in the environment variable `SLACK_OUTGOING_WEBHOOK`.
    
        $ SLACK_INCOMING_WEBHOOK_URI=https://hooks.slack.com/services/REPLACE/THIS \
            SLACK_OUTGOING_WEBHOOK=/your-hook-path \
            ./your-chat-bot
    
    Another setting you can make is that will be displayed next to your bot's name in Slack. Pass the environment variable `SLACK_EMOJI` with your favorite emoji (in [colon syntax](http://www.emoji-cheat-sheet.com/)) to victor. 

### Making Him Your Own

Victor is more of a framework for constructing your own bot so he doesn't do a whole lot out-of-the-box. I'll be adding more default behavior to him as time progresses. You can use the programs located in `examples/` as starting points to create your own executable.
