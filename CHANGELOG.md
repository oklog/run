# Changelog

## [0.0.5] - 2024-05-08

### Added

- Added support for ordered shutdown.

## [0.0.4] - 2024-03-31

### Added

- Added `contrib/ticker`.

## [0.0.3] - 2024-02-21

### Changed

- Renamed the `contrib/waitgroup` constructor.

## [0.0.2] - 2024-02-21

### Added

- Added `contrib/waitgroup`.
- Partial unit test coverage for `run.go`.
- Some additional code comments on `runnable.Run`.

### Fixed

- Fixed syntax error in `group_test.go`.
- Wait until `n` messages come across `closeCh`, not `n-1`. This was a copy/paste issue.

## [0.0.1] - 2024-02-20

### Added

- Initial release.
