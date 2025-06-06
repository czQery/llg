<div align="center" style="text-align: center">
  <img src="https://github.com/czQery/llg/blob/main/.github/banner.png?raw=true" alt="Logo">
  <br>
  
  [![Badge Report]][Report]
  [![Badge Actions]][Actions]
  [![Badge Release]][Release]
  
</div>

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
    # - ./users/:/data/data/users/
    # - ./devices/:/data/data/devices/
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

- pc1.log

```
login;PC (0.0.0.0);user1;01.01.1970;00:00
logoff;PC (0.0.0.0);user1;01.01.1970;00:01
login;PC (0.0.0.0);user2;02.01.1970;00:00
logoff;PC (0.0.0.0);user2;02.01.1970;00:01
```

<!----------------------------------------------------------------------------->

[Report]: https://goreportcard.com/report/github.com/czQery/llg/backend
[Badge Report]: https://goreportcard.com/badge/github.com/czQery/llg/backend

[Actions]: https://github.com/czQery/llg/actions
[Badge Actions]: https://img.shields.io/github/actions/workflow/status/czQery/llg/release.yml

[Release]: https://github.com/czQery/llg/releases/latest
[Badge Release]: https://img.shields.io/github/v/release/czQery/llg

<!----------------------------------------------------------------------------->
