# CodeQL Check
## Introduction
For finding security vulnerabilities in the code, we are using GitHub's semantic code analysis engine, [CodeQL](https://github.com/github/codeql-action) inside Github Actions workflow. It automatically uploads the results to GitHub so they can be displayed in the repository's security tab.

## Github Actions
The workflow file [codeql-analysis](https://github.com/icon-project/icon-bridge/blob/main/.github/workflows/codeql-analysis.yml) runs in two cases:
1. PR or merge to main branch
2. Schedule every week on Thursday

## CodeQL Report
The CodeQL actions updates security details of the repo on [Code Scanning section](https://github.com/icon-project/icon-bridge/security/code-scanning).
