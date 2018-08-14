        +-------+         ____                                     _
        | || || |        / ___|___  _ __ ___  _ __   ___  ___  ___(_)
    +---+---+---+---+   | |   / _ \| '_ ` _ \| '_ \ / _ \/ __|/ _ \ |
    | || || | || || |   | |__| (_) | | | | | | |_) | (_) \__ \  __/ |
    +-------+-------+    \____\___/|_| |_| |_| .__/ \___/|___/\___|_|
                                             |_|



Composei is an interactive command line tool build with golang that helps you create your docker compose file.

Cause I'm too lazy to remember all the possible options for each container xD

# Installation
Run

    go get -u github.com/kariae/composei

# Usage

    $ composei -h
    NAME:
       Composei - Composei is an interactive command line tool build with golang that helps you create your docker compose file.

    USAGE:
       main [global options] command [command options] [arguments...]

    VERSION:
       0.1.0

    COMMANDS:
         generate, g  Generate docker compose file
         help, h      Shows a list of commands or help for one command

    GLOBAL OPTIONS:
       --help, -h     show help
       --version, -v  print the version

# TODOs
- [ ] Add more details to README file.
- [ ] Automate releases generation.
- [ ] Add travis automation for tests.

# Contributing
First, **many thanks** for your contributions, please note that this eco system is a personal preference that I use in most of my PHP projects (using Symfony or other frameworks), if you find any typo/misconfiguration, or just want to optimize more the workflow, please
1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D
