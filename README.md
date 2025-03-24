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

- [x] Implement Software installation Linux
- [ ] Implement Environment Setup Linux
- [ ] Implement Infra Setup Linux

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
