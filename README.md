# Purge

Inspired by macOS `purge` to drop the page cache, dentries and inodes

Build for the NVIDIA Jetson nano, compatible with most Linux distributions.

> **Note**
> The `build` and `install` targets require a working installation of [tinygo](https://tinygo.org/).  
> If you wish to build without tinygo, you can use the `build2` and `install2` targets.

## Install

```
git clone https://github.com/everdrone/purge.git
cd purge
make install
```

## Usage

```
sudo purge
```
