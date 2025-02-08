# Hyloblog

[Platform](https://hyloblog.com) | [Blog](https://blog.hyloblog.com) |  [Discord](https://discord.com/invite/E665nuukYn)

This is a repository containing the Hyloblog software.

## About

<!-- basic description -->
Hyloblog is a static site generator. The name derives from ὕλη, the Greek word
for _matter_.

<!-- how it works -->
Hyloblog generates blog sites (like [this one][blog]) from Markdown files in
Git repos using simple default rules to limit metadata in source files, such as

- The highest-level title in an `*.md` source file is treated as the title of
  the corresponding blog post

- The date and author of the commit when an `*.md` file is first added are
  the date and author of the post

- The path to a given file within the repo is mirrored exactly in the generated
  files, so if you want an `/updates` section in your blog you simply make a
  directory and start writing updates.

  [blog]: https://github.com/hyloblog/blog

<!-- influences -->
Although Hyloblog learns from static-site generators like Jekyll and Hugo, an
important influence is LaTeX.
We want to emphasise the same form/content separation in blogging and site
building that LaTeX accomplishes in typesetting.

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
