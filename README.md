# CLInote

CLInote is a command line client for Evernote inspired by [geeknote](https://github.com/VitaliyRodnenko/geeknote).

CLInote allows you to:

* Create notes
* Send the note content to stdout
* Edit notes in your $EDITOR
* Create new notebooks
* Search for notes

## Installation

From source:
```
go get -v github.com/TcM1911/clinote
```

## Authorize to Evernote via OAuth

Before you can use any features, you need to authorize CLInote to access youre notes. To authorize run the command:
```
clinote user login
```
If you have your default browser defined in the $BROWSER environment variable, CLInote will open the link in your default browser.

## Create a new note

A new note can be created with the command shown below. A title needs to be given for the note. If no notebook is given, the default notebook will be used. The new note can be open in the $EDITOR by using the edit flag.

```
clinote note new --title "note title" [--notebook "notebook name"] [--edit]
```

## Edit note

Notes can be edited using the edit command. If no flags are set, the note is opened
with the editor defined by the environment variable $EDITOR. The first line will be used as the note title and the rest is encoded as the note content.

To change to title, the title flag can be used.

The note can be moved to another notebook by defining the new notebook
with the notebook flag.
```
clinote note edit "note title" [--title "new note title"] [--notebook "new notebook"]
```

## Show note content

You can send the note content to the standard out with the command below:
```
clinote note "note title"
```

## Remove a note

Delete moves the note into the trash. The note may still be undeleted, unless it is expunged.
To expunge the note you need to use the official client or the web client.
```
clinote note delete "note title"
```

## Search for notes

To search for notes, use the list command as shown below.
```
clinote note list [--count 20] [--search "search term"] [--notebook "notebook name"]
```
The search term flag can be used to define a search term
to be used. The search can be restricted to a notebook
by using the notebook flag.

Count can be used to restrict the maximum number of notes
returned.

If no search term is given, a wild card search will be used.
The notes will be sorted by the modified time.

## Create a new notebook

To create a new notebook, use the command below:
```
clinote notebook new "notebook name" [--default] [--stack "Stack name"]
```

## Edit a notebook

To edit a notebook use this command:
```
clinote notebook edit "notebook name" [--name "new notebook name"] [--stack "new stack"]
```

## List all notebooks

To list all notebooks, use the notebook list command:
```
clinote notebook list
```
