# Test Coverage with Codecov

## Introduction
To know how much of source code has been tested to give confidence on code, code coverage check is implemented with Github Actions. The report is generated within Github Actions and then sent to codecov.io for visualization.

## Coverage Check Conditions
[coverage.yml](https://github.com/icon-project/icon-bridge/blob/main/.github/workflows/coverage.yml) file inside the Github Actions workflow checks coverage if these conditions matches:
1. PR to main branch is created/update
2. Code is changed/merged on main branch

## Pull Request Check
PR to main branch is checked everytime code is updated and a comment is dropped based on the percentage coverage by `codecov bot` user.

## Coverage Configuration
[codecov.yml](https://github.com/icon-project/icon-bridge/blob/main/codecov.yml) file at the root of the project is the configuration file for PR check, comment to PR etc.

## Coverage Reporting
The coverage reporting dashboard is publicly available on codecov URL: https://app.codecov.io/gh/icon-project/icon-bridge

Also, the main branch coverage percentage is available on the top of the project's root README file in graph.
