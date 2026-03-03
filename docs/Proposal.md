# \[Project Title\]
Members: Josh Justice, Leo Farmerie, Harshvardhan Singh

## Goal Summary:
Create our own implemantion of nav. Allows users to see into the directory as they navigate the file system.
Recreate grep (and possibly other command line utilites) with improvements + additional features that we wish were a part of the command
*Potential Other Utilities* 
- find
- some kind of process manager like top?
- xargs
- something of our own?
  - 

## Use Case
This is pretty straight forward, these two utilites are going to be used any time you would use what they aim to replace.

## Sketch
**TODO**

## Thoughts on testing
Unit tests can be created to test functionality as well as create situations where users act unfavorably in the interactive environment, making sure that any possible user action at a given state will be handled aprropriately.

## MVP
### polo (After Marco Polo)
- User can peek into subdirs and navigate to their desired destination
- Menu is interactive, navigation occurs with only one execution

### grep (no fun name yet)
- Search with and without regex


## Stretch Goals
### polo
- A stretch goal for polo is implementing a file preview pane that displays a preview of the currently selected file or directory. When the user navigates through entries, the preview pane would dynamically show relevant information such as the first few lines of a text file, metadata for other file types, or the contents of a directory, allowing users to inspect files without opening them.

### grep
- Interactive mode, updates result as regex is updated
  - Corrects users regex errors
- 

## Checkpoint Goals
Because our MVP is so minimal, we feel that the meat of the project will be in the strech goals + enhancements to the MVP. Our goal for the checkpoint is to have the MVP done, as both of these utilities are pretty simple at their base. It's possible that we run into issues with the interactive menu of polo, so its possible that that may not be working 100% by the checkpoint. Otherwise, we feel pretty safe to expect that grep will be minimally complete by that time.
