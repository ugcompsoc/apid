## What's this for?

Currently this is the framework for effectively CompSoc's APIv2 which will include a proper CI/CD system.

## Contributing

### Repo Conventions

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
    ./.githooks/git-hooks-config.sh
```

##### Can I just not?
You can 'just not' if you want. The plan in the future is to implement a workflow that verifies commits follow this convention and, if not followed, will block commits.

If there is a legitimate reason for not following this convention, you can use:
```
    git commit --no-verify
        or
    git push --no-verify.
```
