# Hyloblog

[![Hosted Platform](https://img.shields.io/badge/hosted%20platform)](https://hyloblog.com)
[![Discord Server](https://img.shields.io/badge/discord-gray?style=for-the-badge&logo=discord&logoColor=white)](https://discord.com/invite/E665nuukYn)

This is a repository containing the hyloblog software.

## Building

Some of the Make targets require sudo on Ubuntu, which is supplied in the
following way:

```bash
make up SUDO=sudo
```

To build statically: 

```bash
make up SUDO=sudo BUILDARGS=--static
```

## License and trademark

This repository contains the Hyloblog software, covered under the 
[Apache 2.0 License](LICENSE),
except where noted.

You are free to make your own distribution of the software, but you cannot use
any of the Hyloblog trademarks, as explained in
[our trademark policy](TRADEMARK.md).
