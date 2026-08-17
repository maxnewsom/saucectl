# `ai-authoring` (prototype)

Lets saucectl talk directly to the Sauce Labs AI Test Authoring API — generate
AI-authored test cases from a natural-language prompt, list/inspect/run them,
and pull the generated source code down to your workstation — without hand
writing curl calls.

**Status:** prototype, not wired into a saucectl release. It lives on this
branch/fork only; nothing here has shipped to `saucelabs/saucectl`.

## Building it

This isn't part of a saucectl release yet, so you build it from source and
run it as a separate binary rather than replacing your installed `saucectl`:

```
git clone https://github.com/maxnewsom/saucectl.git saucectl-proto
cd saucectl-proto
go build -o ./saucectl-dev ./cmd/saucectl
```

`saucectl-dev` is covered by the repo's `.gitignore` (`/saucectl*`), so it
never gets committed. Run it with `./saucectl-dev` from inside the repo; your
existing `saucectl` install is untouched.

## Credentials

Same resolution as regular saucectl — no extra setup needed if you're
already using saucectl:

1. `SAUCE_USERNAME` / `SAUCE_ACCESS_KEY` environment variables, or
2. `~/.sauce/credentials.yml` (written by `saucectl configure`)

## Commands

All commands live under `ai-authoring` and take a `-r/--region` flag
(`us-west-1` default, also `us-east-4`, `eu-central-1`).

### `generate` — create a test case from a prompt

```
saucectl-dev ai-authoring generate \
  --name "My first AI test" \
  --intent "Go to the homepage and verify the logo is visible" \
  --test-url "https://www.saucedemo.com" \
  --browser chrome --platform "Windows 11" --browser-version latest
```

Prints the task ID immediately, then polls (`--poll-interval`, default 5s)
until the task reaches `COMPLETED` or `FAILED`. Pass `--wait=false` to return
immediately after kicking off generation instead.

For mobile, omit `--test-url` and pass raw capabilities via `--capabilities`
(a JSON string), e.g.:

```
--capabilities '{"platformName":"Android","appium:automationName":"UiAutomator2","appium:deviceName":"Android GoogleAPI Emulator","appium:platformVersion":"13.0","appium:app":"storage:filename=myapp.apk"}'
```

Other flags: `--max-steps`, `--timeout`, `--test-suite-id`.

### `list` — list saved test cases

```
saucectl-dev ai-authoring list
saucectl-dev ai-authoring list --search checkout --limit 10
```

### `get` — get a single test case

```
saucectl-dev ai-authoring get <testCaseID>
```

### `run` — run a saved test case

```
saucectl-dev ai-authoring run <testCaseID> --build-name "manual test"
```

Prints the run ID and a dashboard URL for each job.

### `pull` — pull generated source code to disk

```
saucectl-dev ai-authoring pull <testCaseID>
```

Run without `--target` first to see the valid export targets for that test
case (they differ for web vs. mobile). Then:

```
saucectl-dev ai-authoring pull <testCaseID> --target javascript_webdriverio -o mytest.js
```

Omit `-o/--output` to print the code to stdout instead of writing a file.

## Known gaps

- No support yet for test suites, schedules, vault/variables, or tags —
  scoped out of this prototype to keep it small. Same client/command pattern
  extends to those if useful.
- The API's OpenAPI spec declares `bearerAuth` (JWT) as its security scheme,
  but the live API accepts HTTP Basic auth with username/access-key (same as
  every other saucectl-backed Sauce Labs API), which is what this
  implementation uses. Confirmed working against the live API.
