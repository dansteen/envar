# envar
Small application that replaces parameter references to environment variables with their values.

# The Problem

Currently, container engines don't process environment variables passed on the command line of their "exec" functions.  Docker has this to say:
> Environment variables are supported by the following list of instructions in the Dockerfile:
>
> ADD
> COPY
> ENV
> EXPOSE
> LABEL
> USER
> WORKDIR
> VOLUME
> STOPSIGNAL

Note that both `ENTRYPOINT` and `CMD` are missing from the above list.   

The Appc Spec (as used by rkt et al) says the following:
> *exec* (list of strings, optional) executable to launch and any flags. ...  These strings are not evaluated in any way and environment variables are not substituted.

However, passing variables into the container at runtime is an extermely valuable way of modifying behavior, at runtime, based on environment.   In a Prod environment, I may want one set of values, while in a Staging environment I may want a different set.   The restrictions above seem especially odd when you take into account the fact that you *can* pass variables in, and they do impact the run of the proccess - they are just not interpreted on the actual command line you submit.

# The Solution

envar is a very small application that reads in a command line, and replaces any environment variables that it finds with their values in the environment.  It then calls `exec` and replaces itself (the PID stays the same) with the command passed in on the command line.  Any variables that are not found in the environment are removed prior to execution. PATH searching is handled as defined [here](https://golang.org/pkg/os/exec/#LookPath).

## Usage
`envar <command> <command_options>`

## Example

```
#> export DURATION=3
#> envar sleep \$DURATION
```
This will execute `sleep 3`. Note that we escape the `$` above because we are executing on the command line, and otherwise the command line will interpret the variable for us.   Inside a container definition you will likely not need that (depending on what you are using to build the definition of course).

# Known Issues

There is a large push in the Docker world to have images be as small as possible.  Unforunately, because envar is written in go, the binary size at build is approximately 3Mb.  This doesn't sound like much, but considering that it's 50 lines of code, and the equivelant code in C would be about 50kb, it seems like a bit much.

You can mitigate this to some extent by building with:
```
go build -ldflags "-s -w"
```
and then running the resulting binary through `upx`:
```
upx envar
``
This will result in a binary of about 600kb, which is still large but better than 3Mb.  Of course, this comes with it's own tradeoffs, but it's the usuall speed vs size calculation. 
