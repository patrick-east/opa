#!/bin/bash

set -xe

RELEASES=$(cat RELEASES)

ORIGINAL_COMMIT=$(git rev-parse --abbrev-ref HEAD)
ROOT_DIR=$(git rev-parse --show-toplevel)
RELEASES_YAML_FILE=${ROOT_DIR}/docs/data/releases.yaml
GIT_VERSION=$(git --version)

echo "Git version: ${GIT_VERSION}"

echo "Saving current workspace state"
STASH_TOKEN=$(uuidgen)
git stash push --include-untracked -m "${STASH_TOKEN}"

function restore_tree {
    echo "Returning to commit ${ORIGINAL_COMMIT}"
    git checkout ${ORIGINAL_COMMIT}

    # Only pop from the stash if we had stashed something earlier
    if [[ -n "$(git stash list | head -1 | grep ${STASH_TOKEN} || echo '')" ]]; then
        git stash pop
    fi
}

function cleanup {
    EXIT_CODE=$?

    if [[ "${EXIT_CODE}" != "0" ]]; then 
        # on errors attempt to restore the starting tree state
        restore_tree

        echo "Error loading docs"
        exit ${EXIT_CODE}
    fi

    echo "Docs loading complete"
}

trap cleanup EXIT

echo "Cleaning generated folder"
rm -rf ${ROOT_DIR}/docs/generated/*

echo "Removing data/releases.yaml file"
rm -f ${RELEASES_YAML_FILE}

for release in ${RELEASES}; do
    version_docs_dir=${ROOT_DIR}/docs/generated/docs/${release}

    mkdir -p ${version_docs_dir}

    echo "Adding ${release} to releases.yaml"
    echo "- ${release}" >> ${RELEASES_YAML_FILE}

    echo "Copying doc content from tag ${release}"
    git checkout ${release}

    # TODO: Remove this check once we are no longer showing docs for v0.10.7 
    # or older those releases have the docs content in a different location.
    if [[ -d "${ROOT_DIR}/docs/content/code/" ]]; then
        # new location
        cp -r ${ROOT_DIR}/docs/content/* ${version_docs_dir}/
    else
        # old location
        cp -r ${ROOT_DIR}/docs/content/docs/* ${version_docs_dir}/
        cp -r ${ROOT_DIR}/docs/code ${version_docs_dir}/
    fi
done

# Go back to the original tree state
restore_tree

# Create the "edge" version from current working tree
echo 'Adding "edge" to releases.yaml'
echo "- edge" >> ${RELEASES_YAML_FILE}

# Link instead of copy so we don't need to re-generate each time.
# Use a relative link so it works in a container more easily.
ln -s ../../content ${ROOT_DIR}/docs/generated/docs/edge
