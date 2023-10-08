<p align="center" style="text-align: center">
  <img src="https://github.com/czQery/llg/blob/main/.github/banner.png?raw=true" alt="Logo">
  <br>
  <a href="https://github.com/czQery/llg/actions">
    <img src="https://img.shields.io/github/actions/workflow/status/czQery/llg/release.yml" alt="build"/>
  </a>
  <a href="https://github.com/czQery/llg/releases/latest">
    <img src="https://img.shields.io/github/v/release/czQery/llg" alt="release"/>
  </a>
  <br>
</p>

# <a href="#installation" id="installation" name="installation">Installation</a>

### <a href="#installation-manual" id="installation-manual" name="installation-manual">Manual</a>

Linux:

```sh
curl -L https://github.com/czQery/llg/releases/latest/download/llg_linux_amd64.tar.gz | tar -xz
cd llg_linux_amd64
./llg
```

Windows:

```sh
curl.exe -L https://github.com/czQery/llg/releases/latest/download/llg_windows_amd64.zip
tar -xf .\llg_windows_amd64.zip
cd llg_windows_amd64
./llg.exe
```

### <a href="#installation-docker" id="installation-docker" name="installation-docker">Docker</a>

Docker Compose:

```yml
version: '3'
services:
  app:
    image: 'ghcr.io/czqery/llg:latest'
    container_name: llg
    restart: always
    volumes:
      - ./data/:/data/data/
    ports:
      - '8893:8893'
    tty: true
```

# <a href="#format" id="format" name="format">Data format</a>

- user1.log

```
login;user1;PC (0.0.0.0);01.01.1970;00:00
logoff;user1;PC (0.0.0.0);01.01.1970;00:01
```

- user2.log

```
login;user2;PC (0.0.0.0);01.01.1970;00:00
logoff;user2;PC (0.0.0.0);01.01.1970;00:01
```
