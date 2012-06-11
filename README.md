**Victor** is a 37Signals Campfire bot written in Go (http://golang.org); inspired by Github's Hubot (http://github.com/github/hubot).

### Playing with Victor

#### 1. Get the code

Use the `go get` command to fetch the latest build of the project:

    $ go get github.com/brettbuddin/victor

#### 2. Install the Sample Program

To quickly start 

    $ go install github.com/brettbuddin/victor/examples/bot

#### 3. Connect to Your Campfire Room

    # Set some important environment variables
    $ export VICTOR_CAMPFIRE_ACCOUNT=
    $ export VICTOR_CAMPFIRE_TOKEN=
    $ export VICTOR_CAMPFIRE_ROOMS=
    
    # Make him join the party
    $ bot -adapter campfire

### Making Him Your Own

Victor doesn't do a whole lot out-of-the-box right now. I'll be adding more default behavior to him as time progresses, but you might want him to do something specific to your team's needs. 

You can use the programs located in `examples/` as starting points to create your own executable.

### Robot Abilities

- `Hear`: Trigger an action based on some criteria heard anywhere in the channel.
- `Respond`: Respond to a direct statement at the robot (e.g. "gobias show me the diff")

### Actions on Messages

- `Send`: Send a bit of text to the channel.
- `Reply`: Reply directly to the person that triggered the action (e.g. "Brett: Yo yo yo. Here's the diff:")
- `Paste`: Paste text to the channel.
