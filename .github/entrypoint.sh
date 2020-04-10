#! /bin/sh
#
# Entry point for Github Action container.
set -e

IMPORT_PATH="${INPUT_IMPORT_PATH:-github.com/${GITHUB_REPOSITORY}}"
# Branch in push mode, or PR# in pull_request mode.
BRANCH=$(echo "${GITHUB_REF}" | cut -d/ -f3)
EMAIL="${INPUT_COMMIT-EMAIL:-posener@gmail.com}"

echo "Processing: ${IMPORT_PATH}@${BRANCH}
Event: ${GITHUB_EVENT_NAME}
"

# Run Goreadme on the current HEAD.
goreadme -import-path="${IMPORT_PATH}" $@ > ${README_FILE}

# Check if README was modified or was added, and don't push changes if nothing changed.
git add ${README_FILE}
if git diff --staged --exit-code --no-color ${README_FILE} > readme_diff.txt
then
    echo "No changes were made to ${README_FILE}, aborting"
    exit 0
fi

if [[ "${GITHUB_EVENT_NAME}" == "push" ]]
then
    echo "Push mode"

    # Configure git before commit.
    git config user.name Goreadme
    git config user.email "${EMAIL}"

    # Commit and push chnages to upstream branch.
    git commit -m "Update readme according to Go doc"
    git push origin HEAD:${BRANCH}
else
    echo "Pull request mode"

    # Prepare comment text.
    BODY=$(cat readme_diff.txt | tail +5 | sed "s/\`/'/g")

    echo "Diff:\n\n${BODY}"

    BODY="
[Goreadme](https://github.com/posener/goreadme) diff for \`${README_FILE}\` file for this PR:

\`\`\`diff
${BODY}
\`\`\`

This change will be automatically pushed when this PR is merged.
"
    
    # Add comment on PR if Github token was given.
    if [[ -z "${GITHUB_TOKEN}" ]]
    then
        echo "
In order to add request comment, set the github_token input.
${BODY}
"
        exit 0
    fi

    # Make the API call to post the comment.

    # Prepare the body for json:
    BODY=$(echo "${BODY}" | sed 's/\"/\\"/g' | sed ':a;N;$!ba;s/\n/\\n/g')
    
    curl "https://api.github.com/repos/${GITHUB_REPOSITORY}/pulls/${BRANCH}/reviews" \
        --fail \
        -H "Content-Type: application/json" \
        -H "Authorization: token ${GITHUB_TOKEN}" \
        -d "{
            \"event\": \"COMMENT\",
            \"body\": \"${BODY}\"
           }"
fi