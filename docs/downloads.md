---
layout: dl
date: 2016-12-03T17:49:57+08:00
title: Downloads
permalink: /dl
---

{% for version in site.mid_versions %}
{% if version == site.mid_versions[0] %}
## Latest version: v{{ version }}
{% else %}
{% if version == site.mid_versions[1] %}
## Older versions

{% endif %}
### v{{ version }}
{% endif %}
##### Linux
* [mid{{ version }}.linux-386.tar.gz](https://github.com/midlang/mid/releases/download/v{{ version }}/mid{{ version }}.linux-386.tar.gz)
* [mid{{ version }}.linux-amd64.tar.gz](https://github.com/midlang/mid/releases/download/v{{ version }}/mid{{ version }}.linux-amd64.tar.gz)

##### Mac OS
* [mid{{ version }}.darwin-386.tar.gz](https://github.com/midlang/mid/releases/download/v{{ version }}/mid{{ version }}.darwin-386.tar.gz)
* [mid{{ version }}.darwin-amd64.tar.gz](https://github.com/midlang/mid/releases/download/v{{ version }}/mid{{ version }}.darwin-amd64.tar.gz)

##### Windows
* [mid{{ version }}.windows-386.zip](https://github.com/midlang/mid/releases/download/v{{ version }}/mid{{ version }}.windows-386.zip)
* [mid{{ version }}.windows-amd64.zip](https://github.com/midlang/mid/releases/download/v{{ version }}/mid{{ version }}.windows-amd64.zip)

##### Source code
* [mid{{ version }}.source.zip](https://github.com/midlang/mid/archive/v0.1.1.zip)
* [mid{{ version }}.source.tar.gz](https://github.com/midlang/mid/archive/v0.1.1.tar.gz)
{% endfor %}
