# ImgPrint
A go tool for printing images to a terminal with support of 24 bit colors

Features include: 
  * Parallell processing of image data
  * Efficient for printing previews of images of any size, usually, in a fraction of a second
  * Printing images made easy for terminals that do not support image outputs

# Building

The only dependency is the go programming language, preferably version 1.21.7 or higher, as it makes use of experimental packages of the standard library

1. Clone the repository:
  ```bash
    git clone https://github.com/marattttt/imgprint && cd imgprint
  ```
2. Build the binary: 
  ```bash
    go build .
  ```
3. Run it!
  ```bash
    ./imgprint 3x3matrix.png
  ```
4. You probably want to move the executable into a folder that is included in your PATH env variable


