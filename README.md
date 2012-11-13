**Victor** is a library for creating your own chat bot for 37Signals' Campfire.

We use Victor as the backbone of our bot, Virbot, within our team Campfire at Virb (http://virb.com).

### Making Him Your Own

Victor is more of a framework for constructing your own Campfire bot so he doesn't do a whole lot out-of-the-box. I'll be adding more default behavior to him as time progresses, but you might want him to do something specific to your team's needs. You can use the programs located in `examples/` as starting points to create your own robot executable.

### Function Overview

#### Listening for Actions

- `Hear`: Trigger an action based on some criteria heard anywhere in the channel.
- `Respond`: Respond to a direct statement at the robot (e.g. "victor show not shipped")

#### Responding

Each listener that's triggered will provide a Context object. The Contex object contains information about the original message, matches from the regex that triggered the listener and a few helper methods for responding to the context from which the message originated. These methods include:

- `Send`: Send a bit of text to the channel.
- `Reply`: Reply directly to the person that triggered the action (e.g. "Brett: Yo yo yo. Howdy?")
- `Paste`: Paste text to the channel.
- `Sound`: Plays a Campfire sound
