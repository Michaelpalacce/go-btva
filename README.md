# go-btva

Sets UP BTVA for local development. 💻

## Overview 

[Build Tools For VMware Aria](https://github.com/vmware/build-tools-for-vmware-aria) is a tool that enables easier development for automation on VCF.
BTVA comes with pre-requisites and even third party systems in order for everything to work correctly.

This tool aims to ease that initial setup.

## Installing 🚀

Go to the [releases](https://github.com/Michaelpalacce/go-btva/releases) and download the correct binary depending on your OS.
- The binaries are generated entirely in a [Github
  Action](https://github.com/Michaelpalacce/go-btva/blob/main/.github/workflows/build.yaml). The sha is shown in the [job run](https://github.com/Michaelpalacce/go-btva/actions/runs/14332857255/job/40172667918#step:4:262) and is also uploaded in the releases page.

## What This Is Not

- TBD, it's everything

## What This Is

- [x] 📝 Provides clear instructions on what will be done and prompts the user for actions that may affect his system(s)
- Setup minimal infra on a vm
    - [x] 🏗️ Run the minimal infrastructure installer
    - [x] 🔑 Fetch Nexus Password
    - [x] 🔑 Fetch Gitlab Password
    - [x] 🔑 Create New Gitlab Public Access Token
    - [x] 🏁 Register Gitlab Runner
- Setup dev environment
    - Supported os
        - [x] Linux (Ubuntu)
        - [x] MacOS
        - [ ] ⚒️ Windows
    - Configure
        - [x] 🏗️ Configure `settings.xml` for nexus and aria
        - [x] 🗒️ Provide Custom non minimal infra settings to integrate with other Artifact Managers
- Software
    - Ability to install software agnostic to the environment
        - [x] ⚡ Install Java
        - [x] ⚡ Install mvn
        - [x] ⚡ Install NodeJS with `fnm`
        - [x] ⚡ Install VSCode with recommended extensions
- At the end, give details to the user what their next steps should be
    - [x] 📝 Give Gitlab Credentials
    - [x] 📝 Give Nexus Credentials

## Action Items

### Must Haves

- [ ] Ask user if state should be used if ran interactively
- [ ] Non-interactive run should not ask for anything and instead fail if a question is needed.

### Good To Haves

- [ ] Modify the zprofile for java
- [ ] Working on windows
- [ ] Configure Other Artifact Managers?

### Nice to haves

- [ ] Unit test the solution

## State

State is managed by a state file that is created where the tool is ran. After the initial run of the tool, CLI arguments are ignored and
instead the ones stored in the state file are used. If you want to do any changes, do so in the state file. As the whole process is
idempotent, you can also remove the state file and re-run with the desired arguments.

> State contains sensitive information for now. Be carefull when opening it.


<details>
    <summary>Finished state</summary>
    <img src="assets/state-finished.png"/>
</details>

## Development 📦

We use `make` to run the program for dev

## Running ⚡

```sh
make run
```

## Cleanup 🧹

Cleanup scripts are provided for linux to ease testing. There is a generic `cleanup` goal that can be used to cleanup everything, or more
specific `cleanup-mvn` for example to cleanup specific components only.

```sh
make cleanup
```

## Tests ✅

```sh
make test
```

## Makefile

The `Makefile` contains a bunch of different helper methods. You can run `make help` to get a description of what is available.

## FAQ ❓

### Why is there no Air Gapped Installation? 🌌

Doing a full air gapped installation may be viewed as a security issue.
We need to download all the dependencies for software and store them as install binaries. Potentially these binaries can be embedded in the executable, 
however this introduces an attack vector as enterprises need to trust that those binaries are safe.

Alternative would be to give clients the ability to specify their own binaries, but then at that point there is little point for this tool
to handle the software installation.
