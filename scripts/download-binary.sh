platform='*'
unamestr=$(uname)
if [[ "$unamestr" == 'Linux' ]]; then
   platform='linux'
elif [[ "$unamestr" == 'Windows' ]]; then
   platform='Windows'
elif [[ "$unamestr" == 'Darwin' ]]; then
   platform='Darwin'
fi

ARCH=$(arch)

# Get the latest release version
VERSION=$(curl https://api.github.com/repos/kilianstallz/stage-sync/releases | jq '.[0]' | grep 'tag_name' | grep -Eo 'v[^\"]*')

# Get the latest release download url for the current os and architecture
DL=$(curl https://api.github.com/repos/kilianstallz/stage-sync/releases | jq '.[0].assets' | grep "browser_download_url" | grep -Eo 'https://[^\"]*' | grep "$platform" | grep "$ARCH" | sed -n '1p')

# Download the latest release and extract
wget -qO- "$DL" | tar -xz