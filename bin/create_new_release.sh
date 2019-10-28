#!/bin/bash -xe
git pull --tags >/dev/null

echo "generate a rc build number"
BUMP_SEMVER_PATCH=$(git tag -l | grep -v "-" | tail -1 | awk -F. '{print $1"."$2"."$3+1}')
BUMP_SEMVER_RC=$(git tag -l | grep "${BUMP_SEMVER_PATCH}" | grep -e "-rc" | tail -1 | awk -F"-rc." '{print $2+1}')
SEMVER=${BUMP_SEMVER_PATCH}-rc.${BUMP_SEMVER_RC}
echo "tag id is: "${SEMVER}
echo "creating release"
github-release release -t ${SEMVER} -p
echo "uploading files"
for file in `ls build | grep '^hcunit'`
do
  github-release upload -t ${SEMVER} -f build/${file} -n ${file}
done
