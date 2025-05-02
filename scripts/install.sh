#!/bin/bash
set -e

# Install the RSS reader service and UI components
echo "Installing RSS Reader..."

# Determine the installation directory
INSTALL_DIR="$HOME/.rss-reader"
mkdir -p "$INSTALL_DIR/bin"

# Copy binaries
cp "build/rss-service" "$INSTALL_DIR/bin/"
cp "build/rss-ui" "$INSTALL_DIR/bin/"

# Create launch agent for auto-start
PLIST_PATH="$HOME/Library/LaunchAgents/com.yourusername.rss-reader.plist"

cat > "$PLIST_PATH" << EOL
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.yourusername.rss-reader</string>
    <key>ProgramArguments</key>
    <array>
        <string>${INSTALL_DIR}/bin/rss-service</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardErrorPath</key>
    <string>${INSTALL_DIR}/service-error.log</string>
    <key>StandardOutPath</key>
    <string>${INSTALL_DIR}/service-output.log</string>
</dict>
</plist>
EOL

# Load the launch agent
launchctl load "$PLIST_PATH"

# Create a script to launch the UI
cat > "$INSTALL_DIR/bin/launch-ui.sh" << EOL
#!/bin/bash
"$INSTALL_DIR/bin/rss-ui" "\$@"
EOL

chmod +x "$INSTALL_DIR/bin/launch-ui.sh"

# Create Applications symlink
mkdir -p "$HOME/Applications"
ln -sf "$INSTALL_DIR/bin/launch-ui.sh" "$HOME/Applications/RSS Reader.app"

echo "RSS Reader installed successfully!"
echo "The service is now running in the background."
echo "You can access the UI from the Applications folder or from the system tray."
