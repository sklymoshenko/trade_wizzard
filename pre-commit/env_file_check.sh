#!/bin/bash

# Function to check for changes in .env file
check_env_changes() {
    # Check if .env file is tracked by git and has changes
    if git status --short | grep -q "^.M .env"; then
        echo "Warning: You have uncommitted changes in your .env file!"
        exit 1
    fi

    # Check if .env file is staged for commit
    if git status --short | grep -q "^M  .env"; then
        echo "Warning: You have staged changes in your .env file!"
        exit 1
    fi

    # Check if .env file is untracked
    if git status --short | grep -q "^?? .env"; then
        echo "Warning: You have an untracked .env file!"
        exit 1
    fi

    echo ".env file is clean."
}

# Run the check
check_env_changes