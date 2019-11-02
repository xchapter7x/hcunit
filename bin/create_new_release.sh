#!/bin/bash -xe
if [ -z ${GITHUB_REPO+x} ]; then echo "GITHUB_REPO is unset"; exit 1; fi
if [ -z ${GITHUB_TOKEN+x} ]; then echo "GITHUB_TOKEN is unset"; exit 1; fi
if [ -z ${GITHUB_USER+x} ]; then echo "GITHUB_USER is unset"; exit 1; fi
git pull --tags >/dev/null || true

echo "generate a rc build number"
BUMP_SEMVER_PATCH=$(git tag -l | grep -v "-" | tail -1 | awk -F. '{print $1"."$2"."$3+1}')
if [[ ${BUMP_SEMVER_PATCH} -eq "" ]]; then
  BUMP_SEMVER_PATCH="0.1.0"
fi
BUMP_SEMVER_RC=$(git tag -l | grep "${BUMP_SEMVER_PATCH}" | grep -e "-rc" | tail -1 | awk -F"-rc." '{print $2+1}')
if [[ ${BUMP_SEMVER_RC} -eq "" ]]; then
  BUMP_SEMVER_RC="0"
fi
SEMVER=${BUMP_SEMVER_PATCH}-rc.${BUMP_SEMVER_RC}
echo "tag id is: "${SEMVER}
echo "creating plugin tarball"
tar -czvf build/hcunit_plugin.tgz --directory=build hcunit_osx hcunit.exe hcunit_unix plugin.yaml
echo "creating release"
github-release release -t ${SEMVER} -p
echo "uploading files"
for file in `ls build | grep '^hcunit'`
do
  github-release upload -t ${SEMVER} -f build/${file} -n ${file}
done
