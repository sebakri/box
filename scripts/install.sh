#!/bin/sh
set -e

# Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "${OS}" in
  linux*)   OS='linux';;
  darwin*)  OS='darwin';;
  msys*|cygwin*|mingw*) OS='windows';;
  *)        echo "Unsupported OS: ${OS}"; exit 1;;
esac

# Detect Architecture
ARCH="$(uname -m)"
case "${ARCH}" in
  x86_64) ARCH='amd64';;
  arm64|aarch64) ARCH='arm64';;
  *)      echo "Unsupported architecture: ${ARCH}"; exit 1;;
esac

SUFFIX=""
if [ "${OS}" = "windows" ]; then
    SUFFIX=".exe"
fi

REPO="sebakri/box"
BINARY="box-${OS}-${ARCH}${SUFFIX}"

# Get latest release tag
echo "Detecting latest version..."
LATEST_TAG=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "${LATEST_TAG}" ]; then
    echo "Failed to detect latest version. Please check your internet connection or GitHub API limits."
    exit 1
fi

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}/${BINARY}"
INSTALL_DIR="/usr/local/bin"

# Allow overriding install directory
if [ -n "${BOX_INSTALL_DIR}" ]; then
    INSTALL_DIR="${BOX_INSTALL_DIR}"
fi

echo "Downloading box ${LATEST_TAG} for ${OS}/${ARCH}..."
curl -L -o "box" "${DOWNLOAD_URL}"
chmod +x "box"

if [ -w "${INSTALL_DIR}" ]; then
    mv "box" "${INSTALL_DIR}/box"
    echo "Successfully installed box to ${INSTALL_DIR}/box"
else
    echo "Install directory ${INSTALL_DIR} is not writable. Trying with sudo..."
    sudo mv "box" "${INSTALL_DIR}/box"
    echo "Successfully installed box to ${INSTALL_DIR}/box"
fi

echo "You can now run 'box' from your terminal."
