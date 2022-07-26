# Contributing to ICON Bridge

The following is a set of guidelines for contributing to ICON Bridge.

These guidelines are subject to change. Feel free to propose changes to this document in a pull request.

## Pull Request Checklist

Before sending your pull requests, make sure you do the following:

-   Read the [contributing guidelines](CONTRIBUTING.md).
-   Check if your changes adhere to the [guidelines](https://github.com/icon-project/community/blob/main/guidelines/technical-development/development-guidelines.md).
-   Run the [unit tests](#running-unit-tests).

## Code of Conduct

The ICON Bridge project is governed by the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/2/1/code_of_conduct/code_of_conduct.md). Participants are expected to uphold this code.

## Questions

> **Note:** Github Issues are reserved for feature requests and bug reporting. Please don't create a new Github issue to ask a question.

We have a vibrant developer community and several community resources to ask questions in.

### Community Resources

* [Github Discussions](https://github.com/icon-project/icon-bridge/discussions)
* [ICON Official Discord](https://discord.gg/qa9m4bgKHE)

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check  **[our list of issues](https://github.com/icon-project/icon-bridge/issues)** to see if an
issue already exists.
> **Note:** For existing issues, please add a comment to the existing issue instead of opening a new issue. If the issue is closed and
> the problem persists, please open a new issue and include a link to the original issue in the body of your new one.

When you are creating a bug report, please fill out [the required template](https://github.com/icon-project/icon-bridge/blob/main/.github/ISSUE_TEMPLATE/bug.md) and include as many details as possible.

### Contributing Code

If you want to contribute, start working through the icon-bridge repository, navigate to the Github "issues" tab and start looking through issues. We recommed issues labeled "good first issue". These are issues that we believe are particularly well suited for newcomers. If you decide to start on an issue, leave a comment so that other people know that you're working on it. If you want to help out, but not alone, use the issue comment thread to coordinate.

Please see the [ICON Foundation Development Guidelines](https://github.com/icon-project/community/blob/main/guidelines/technical-development/development-guidelines.md)
for information regarding our development standards and practices.

### ICON Bridge Core Development Process

This section is intended for developers and project managers that are involved with core ICON Bridge and BTP integrations.

We use Zenhub for project management.

Issues that are not currently being worked on but determined to be part of this month's release should be in the "Release Backlog" column. Any other issue that is not being worked on that is not part of this month's release goes in the "Icebox" column.

Issues currently being worked on should be in the "In Progress" column, assigned this week's sprint, and assigned a time estimate (# of days to complete the issue). It's ok if the time estimate is not accurate, as it is only an estimate.

Issues currently being worked on should each have an associated branch. If an issue needs multiple branches, the issue is probably too large and should be broken down into multiple smaller issues. If the issue is describing a feature, the naming convention of the branch should be "feature/[description]-[issue #]". If the issue is describing a bug fix, the naming convention of the branch should be "fix/[description]-[issue #]".

After the developer working on the branch determines the feature (or fix) branch satisfies the acceptance criteria of the associated issue and has sufficiently tested the added code, the developer should submit a pull request with at least 1 reviewer assigned for the feature branch to be merged into the main branch. When the branch is submitted for pull request, the associated issue should go into the "Review / QA" column in Zenhub.

After the pull request has been merged, the feature branch should be deleted and the issue should be closed. When the Git issue is closed, the issue should go into the "Completed" column in Zenhub.

> **Note:** The process of moving an issue from one column to another is currently entirely manual. There are clear trigger points of when an issue should move from one column to another, so we intend to automate this in the future.
