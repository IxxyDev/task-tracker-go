# Task Tracker CLI

A simple command-line task tracker written in Go. It allows you to add, list, update, and delete tasks directly from your terminal.

This project is designed to demonstrate core Go skills, including file system interaction, JSON handling, and command-line argument parsing, without relying on external libraries.

## Installation

1.  Ensure you have [Go](https://go.dev/doc/install) (version 1.18 or newer) installed on your system.
2.  Clone the repository:
    ```sh
    git clone <YOUR_REPOSITORY_URL>
    cd <YOUR_PROJECT_DIRECTORY>
    ```
3.  Build the project:
    ```sh
    go build -o task-cli
    ```
    This will create an executable file named `task-cli` in the project directory.

## Usage

The application is controlled via command-line arguments.

**Add a new task:**
```sh
./task-cli add "Buy milk"
```
Output: `Task added (ID: 1)`

**List tasks:**
```sh
# List all tasks
./task-cli list

# Filter tasks by status
./task-cli list done
./task-cli list in-progress
./task-cli list todo
```

**Update a task's description:**
```sh
./task-cli update 1 "Buy milk and bread"
```
Output: `Task 1 updated.`

**Update a task's status:**
```sh
# Mark task as done
./task-cli mark-done 1

# Mark task as in-progress
./task-cli mark-in-progress 1
```
Output: `Status of task 1 changed to 'done'.`

**Delete a task:**
```sh
./task-cli delete 1
```
Output: `Task 1 deleted.`

## Task Properties

Each task is stored as a JSON object with the following properties:
-   `id`: A unique integer identifier.
-   `description`: A string describing the task.
-   `status`: The current status (`todo`, `in-progress`, or `done`).
-   `createdAt`: The timestamp when the task was created.
-   `updatedAt`: The timestamp of the last modification.
