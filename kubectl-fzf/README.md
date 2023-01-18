## Fuzzy search with Mac background task

1. Edit the `plist` file for username and path to `sh` file
2. Move the `plist` file to `~/Library/LaunchAgents/com.kubectl.fzf-server.plist`
3. Run `launchctl load ~/Library/LaunchAgents/com.kubectl.fzf-server.plist`
4. Check if it's already running with `launchctl list | grep com.ku` (should have pid in first column, not "-")