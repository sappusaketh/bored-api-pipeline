#!/bin/bash

###############################################################################
# INSTALL DEVELOPMENT TOOLS
###############################################################################

# Default values
package_manager=""
repo_root=$(git rev-parse --show-toplevel) || { echo 'ERROR: Failed to get git repository root' >&2; return 1;}
scripts_home="${repo_root}/scripts"
# shellcheck source=/dev/null
source "${scripts_home}/lib/bash-logger.sh" || { echo 'ERROR: Failed to load bash logger' >&2; return 1;}

# Function to install tools using apt-get (for Ubuntu)
install_with_apt() {
  local tool="$1"
  sudo apt-get update
  sudo apt-get -y install "$tool"
}

# Function to install tools using brew (for macOS)
install_with_brew() {
  local tool="$1"
  brew install "$tool"
}


set_package_manager() {
  # Determine the OS type
  os_type=$(uname -s)
  local supported_package_managers=()
  if [[ "$os_type" == "Darwin" ]]; then
    # macOS
    supported_package_managers=("brew")
  elif [[ "$os_type" == "Linux" ]]; then
    # Linux (assuming Ubuntu)
    supported_package_managers=("apt-get")
  else
    WARNING "Unsupported OS type: $os_type"
  fi

  for pm in "${supported_package_managers[@]}"; do
    if command -v "$pm" >/dev/null 2>&1; then
      package_manager="$pm"
    fi
  done
  if [[ -z $package_manager ]]; then
    WARNING "one of supported_package_managers ${supported_package_managers[*]} is not installed. Missing packages won't be installed automatically."
  elif [[ -n $AUTO_INSTALL ]]; then
    NOTICE "Found supported package manager $package_manager. Missing packages will be installed automatically"
  fi
}


# Set the package manager
set_package_manager

# Check for pre-req tools
for _tool in docker make git minikube kubectl helm go
do
  if ! command -v "${_tool}" >/dev/null 2>&1; then
    WARNING "${_tool} is not installed or available in PATH." >&2
    # Install missing tool using the appropriate package manager
    # if you want to auto install on local then set AUTO_INSTALL to any value
    if [[ -n $AUTO_INSTALL ]]; then
      if [[ "$package_manager" == "brew" ]]; then
        install_with_brew "$_tool"
      elif [[ "$package_manager" == "apt-get" ]]; then
        install_with_apt "$_tool"
      else
        echo "Unsupported package manager: $package_manager"
        exit 1
      fi
    fi
  fi
done


# Done
exit 0

###############################################################################