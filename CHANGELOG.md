# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2020-05-23
### Added
- Error handling in .Process() interface method.

### Changed
- Refactoring to prepare current structure to support proper packages;
- Fibonnaci algorithm to a faster one (variable substitution rather than using recursion);
- Time function now produces another function which diffs with the time of its creation;
- Module name.

### Fixed
- Some unecessary methods and stdout printing;
- Interfaces and structures are now receiving any values as its input.

### Fixed
- README Golang version on installation requirement step.

## [0.1.0] - 2020-05-21
### Added
- Fibonnaci calculation that leverages computation to as many cores as the user wants (respecting CPU limits);
- Changelog file;
- Gitignore file;
- Licensing information.
