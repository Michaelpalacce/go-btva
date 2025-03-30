# go-btva

Sets UP BTVA for local development.

## Responsibilities

- Setup minimal infra on a vm
- Setup dev environment
    - Supported os
        - Linux
        - Darwin (soon)
    - Configure
- Software
    - Ability to install software agnostic to the environment
- State
    - Allows us to resume and skip parts
- At the end, give details to the user what to do if needed

## Action Items

- [ ] Make minimal infra install idempotent

### Nice To Haves

- [ ] Window pop up
- [ ] Executable

# State

State is managed by a state file that is created where the tool is ran. After the initial run of the tool, CLI arguments are ignored and
instead the ones stored in the state file are used. If you want to do any changes, do the changes in the state file.

This tool **is not** a desired state machine. Instead it allows for a resumable process. You can however modify the state if you wish
certain steps to be repeated. For example, to repeat something, like installing node, you can remove that from the state and it will be
re-applied.

> State contains sensitive information for now. Be carefull when opening it.


<details>
    <summary>Finished state</summary>
    <img src="assets/state-finished.png"/>
</details>

# Development

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

