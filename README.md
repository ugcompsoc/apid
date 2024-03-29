## What's this for?

Currently, this is the framework for effectively CompSoc's APIv2 which will include a proper CI/CD system.

This is a REST API developed for the benefit of the societies and Computer Society members of the University of Galway. The end goal is that this REST API be available to the above mentioned to provision, and manage their CompSoc resources. This REST API will primarily used by [dash.compsoc.ie](https://dash.compsoc.ie) and [compsoc.ie](https://compsoc.ie). The first being the website where the above mentioned do all their administration through our portal, and the latter for querying when upcoming events we have.

## Why am I here?
I don't know. While you're here though, check out this [cool video](https://www.youtube.com/watch?v=dQw4w9WgXcQ).

![XKCD Santa Sudo Meme](https://imgs.xkcd.com/comics/incident.png "He sees you when you're sleeping, he knows when you're awake, he's copied on /var/spool/mail/root, so be good for goodness' sake.")

## Contributing

### Running the code

A handy docker compose file has been created that contains all of the dependencies needed to run this API. Run this in your terminal to start the API on port 8080 as well as any dependencies needed to run the service (A.K.A. MongoDB)

  docker compose -f docker-compose-local.yml up

### Workflows

There are many [workflows](.github/workflows) in this repository that handle testing the code, building the docker image, deploying the docker image to the respective environments, and handling deleting old deployments. These workflows are ran at various points of the software development cycle, the mains ones to mention are:
- When a commit is pushed to a branch, the branch is deployed to its own environment at the URL: `{branch_name}.dev.apid.testbox.compsoc.ie (DEV)`.
- When a PR is merged or commit is pushed to main, the same happens except it is deployed at the URL: `dev.apid.testbox.compsoc.ie (TEST)`.
- When a tag is created, the same happened except it is deployed at the URL: `apid.testbox.compsoc.ie (PROD)`.

On every branch or tag event, unit and API tests are also ran.

### Unit Testing

The unit tests in this repo follow the safe standard that all Go projects should. Unit tests are kept next to their relating source files and are named the same as the source file with a preceeding `_test`.

Every change to source code should be covered by a unit test (negative and positive tests)

### Code Coverage

A code coverage percentage of 50% is to be upheld or commits and build will be rejected. If your code coverage is below this, just go write some unit tests; make them good, don't just make them for the craic.

Currently the `cmd` package is ignored from code coverage (TODO [#21](https://github.com/ugcompsoc/apid/issues/21)).

The `internal/services/database_test_utils` package is permanetly ignored from code coverage as currently do not see a benefit to testing test utils.

### API Tests

Each new route should have a coresponding Venom API test. These tests should ideally only cover happy paths. These tests belong in the `venom` folder.

### Swagger

Swagger is responsible for collating all the route documentation and displaying it at `docs/index.html`. Every time there is a change to documentation of a route you should run `swag init --dir "cmd,internal/server,internal/helpers"`, this will update the relevant files so it will display the correct information. It is important that you run this command after you merge any changes into your branch that changed the route documentation.

### Repo Conventions

#### Workflows

All workflows should pass on your commit, PR, whatever it be. If they don't, revert it, don't merge it. Figure out why they aren't passing.

#### Conventional Commits
This repo makes use of the 'conventional commits' convention. This means that all commits should be formatted using this convention. 

For example:
    feat(router): added events route

An example of an incorrect commit message would be:
    i fixed an error that caused the events route to not work lol

All commit message should follow the follow standard:
```
    <type>: <subject>
        or
    <type(<scope>):> <subject>
```

A type can be: build, chore, ci, docs, feat, fix, perf, refactor, revert, style, test, or wip.
A scope should encompass the change in one to two words. It would usually refer to the location of the change, e.x. the controller, model, or router.
The subject should be limited to 50 characters and can describe your change. If you need to add more detail to the commit message, you can use a line break (\n).

There is a git hook that has been set up to help you create messaging and won't let you commit unless you do it right. All you need to do is execute the following script once after you clone the repo:
```
    .githooks/git-hooks-config.sh
```

##### Can I just not?
You can 'just not' if you want. ~~The plan in the future is to implement a workflow that verifies commits follow this convention and, if not followed, will block commits.~~

If there is a legitimate reason for not following this convention, you can use:
```
    git commit --no-verify
        or
    git push --no-verify.
```
