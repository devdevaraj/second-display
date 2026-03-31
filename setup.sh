#!/bin/bash
set -e

echo "Updating package list..."
sudo apt-get update

echo "Installing dependencies for Virtual Display Management System..."
# sway: Headless compositor
# wf-recorder: Screen recording for wlroots-based compositors
# ffmpeg: Encoding and streaming
echo "Installing Sway, wf-recorder and FFmpeg..."
sudo apt-get install -y sway wf-recorder ffmpeg

# Build tools and basic requirements
echo "Installing build tools..."
sudo apt-get install -y build-essential pkg-config npm golang

echo "========================================="
echo "Dependencies installed successfully!"
echo "Note: It's recommended to install recent versions of Go and Node.js via their official managers (nvm, gvm) if Debian repos are outdated."
echo "========================================="
