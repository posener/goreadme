#! /bin/sh

PACKAGE_NAME="github.com/${GITHUB_REPOSITORY}"
BRANCH=$(echo "${GITHUB_REF}" | cut -d/ -f3)
README="${INPUT_README_FILE}"
EMAIL="${INPUT_COMMIT-EMAIL:-posener@gmail.com}"

echo "Processing ${PACKAGE_NAME}@${BRANCH}"

# Run Goreadme on the current HEAD.
goreadme -package-name="${PACKAGE_NAME}" $@ > ${README}

# Check if README was modified, and don't push changes if nothing changed.
git add ${README}
if git diff --staged --exit-code
then
    echo "No changes were made to ${README}, aborting"
    exit 0
fi

# Configure git before commit.
git config user.name Goreadme
git config user.email "${EMAIL}"

# Commit and push chnages to upstream branch.
git commit -m "Update readme according to Go doc"
git push origin HEAD:${BRANCH}