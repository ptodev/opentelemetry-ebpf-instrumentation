# Working with Upstream Repositories

## Adding Upstream Repositories

When working with a fork of this repository, you may want to track changes from multiple upstream sources. Unlike managing the fork relationship itself (which is done through the GitHub UI), adding upstream repositories for local development is done using Git commands.

### Important Note: GitHub UI vs Git Commands

**You cannot add upstream repositories through the GitHub UI.** The GitHub UI only allows you to:

- Create a fork relationship (one-time action when forking)
- Sync your fork with the original repository (using the "Sync fork" button)

To work with multiple upstream repositories locally, you must use Git commands on your local machine.

## Adding the Grafana Fork as an Upstream Remote

If you want to track changes from the Grafana fork of this project, follow these steps:

### 1. Check Your Current Remotes

First, see what remotes you already have configured:

```sh
git remote -v
```

This will typically show your `origin` remote (your fork or the repository you cloned from).

### 2. Add the Grafana Repository as a Remote

Add the Grafana fork as a new remote called `grafana`:

```sh
git remote add grafana https://github.com/grafana/opentelemetry-ebpf-instrumentation.git
```

### 3. Fetch Changes from the Grafana Remote

Fetch the branches and commits from the Grafana repository:

```sh
git fetch grafana
```

### 4. Verify the Remote Was Added

Check that the remote was added successfully:

```sh
git remote -v
```

You should now see both your `origin` and the new `grafana` remote.

## Working with Multiple Upstreams

### Common Upstream Remotes for This Project

Depending on your workflow, you might want to track:

1. **open-telemetry/opentelemetry-ebpf-instrumentation** - The official OpenTelemetry repository
2. **grafana/opentelemetry-ebpf-instrumentation** - The Grafana fork with additional features/fixes

### Adding the Official OpenTelemetry Upstream

If you cloned from a fork and want to track the official OpenTelemetry repository:

```sh
git remote add upstream https://github.com/open-telemetry/opentelemetry-ebpf-instrumentation.git
git fetch upstream
```

### Syncing with Upstream Repositories

To pull changes from an upstream remote into your current branch:

```sh
# From the official OpenTelemetry upstream
git pull upstream main

# From the Grafana fork
git pull grafana main
```

To view changes without merging:

```sh
# View branches from a remote
git branch -r

# View commits from a specific remote branch
git log grafana/main

# Compare your branch with an upstream branch
git diff main..grafana/main
```

## Best Practices

1. **Naming Convention**: Use descriptive names for remotes:
   - `origin`: Your fork or the repository you push to
   - `upstream`: The official/primary upstream repository
   - `grafana`: The Grafana fork
   - Custom names as needed for other forks

2. **Regular Syncing**: Periodically fetch from all remotes to stay up-to-date:

   ```sh
   git fetch --all
   ```

3. **Cherry-picking**: If you want specific commits from an upstream:

   ```sh
   git cherry-pick <commit-hash>
   ```

4. **Creating Branches**: When working on features, create branches that track specific upstreams:

   ```sh
   git checkout -b feature-from-grafana grafana/main
   ```

## GitHub UI Features

While you cannot add remotes through the GitHub UI, you can:

- **Fork a repository**: Click the "Fork" button on the repository page
- **Sync your fork**: Use the "Sync fork" button to update your fork with the upstream default branch
- **Compare branches**: Use the compare view (e.g., `https://github.com/owner/repo/compare/branch1...owner2:repo:branch2`)
- **Create pull requests**: Between different forks and branches

## Troubleshooting

### "Remote already exists" Error

If you get an error that the remote already exists, you can:

- View existing remotes: `git remote -v`
- Remove a remote: `git remote remove <name>`
- Rename a remote: `git remote rename <old-name> <new-name>`
- Update a remote URL: `git remote set-url <name> <new-url>`

### Authentication Issues

When using HTTPS URLs, you may need to:

- Use a personal access token instead of password
- Configure Git credential helper
- Switch to SSH URLs: `git@github.com:grafana/opentelemetry-ebpf-instrumentation.git`

## See Also

- [Git Remote Documentation](https://git-scm.com/docs/git-remote)
- [GitHub: Working with Forks](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks)
- [CONTRIBUTING.md](../CONTRIBUTING.md) - General contribution guidelines
