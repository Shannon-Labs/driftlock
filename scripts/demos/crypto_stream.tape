# Crypto Streaming Anomaly Detection Demo
# Produces video-friendly output with Unicode boxes, checkmarks, and progress indicators

Output scripts/demos/crypto_stream.gif
Output scripts/demos/crypto_stream.mp4

Set Shell "zsh"
Set FontSize 16
Set Width 1200
Set Height 700
Set Theme "Dracula"
Set TypingSpeed 40ms
Set Padding 20
Set Framerate 30
Set PlaybackSpeed 1.0

# Start clean
Hide
Type "clear"
Enter
Show
Sleep 500ms

# Run the demo (using pre-built binary to avoid cargo output)
Type "./target/release/examples/crypto_stream"
Enter

# Wait for demo execution (runs in ~17 seconds)
Sleep 20s

# Let output settle
Sleep 2s
