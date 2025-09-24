# chebread.github.io
This is my personal blog.

## How to Build
```shell
pnpm run build
```

## How to Run
### 1. Install Dependencies
```shell
go install github.com/cortesi/devd/cmd/devd@latest

pnpm install
```

### 2. Run the Development Server
```shell
pnpm run dev
```

## How to Deploy
To publish a new post or deploy changes, you must push a specially formatted Git tag. The deployment is handled automatically by GitHub Actions.

### 1. Commit Your Changes
First, make sure all your new content and changes are committed to the main branch.

```shell
# Stage all changes
git add .

# Commit the changes
git commit -m "post: foo: boo"

# Push the commit to the main branch
git push origin main
```

### 2. Create a Deploy Tag
The deployment workflow is triggered by a tag. Create a new tag based on the current date. The format is post/YYYY-MM-DD-XX, where XX is a two-digit number for posts on the same day (e.g., 01, 02).

```shell
# Example for the first post on September 23, 2025
git tag post/2025-09-23-01
```

### 3. Push the Tag to Deploy
Push the newly created tag to GitHub. This is the final step that will trigger the automatic build and deployment process.

```shell
git push origin post/2025-09-23-01
```

You can now go to the Actions tab in your GitHub repository to watch the deployment progress. Your site will be live with the new content in a few minutes.

## LICENSE
MIT LICENSE &copy; 2025 Cha Haneum