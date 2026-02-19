*\[Project Title\]
Members: Josh Justice & Leo Farmerie

** Goal Summary:
Create our own implemantion of nav. Allows users to see into the directory as they navigate the file system.
Recreate grep (and possibly other command line utilites) with improvements + additional features that we wish were a part of the command
*Potential Other Utilities* 
- find
- some kind of process manager like top?
- xargs
- something of our own?
  - 

** Use Case
This is pretty straight forward, these two utilites are going to be used any time you would use what they aim to replace.

** Sketch
**TODO**

** Thoughts on testing
Unit tests can be created to test functionality as well as create situations where users act unfavorably in the interactive environment, making sure that any possible user action at a given state will be handled aprropriately.

** MVP
*** polo (After Marco Polo)
- User can peek into subdirs and navigate to their desired destination
- Menu is interactive, navigation occurs with only one execution

*** grep (no fun name yet)
- Search with and without regex


** Stretch Goals
*** polo
- 

*** grep
- Interactive mode, updates result as regex is updated
  - Corrects users regex errors
- 