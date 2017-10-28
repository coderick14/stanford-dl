![Built with love](http://forthebadge.com/images/badges/built-with-love.svg)
![Kinda SFW](http://forthebadge.com/images/badges/kinda-sfw.svg)
# stanford-dl

A simple multithreaded video/pdf downloader from ***Stanford Engineering Everywhere*** that *just works*.  

> NOTE : This tool is to be used strictly for educational purposes.

#### Installation
###### Build from source

```go
go get github.com/coderick14/stanford-dl
```
###### Download compiled binaries
Just download the binary for your required OS and Architecture from [releases](https://github.com/coderick14/stanford-dl/releases).

#### Usage
Type `stanford-dl --help` to get usage details. 
```
stanford-dl -course COURSE_CODE [-type {video|pdf}] [-all] [-lec lectures] [--help]

--course    Course name e.g. CS229, EE261
--type 	    Specify whether to download videos or pdfs. Defaults to PDF.
--all       Download for all lectures
--lec       Comma separated list of lectures e.g. 1,3,5,10
--help      Display this help message and quit
```

#### Examples
- Get all transcripts (PDFs) for a course
```
stanford-dl -course CS229 -type pdf -all
```
- Get only certain lectures
```
stanford-dl -course CS229 -type pdf -lec 1,3,5,10
```
- Get all videos for a course
```
stanford-dl -course CS229 -type video -all
```
- Get only certain lectures
```
stanford-dl -course CS229 -type video -lec 1,3,5,10
```

#### Contributions
This script has minimum functionality now. Please feel free to contribute for adding more features :)
