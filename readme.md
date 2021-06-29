# Artvu Ansi Gallery
ANSI art viewer for old-skool, and new-school, BBSs (like Mystic, Talisman, WWiV, ENiGMA 1/2, etc.). Supports 80 and 132 terminal widths.

Required command-line flag is "-root" - which is the path to your folder containing Ansi art:

./artvu-ansi-gallery -root /path/to/art

Artvu scans for files and folders and turns them into directory-style listings that can be navigated via lightbar interface over telnet/ssh connection (stdout) using a terminal program like SyncTerm, MagiTerm or NetRunner. 

It incorporates SAUCE metadata into the file listings. 

Linux, Mac, RPi compatible. No Windows support -- yet.

## Build ##
Build ArtVu like any Go program. Or, use a pre-built binary from the /release folder.

## Known issues ##

- Viewing 80c art in 132c terminal mode: if the art was created at full 80c width (which is often the case), the file will contain "implicit" line endings (vs. actual line endings like "\n") so all the lines will be jammed together. There's no real fix for this at the moment, other than making sure art is saved to 79c instead of 80c.

- If the art file contains cursor codes/animation, it won't display properly.

- There's some degree of flicker when using the lighbar, see TO DO below...

## TO-DO ##
- [ ] Create a single object (struct) that contains directory/file hierarchy and SAUCE metadata - right now it scans only the "active" directory and creates 2 simple slices holding the directory and file names. Refactoring will allow for more granular sorting and filtering
- [ ] Add "Preview" data panel for terms wider than 80 cols
- [ ] Reduce flicker in dir list - Refactor lightbar function to *only* re-draw the entire list if scrolling beyond the terminal height. Otherwise, handle only affected rows
- [ ] Allow users to manually scroll up/down in ANSI art
- [ ] Add "-local" command line flag to allow users to browse locally - needs to handle/print ANSI differently
- [ ] Add Windows compatibility
- [ ] Add support for PgUp and PgDn keys/functions
- [ ] Add configurable delay (hard coded right now)