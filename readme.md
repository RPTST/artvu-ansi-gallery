# Artvu Ansi Gallery
ANSI art viewer for old-skool, and new-school, BBSs (like Mystic, Talisman, WWiV, ENiGMA 1/2, etc.). Supports 80 and 132 terminal widths. Best run as an stdout-based (linux native) door program - refer to the example start.sh file.

Required command-line flag is "-root" followed by the path to the folder containing Ansi art:

./artvu-ansi-gallery -root /path/to/art

Everything under this root will be viewable by users, so be aware.

Artvu scans for files and folders that contain *.ans, *.asc or *.diz files under the root and turns each directory that contains those files into a listing tree that can be navigated via a lightbar over telnet/ssh connection (stdout) using a terminal program like SyncTerm, MagiTerm or NetRunner. 

It incorporates SAUCE metadata into the file listings. 

Linux, Mac, RPi compatible. No Windows support -- yet.

## Build ##
Build ArtVu like any Go program. Or, use a pre-built binary from the /release folder.

## Known issues ##

- Viewing 80 column art in 132 column terminal mode: if the .ans file was created at full width (80), then the file contains "implicit" line endings (vs. actual line endings like "\n") so all the lines will be jammed together when viewed at larger widths. There's no real fix for this at the moment, other than making sure art is saved to 79c instead of 80c.

- If the art file contains cursor codes/animation, it won't display properly.

- There's some degree of flicker when using the lighbar, see TO-DO below...

## TO-DO ##
- [ ] Create a single object (struct) that contains directory/file hierarchy and SAUCE metadata - right now it scans only the "active" directory and creates 2 simple slices holding the contents, one for files and the other for directories. Refactoring to a single object will allow for more granular displaying, sorting and filtering
- [ ] Add "Preview" data panel for terms wider than 80 cols (blank/placeholder for now)
- [ ] Reduce lightbar flicker in list - *only* re-draw the list if scrolling outside the terminal height. Otherwise, just handle the affected rows
- [ ] Allow users to manually scroll up/down in ANSI art
- [ ] Add "-local" command line flag to allow users to browse locally - but needs to handle/print ANSI differently
- [ ] Add Windows compatibility
- [ ] Add support for PgUp and PgDn keys/functions
- [ ] Add configurable delay (hard coded right now)