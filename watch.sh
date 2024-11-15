#!/usr/bin/env bash


# Directories to exclude
EXCLUDE_DIRS=(".githooks" ".github" "logs")

# Build the exclude patterns for inotifywait
EXCLUDE_PATTERN=$(printf "|%s" "${EXCLUDE_DIRS[@]}")
EXCLUDE_PATTERN=${EXCLUDE_PATTERN:1} # Remove leading |

# Kill background jobs on exit
cleanup() {
    if [[ -n $RUN_PID ]]; then
        kill $RUN_PID
        wait $RUN_PID 2>/dev/null
    fi
    exit
}
trap cleanup SIGINT SIGTERM

# Function to restart the program
restart_program() {
    echo "Change detected. Restarting program..."
    # Stop the running program
    if [[ -n $RUN_PID ]]; then
        kill $RUN_PID
        wait $RUN_PID 2>/dev/null
    fi

    # Rebuild and restart
    make build || { echo "Build failed"; return; }
    make run &
    RUN_PID=$!
}

# Start the program initially
restart_program

# Watch for file changes, ignoring excluded directories
while true; do
    CHANGED_FILE=$(inotifywait -re modify --exclude "(${EXCLUDE_PATTERN})" . 2>/dev/null)
    if [[ $? -eq 0 ]]; then
        restart_program
    fi
done
