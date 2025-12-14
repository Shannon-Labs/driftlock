# Kaggle Dataset Benchmarks Demo
# Full benchmark with summary stats and Unicode box formatting

Output scripts/demos/kaggle_benchmarks.gif
Output scripts/demos/kaggle_benchmarks.mp4

Set Shell "zsh"
Set FontSize 14
Set Width 1400
Set Height 900
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

# Run the benchmark (using pre-built binary to avoid cargo output)
Type "./target/release/examples/kaggle_benchmarks"
Enter

# Wait for all benchmarks to complete (~40 seconds)
Sleep 45s

# Let final output settle
Sleep 3s
