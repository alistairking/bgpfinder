# bgpfinder

These are WIP notes about the bgpfinder project. Eventually this can be tidied up as user documentation.

## Technology

- Language: Go. It's what I know best these days, and it's easy to do HTTP things with it.
- DB: 
  - I'm open to suggestions
  - We don't need to pick just yet
  - Maybe a combo of redis and postgres?
  - Nothing newfangled.
- Build & Deploy:
  - Let's get Github Actions set up ASAP
  - Use to run tests and build docker images

## Packages

### base package: finder (used as a library)
- keep it simple. just support what we need. don’t make it too generic (hard coded stuff for a project is fine)
- handle only http urls (deal with kafka etc elsewhere, if you want a local archive, do it with object storage)
- defines Finder interface that other packages may implement
- just hit website and do parsing. it’s up to the caller to figure out lengths etc
- don’t worry much about edge cases where timestamps in files are off (leave it to the caller to ask for more time than they need if they care)
- intervals are inclusive,exclusive
  - e.g., midnight-1am returns 00,15,45 for rv collectors
- simple CLI that enables testing while development as well as use in scripts
  - e.g., find files matching a query and then pipe into xargs wget

#### Finder interface:
- list of projects
- list of collectors (ideally dynamic)
  - (optionally) overall approx time range for collector
- give me all the URLs for a project/collector/time window

### next package: server
- wrap finder in http server
- backwards compat with bgpstream broker API
- still no db. just fire up finder for every request
- designed to be run as sidecar to bgpstream instance(s)
- allows bgpstream to work without extra infra (or if broker infra goes down)
- hopefully will make production use more appealing

### and then: cache/database
- similar to current downloader
- runs as daemon
- uses finder package
- periodically schedules finds across proj/coll/time
  - schedule finds for recent windows regularly, less frequent for old windows
- compares found URLs to DB, updates accordingly.
- implements finder interface so that server and/or other embedded users can benefit from cache

### finally: public instance(s)
- run finder server in several stable places
- update bgpstream broker code to call finder server
- all bgpstream features should work in this way
