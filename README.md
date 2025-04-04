# go-btva

Sets UP BTVA for local development.

## Overview

[Build Tools For VMware Aria](https://github.com/vmware/build-tools-for-vmware-aria) is a tool that enabling easier development for automation on VCF.
BTVA comes with pre-requisites and even third party systems in order for everything to work correctly.

This tool aims to ease the initial setup.

## What This Is Not

- TBD, it's everything

## What This Is

- Setup minimal infra on a vm
    - [x] Run the minimal infrastructure installer
    - [x] Fetch Nexus Password
    - [x] Fetch Gitlab Password
    - [x] Create New Gitlab Public Access Token
    - [x] Register Gitlab Runner
- Setup dev environment
    - Supported os
        - [x] Linux (Debian)
        - [x] Darwin
        - [ ] Windows
    - Configure
        - [x] Configure `settings.xml` for nexus and aria
- Software
    - Ability to install software agnostic to the environment
        - [x] Install Java
        - [x] Install mvn
        - [x] Install NodeJS with `fnm`
        - [ ] Install VSCode with recommended extensions
- At the end, give details to the user what their next steps should be
    - [x] Give Gitlab Credentials
    - [x] Give Nexus Credentials

## Action Items

### Must Haves

- [ ] Build process that publishes runnables

### Good To Haves

- [ ] Create a demo project

### Nice To Haves

- [ ] Encrypt state file variables that are secret
- [ ] Window pop up
- [ ] Executable
- [ ] Working on windows

## State

State is managed by a state file that is created where the tool is ran. After the initial run of the tool, CLI arguments are ignored and
instead the ones stored in the state file are used. If you want to do any changes, do so in the state file. As the whole process is
idempotent, you can also remove the state file and re-run with the desired arguments.

> State contains sensitive information for now. Be carefull when opening it.


<details>
    <summary>Finished state</summary>
    <img src="assets/state-finished.png"/>
</details>

## Development

We use `make` to run the program for dev

## Running

```sh
make run
```

## Cleanup

Cleanup scripts are provided for linux to ease testing. There is a generic `cleanup` goal that can be used to cleanup everything, or more
specific `cleanup-mvn` for example to cleanup specific components only.

```sh
make cleanup
```

## Tests

```sh
make test
```

## Makefile

The `Makefile` contains a bunch of different helper methods. You can run `make help` to get a description of what is available.

