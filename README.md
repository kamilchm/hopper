# Usage

## Installing hop

To install single hop or list of hops, issue the following command:

    $ hop install hop_name1 hop_name2

You could also use wildcards in names likes:

    $ hop install hop_*

It will skip all the existing hops.

## Fixing hops

    $ hop install --fix node

## Removing hops

    $ hop remove hop_*

## Finding out hops

    $ hop info httpie

## Hop definitions

You define hops in a simple YAML files and placing it in the
``~/.hopper/hops`` directory. See example hops file in ...

## Hop precedence

## Local vs User mode

Sometimes there's a cases where you want to use projest specific
commands that are local to selected workspace. Let's say you are
working on a frontend project than needs the Node.js to build.
You could define ``node`` hop for that project only. All you need
is to provide ``hop.yaml`` file in your project worksapce directory
and define ``node`` hop there. If you call hopper there it will
be started in local mode, which means that it will install all hops
in the local dir instead of user profile and it'll will you hops
definiotions from local ``hops.yaml`` file.

# Goals

- [x] run local hops from hop.yaml
- [x] run hops with access to host cwd
- [x] run hops using docker API
- [x] run hops on stdin using unix pipes
- [x] use hop stdout with unix pipes
- [x] unit tests
- [ ] install hops
  - [ ] local hops in cwd with hop.yaml
  - [ ] user hops defined in ~/.hopper/hops/*.yaml
- [ ] sexy README
- [ ] uninstall hops
  - [ ] all hops in cwd
  - [ ] all user hops
- [ ] update hops
  - [ ] remove dead hops
  - [ ] add new hops
- [ ] run hops with host $HOME
- [ ] ``install --local/user/global``
- [ ] using hopper ENV vars in hops
- [ ] support for rkt

## Similiar tools

* https://github.com/tailhook/vagga
* https://github.com/jamiemccrindle/dockerception
