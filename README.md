# jot

Quick note taking tool linked with Git.

## About

In a way to better understand Git and its internals I used git2go instead of other alternatives. This may pose an installation issue moving forward but we'll see.

## Dependencies

### Ubuntu

Following packages are required: `sudo apt-get install libgpm2 libgpm-dev libgit2-glib-1.0-dev`

#### Lib2Git from Source

```
sudo apt-get -y install cmake libssl-dev
wget https://github.com/libgit2/libgit2/releases/download/v1.1.0/libgit2-1.1.0.tar.gz
tar -xzf libgit2-1.1.0.tar.gz
cd libgit2-1.1.0 && mkdir build && cd build && cmake .. -DCMAKE_INSTALL_PREFIX=/usr && sudo cmake --build . --target install
rm -rf libgit2*
```

## Resource

- [Git Plumbing and Porcelain](https://git-scm.com/book/en/v2/Git-Internals-Plumbing-and-Porcelain)
- [git2go](https://github.com/libgit2/git2go)
