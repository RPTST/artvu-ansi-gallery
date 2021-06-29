# Artvu Ansi Gallery

ANSI art viewer for BBSs (like Mystic, Talisman, WWiV, ENiGMA 1/2, etc.). Supports 80 and 132 terminal widths.

Required command-line flag, path to art folder:

./artvu-ansi-gallery -root /path/to/art

## Notes ##

Artvu scans for files and folders and turns that into directory-style lists that can be navigated with a lightbar. It incorporates SAUCE metadata into file listings. It can run on a tcp connection using a BBS's door function (stdout).

Linux, Mac, RPi compatabile.

## Build ##

Build like any go program. Or, use a pre-built binary from the /release folser.

## TO-DO ##
- [ ] Create a single object holding directory and file data - Refactor directory/file scanning to create a single go struct of the entire directory tree that can hold metadata, like SAUCE data (right now it scans only the "active" directory and creates 2 simple slices holding the directory and file names). This will allow for more granular sorting and filtering.
- [ ] Add "Preview" data panel for terms wider than 80 cols
- [ ] Reduce flicker - Refactor lightbar function to *only* re-draw the entire list if scolling beyond the terminal height. Otherwise, handle only affected rows.
- [ ] Allow users to manually scroll up/down in ANSI
- [ ] Add --local command line flag to allow users to browse locally - needs to handle/print ANSI differently
  
