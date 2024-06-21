```
          ██              ████
        ██   ██ ███████ ██    ██
       █ ██  ██         █  █████
      ████ █            ██ ██ ██
       ██ ██  █        █  █ ███
         ██   █  ████  █   █ ██
        ███   ██       ███ █  █            ███
        █ ███             ███ █      █     █ ███ █
        ██          █       ███    ████    ██    ██
        █████████████████████ █    █  ██  ██     ███
        █ ██ █          ████  █    ██ ██      ██  █
        █   ██ ███████████    ██   ██  █         ██
        █      █████████       ████████ ████    █████
      ██                        ██    █      ██    ██
      █     █            ██      █             █ █ ██
      █     ██           █       █    ███       ██ █
      █      █          █        ██       ██     ██
      █       █        ██        ████       ██    █
      █       █        █        █ █   ███  ██ ██  ████
     ███      ██████████       █   █      █    █  █ ███
    █ ████  ████      ██     ██     ██          █ █
   █  ██ ██████        ██████        ████████   ██
  █    █       ████████    █          ██ █████  ██
  █     ███               █            █   ███   █
 ██     ██   ████████████ █            █ █████  ██
  █       ██              █            ███████ ██
  ██        █████      ████            █    ████
    ███       ████         █         ██   ██
   ███      ██  █   ██████████     ███████
 ██  ██  █  ████           ██         █
     █████                 █████████████
                            ███  ██   █
                               ██  ███
```

# Slowbro

An open source implementation of a [Turborepo](https://turbo.build/) remote cache server.

Slowbro is totally single tenant. It returns fake values for API calls and houses a single cache, a single user, a single team, and a single token.

## Demo

[![asciicast demo of Slowbro](https://asciinema.org/a/664941.svg)](https://asciinema.org/a/664941)

## Usage

To start Slowbro, you probably want to specify the `-token` flag at a minimum, otherwise you will be exposing your instance with the default token.

Slowbro can also be configured with a config file.

```
Usage of slowbro:
  -cache string
    	cache directory (default "./cache")
  -config string
    	config file (optional)
  -database string
    	sqlite database file (default "./metadata.db")
  -debug
    	enable debug middleware
  -listen string
    	listen address (default "localhost:8080")
  -login
    	enable login/link workflow
  -token string
    	API token to accept from turbo (default "867-5309")
```

## Turborepo

The easiest (and recommended) way to use Slowbro is to set a few environment variables where you run Turborepo. Set the following:

```
TURBO_API=http://your-slowbro-server-here
TURBO_TOKEN=your-token-here
TURBO_TEAM=literally-anything-you-want-its-ignored-but-required
```

### `-login`

Alternatively you can enable the `-login` flag on Slowbro and use the regular Turborepo workflow,

```
turbo login --login http://your-slowbro-server-here
turbo link --api http://your-slowbro--server-here
turbo run build --api http://your-slowbro-server-here
```

Keep in mind that while you are running with the `-login` flag on, _anyone_ can connect to your instance and get the token.  Using Slowbro this way will also clobber your Turborepo config, so you can not use it alongside another project which uses Vercel remote caching. The login mode is implemented for completeness, but is not recommended.
