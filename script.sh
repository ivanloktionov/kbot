#!/bin/bash

# MacOS
brew install gitleaks

# From Source
cd
git clone https://github.com/gitleaks/gitleaks.git
cd gitleaks
make build
