# `KVDB`

KVDB is a key-value database with support for multiple maps in one database.
It means you can create separate maps with their own keys. 
This approach is chosen to prevent multiple services to overwrite their keys.

## Installation

#### GitHub Releases

Download an archive from [GitHub Releases]("https://github.com/kostya-zero/kvdb/releases") and extract the KVDB binary to directory that is added to `PATH`.

#### Docker

Clone this repository and use `docker` CLI to build and run container.

```shell
docker build -t kvdb .
docker run -p 5511:5511 kvdb
```

## Usage

TBW

## API overview

TBW

## License

KVDB is licensed under MIT License. Learn more in [LICENSE](LICENSE) file.
