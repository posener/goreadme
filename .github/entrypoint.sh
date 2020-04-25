#! /bin/sh
#
# Entry point for Github Action container.
set -e

IMPORT_PATH="${INPUT_IMPORT_PATH:-github.com/${GITHUB_REPOSITORY}}"
# Branch in push mode, or PR# in pull_request mode.
BRANCH=$(echo "${GITHUB_REF}" | cut -d/ -f3)
EMAIL="${INPUT_COMMIT-EMAIL:-posener@gmail.com}"

REVIEWS_URL="https://api.github.com/repos/${GITHUB_REPOSITORY}/pulls/${BRANCH}/reviews"
COMMENT_HEADER="[Goreadme](https://github.com/posener/goreadme) diff for \`${README_FILE}\` file for this PR:"

echo "Processing: ${IMPORT_PATH}@${BRANCH}
Event: ${GITHUB_EVENT_NAME}
"

get_existing_comment_id() {
    CURRENT_REVIEWS="$(curl -sS "${REVIEWS_URL}")"
    LINE="$(echo "${CURRENT_REVIEWS}" | jq ".[] | .body" | grep -Fn "${COMMENT_HEADER}" | cut -d: -f1)"
    if [[ -n "${LINE}" ]]
    then
        echo "${CURRENT_REVIEWS}" | jq ".["$(( LINE - 1 ))"] | .id"
    fi
}

# Run Goreadme on the current HEAD.
goreadme -import-path="${IMPORT_PATH}" "$@" > "${README_FILE}"

# Check if README was modified or was added, and don't push changes if nothing changed.
git add "${README_FILE}"
if git diff --staged --exit-code --no-color "${README_FILE}" > readme_diff.txt
then
    echo "No changes were made to ${README_FILE}, aborting"
    exit 0
fi

# Diff content.
DIFF=$(cat readme_diff.txt | tail +5 | sed "s/\`/'/g")
echo "Diff:

${DIFF}
"

if [[ "${GITHUB_EVENT_NAME}" == "push" ]]
then
    echo "Push mode"

    # Configure git before commit.
    git config user.name Goreadme
    git config user.email "${EMAIL}"

    # Commit and push chnages to upstream branch.
    git commit -m "Update readme according to Go doc"
    git push origin HEAD:"${BRANCH}"
else
    echo "Pull request mode"
    # Add comment on PR if Github token was given.
    if [[ -z "${GITHUB_TOKEN}" ]]
    then
        echo "In order to add request comment, set the github_token input."
        exit 0
    fi

    # Prepare comment text.

    COMMENT_TEXT="${COMMENT_HEADER}

\`\`\`diff
${DIFF}
\`\`\`

This change will be automatically pushed when this PR is merged.
"

    # Prepare the comment text for json:
    COMMENT_TEXT=$(echo "${COMMENT_TEXT}" | sed 's/\"/\\"/g' | sed ':a;N;$!ba;s/\n/\\n/g')
    ACTION="POST"

    # Check if there is already a review that contains goreadme contnet
    COMMENT_ID="$(get_existing_comment_id)"
    if [[ -z "${COMMENT_ID}" ]]
    then
        echo "Creating new comment"
        CURL_BODY="{\"body\": \"${COMMENT_TEXT}\", \"event\": \"COMMENT\"}"
    else
        echo "Updating comment with ID: ${COMMENT_ID}"
        ACTION="PUT"
        CURL_BODY="{\"body\": \"${COMMENT_TEXT}\"}"
    fi
    
    # Perform comment post/update request.
    if ! curl -X "${ACTION}" "${REVIEWS_URL}/${COMMENT_ID}" \
        --fail -sS \
        -H "Content-Type: application/json" \
        -H "Authorization: token ${GITHUB_TOKEN}" \
        -d "${CURL_BODY}"
    then
        echo "Failed updating comment.

URL:    ${REVIEWS_URL}/${COMMENT_ID}
Action: ${ACTION}
Body:   ${CURL_BODY}
"
        exit 1
    fi
fi