# Where my space?
I recently had a need to pick through my hard drive to see what's eating up all my space.
So I made this to play around with different ways of finding out.

## Limitations
This really only works correctly on Windows, because I haven't taken the time to implement the same kind of syscall for other OSs.
Tests ***will*** fail on other OSs.

See [data_windows](data/data_windows.go) for details.

## Operations
There are a few things this can do.

* Show the total and free space on the current drive (Windows only).
* Show the minor, moderate, and major impacts to storage by walking directories starting from the target (or current) directory.
* List the files and directories in the target (or current) directory from biggest to smallest.
  * A result limit can be specified with something like `-n 5` to emit no more than 5 items.
