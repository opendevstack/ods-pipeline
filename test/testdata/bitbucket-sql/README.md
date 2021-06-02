# Bitbucket SQL

This SQL dump is from an "empty" Bitbucket instance.

It uses a [3 hour timebomb license](https://developer.atlassian.com/platform/marketplace/timebomb-licenses-for-testing-server-apps/).

There is only one user, `admin`, with password `admin`. For this user, a personal access token is configured with admin privileges.

There is one project, `ODSPIPELINETEST`.

## Generating a new dump

Start a postgres container like this:
```
docker run --name bitbucket-postgres -e POSTGRES_PASSWORD=jellyfish -e POSTGRES_USER=bitbucketuser -e POSTGRES_DB=bitbucket -d -p 5432:5432 postgres:12
```

Start a Bitbucket container like this:
```
docker run -d --name bitbucket-test -e ELASTICSEARCH_ENABLED=false -p 7990:7990 -p 7999:7999 atlassian/bitbucket:latest
```

After it has started, visit http://localhost:7990 and enter connect to the Postgres database. Then enter the timebomb license and create an admin user. Login with that admin user and create the PAT and a project.

Next, stop the Bitbucket container and enter the Postgres container, dump the database and copy it to your host machine:
```
docker exec -it bitbucket-postgres bash
$ pg_dump -U bitbucketuser bitbucket > bitbucket.sql
docker cp bitbucket-postgres:/bitbucket.sql .
```
