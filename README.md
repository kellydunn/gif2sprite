# gif2sprite

A command line interface for splitting animated gifs into sprite sheets and associated metadata for animations, framing, and other video effects.

This tool is written in the [`golang`](http://golang.org) programming language, sheerly out of personal preference.  

## installation

The simplest way to download this tool at the moment is to use the `go` command line tool to install the binary on your system:

```
$ go install github.com/kellydunn/gif2sprite
```

In the future, I will aim to upload binary realizes with each git release.

## usage

The command essentially splits each frame of an animated gif and stiches them into a single `.png` sprite with a transparent background.  It also extracts the metadata from each frame and places it in a `.json` file.

To do this, issue the following command:

```
$ gif2sprite --input-dir=<source directory> --output-dir=<output directory>
```

Assuming your source directory looks like this:

```
<source directory>
 ├── gif1.gif
 └── gif2.gif
```

The result of the command above will create the following folder structure:

```
<output directory>
 ├── gif1
 │   ├── gif1.json
 │   └── gif1.png
 └── gif2
     ├── gif2.json
     └── gif2.png
```