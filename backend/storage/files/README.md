# File Storage
We will support the default file storage, but will be abstracting additional
functionality ontop of it and eventaully encapsulating it inside an archive
which we modify directly as a single file, instead of dealing with several
different files. This will allow us to have memory mapped portions, file IO
portions, and other interesting features. 
