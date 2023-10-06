# Hangman

Play hangman online through SSH!


## Setup

### Running the server
To run the server, make sure to `go run .` from the root of the project. Then open another terminal window and `ssh -p 1337 localhost` to connect to it.

### SSH Configs
Add the following to your `~/.ssh/config` to avoid having to clear out localhost entries in your `~/.ssh/known_hosts` file:
```
Host localhost
    UserKnownHostsFile /dev/null
```

## TODO

- [X] Create a simple ssh server that you can connect to that returns "Hello, from server!".
- [ ] Create the game handler middleware!
- [ ] Game logic : create the game using bubbletea for single player first!
- [ ] Game logic : use API for fetching words ?
- [ ] Online multiplayer part : use redis or some database for holding score and other stuff ?
- [ ] Online multiplayer part : handle many connections ?
- [ ] Online multiplayer part : room creation for others to join.
