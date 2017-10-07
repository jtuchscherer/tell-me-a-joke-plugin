#!/bin/bash

set -e

(cf uninstall-plugin "tell-me-a-joke-plugin" || true) && go build -o tell-me-a-joke-plugin main.go && cf install-plugin tell-me-a-joke-plugin
