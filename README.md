# hugo-helper

[![wercker status](https://app.wercker.com/status/a783ec03618f8bcaadfb13c54b45c666/s/master "wercker status")](https://app.wercker.com/project/byKey/a783ec03618f8bcaadfb13c54b45c666)

Project hugo-helper is the consumer of several hugo related events.

It consumes these events from a rabbitMQ server and execute corresponding tasks.

These events include:

- New Site: Use `hugo new site <path>` command to create a new source directory;

- Generate: User `hugo [options]` command to (re-)generate a `public` directory from the source.
