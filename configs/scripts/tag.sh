
#!/bin/sh

set -e
echo "Fetching tags..."
git fetch --tags

# Get the most recent tag for this component
LATEST_TAG=$(git tag -l "apps/backend/*" --sort=-version:refname | head -n 1 || echo "")
echo "Latest backend tag: ${LATEST_TAG:-none}"

# Get the list of files changed since the latest tag (or all files if no tag exists)
if [ -n "$LATEST_TAG" ]; then
CHANGED_FILES=$(git diff --name-only "${LATEST_TAG}" HEAD)
echo "Changed files since ${LATEST_TAG}:"
else
CHANGED_FILES=$(git ls-files)
echo "No previous tag found. Checking all files:"
fi
echo "$CHANGED_FILES"

# Check if any changed file is in /api or /apps/backend
if echo "$CHANGED_FILES" | grep -Eq '^(api/|apps/backend/)'; then
	echo "Changes detected in /api or /apps/backend. Proceeding with tagging."

	# Read version from component.json
	COMPONENT_JSON="apps/backend/component.json"

	if [ ! -f "$COMPONENT_JSON" ]; then
	echo "Error: $COMPONENT_JSON not found!"
	exit 1
	fi

	# Extract version from component.json using jq
	VERSION=$(jq -r '.version' "$COMPONENT_JSON")

	if [ -z "$VERSION" ] || [ "$VERSION" = "null" ]; then
	echo "Error: No version field found in $COMPONENT_JSON"
	exit 1
	fi

	echo "Version from component.json: $VERSION"

	# Create new tag with format apps/backend/vX.Y.Z
	NEW_TAG="apps/backend/v${VERSION}"
	echo "New tag will be: $NEW_TAG"

	# Check if tag already exists
	if git rev-parse "$NEW_TAG" >/dev/null 2>&1; then
	echo "Tag $NEW_TAG already exists. Skipping tag creation."
	echo "If you want to create a new tag, update the version in $COMPONENT_JSON"
	exit 0
	fi

	# Create and push the new tag
	git tag "$NEW_TAG"
	git push origin "$NEW_TAG"
	echo "Successfully created and pushed tag: $NEW_TAG"
else
	echo "No changes in /api or /apps/backend. No tagging necessary."
fi