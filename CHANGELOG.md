# Changelog

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

