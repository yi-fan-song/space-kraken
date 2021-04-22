# Space-Kraken

Space-kraken is a client for playing the game [spacetraders.io](https://spacetraders.io) written in go.

# Releases

there may be a compiled binary for your operating system and architecture.

# Development

You need a working go development environment and [mage](https://github.com/magefile/mage).
Alternatively, you can build a docker container with the necessary dependencies.
Using a docker container is not necessary, but it helps ensure consistency of dependencies versions.

### Instructions for the use of the docker container

1. Build your image (or pull it from docker hub)
```bash
$ docker build -t space-kraken .
```
1. Run the container
```bash
$ docker run --rm -it \
	-p 8000:8000 \
	-v $PWD:/root/space-kraken \
	-w /root/space-kraken \
	--name space-kraken \
	space-kraken:latest
```

# Contributing

Contributions are welcome, a contributing guide is still TODO.

# License

Copyright © 2021 [Yi Fan Song](mailto:yfsong00@gmail.com)  

[GNU General Public License](https://www.gnu.org/licenses/)

