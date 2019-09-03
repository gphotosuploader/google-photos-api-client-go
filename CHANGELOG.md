# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/) and this project adheres to [Semantic Versioning](https://semver.org/).

## Unreleased
### Added
- Code quality reporting using [Codebeat](https://codebeat.co/projects/github-com-gphotosuploader-google-photos-api-client-go-master).
- Code coverage reporting using [codecov](https://codecov.io/gh/gphotosuploader/google-photos-api-client-go).

### Changed
- CI platform is now [drone.io](https://drone.io).

## 1.0.4
### Fixed
- Fix `AlbumByName` to check against all Google Photos album list (#12).

## 1.0.3
### Added
- Add resumable file uploads. You can use new `UploadFileResumable` method to upload files that can be resumed. See [documentation](https://godoc.org/github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos) for more details.

## 1.0.2
### Added
- Add makefile for easy test and linting
- Add CI using travis-ci.com
- Add travis and goodocs badges to README

## 1.0.1
### Added
- Add semantic versioning to this package
- Add Go module support
- Add MIT license

### Changed
- Clean up of useless files / directories
- Update package documentation

### Fixed
- Fix issue #8 on parent repository [here](https://github.com/nmrshll/google-photos-api-client-go/issues/8)

### Removed
- Removed `Makefile` support

## 1.0.0
Initial release after clone it from [original repository](https://github.com/nmrshll/google-photos-api-client-go). Latest commit was [3dac07f](https://github.com/nmrshll/google-photos-api-client-go/commit/3dac07f1b07f249ac4a9805b9d60afe0f68c34b2)
