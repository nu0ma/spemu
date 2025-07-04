# GitHub Repository Setup for Automated Releases

This repository uses automated versioning and releases. To enable this functionality, you need to configure your GitHub repository settings.

## Required GitHub Repository Settings

### 1. Enable GitHub Actions to create Pull Requests

1. Go to your repository on GitHub
2. Click on **Settings** tab
3. In the left sidebar, click **Actions** → **General**
4. Scroll down to **Workflow permissions**
5. Select **Read and write permissions**
6. Check **Allow GitHub Actions to create and approve pull requests**
7. Click **Save**

### 2. Branch Protection (Optional but recommended)

1. Go to **Settings** → **Branches**
2. Click **Add rule** for the `main` branch
3. Configure the following:
   - ✅ Require pull request reviews before merging
   - ✅ Require status checks to pass before merging
   - ✅ Require branches to be up to date before merging
   - ✅ Include administrators

## How the Automated Release System Works

1. **Development**: Create feature branches and use [Conventional Commits](https://www.conventionalcommits.org/)
2. **Merge to main**: When a PR is merged to main, release-please analyzes commits
3. **Release PR**: If there are releasable changes, a release PR is automatically created
4. **Release**: When the release PR is merged, a new version is tagged and released
5. **Binaries**: Cross-platform binaries are automatically built and attached to the release

## Commit Message Examples

```bash
feat: add new feature (minor version bump)
fix: resolve bug (patch version bump)
feat!: breaking change (major version bump)
docs: update documentation (no version bump)
ci: update GitHub Actions (no version bump)
```

## Manual Release (if needed)

If you need to manually trigger a release:

```bash
# Create a release PR manually
gh workflow run release-please.yml
```

## Troubleshooting

### Issue: "GitHub Actions is not permitted to create or approve pull requests"

**Solution**: Follow step 1 above to enable GitHub Actions to create pull requests.

### Issue: Release PR is not created

**Possible causes**:
- No releasable commits since last release
- Incorrect commit message format
- GitHub Actions permissions not configured

**Check**: Go to Actions tab and review the release-please workflow logs.