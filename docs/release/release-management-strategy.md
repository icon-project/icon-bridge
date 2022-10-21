# Release management strategy

We follow a release management strategy called [the release train](https://martinfowler.com/articles/branching-patterns.html#release-train).

### Release schedule

1. On the first Monday of every month, we plan what will be released that month.
2. On the third Monday of every month, we freeze the main branch and identify what will make it to release. If the code base is healthy,
what gets released will be exactly what is frozen on the main branch.
4. On the fourth Monday of every month, we tag the commit in main we are releasing with using the correct version bump.
5. Hotfixes are developed and deployed adhoc so that breaking issues do not have to wait to be resolved.
6. Hotfixes are developed in the main branch. If there are changes in the main branch between the last release and the hotfix,
then a release branch with the latest release gets cut and the hotfix gets cherrypicked to that release branch.
The release branch gets closed when the next release occurs.
7. Features that miss a release are simply moved to the following month.

### Release changelog

Every pull request squashes commits in accordance with the
[commit message guidelines](https://github.com/icon-project/community/blob/main/guidelines/technical/software-development-guidelines.md#commit-messages).
We simply look at the main branch commit history and cherrypick if there is anything in main that we are not confident about releasing as needed.

### Acceptance criteria

We can feel confident about release if code only makes its way to main AFTER:
* Following the [development guidelines](https://github.com/icon-project/community/blob/main/guidelines/technical/software-development-guidelines.md),
and particularly the testing and documentation sections.
* Passing CI checks
* Passing code review that the code flows logically and follows the
* [development guidelines](https://github.com/icon-project/community/blob/main/guidelines/technical/software-development-guidelines.md)
* Passing final sign off. We go through each commit in the history for that month at the start of a code freeze.
* During sign off, we note any external dependencies and plan out what needs to be done with them.
* For example, if we know a server needs to be configured, we identify who is involved with performing
* the configuration and schedule when that configuration will occur so that everyone involved is available at that time.

### Release versioning

See how [commits bump versions](https://github.com/icon-project/community/blob/main/guidelines/technical/software-development-guidelines.md#version-bumps). Because we squash commits on pull request, each commit is equivalent to approximately one pull request on main. We consider main's entire commit history for that month when determining the release version. Because the icon-bridge repo is a monorepo, each package has its own version.
