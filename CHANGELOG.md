# Changelog

## 0.4.1

### Bug fix release

* Added overflow check for cached note list
* Fix bug when moving note to different notebook

## 0.4.0

### BoltDB storage

This release adds support for BoltDB storage.

* Notebook list is cached. The cache is refreshed if it's older than 24 hours or when you force a refresh using the new flag `-s`.
* Searches of notes are cached. See README for use cases.

## 0.3.0

### Improved Markdown support

Switched to a different HTML to Markdown library. Should parse the note content better.

### Code refactoring

Moved all Evernote code into one package.

## 0.2.0

### Raw content support

A notes content can now be viewed and edited in XML form instead
of markdown. Call the command edit, new, or note with the `--raw`
flag to edit, create, or view the note content in XML format.

### Start the browser in the background

The external browser during login is now started in the background.

### Bugfix empty notes returned

The XML decoder was to strict which could cause it to fail and an
empty note content was returned. The new decoder is less strict.

## 0.1.0

Initial release

