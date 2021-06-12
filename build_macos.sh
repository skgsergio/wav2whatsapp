#!/bin/sh

APP="wav2whatsapp.app"

# Copy Info.plist
mkdir -p "$APP/Contents"
cp resources/Info.plist "$APP/Contents/Info.plist"

# Copy icon.icns
mkdir -p "$APP/Contents/Resources"
cp resources/icon.icns "$APP/Contents/Resources/icon.icns"

# Build app
mkdir -p "$APP/Contents/MacOS"
GOOS=darwin go build -o "$APP/Contents/MacOS/wav2whatsapp"

# Create tar
tar cvf wav2whatsapp_macos.tar "$APP"
