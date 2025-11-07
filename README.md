An expense tracker cli app for this project https://roadmap.sh/projects/task-tracker.


```bash
# Adding a new task
task add "Buy groceries"
# Output: Task added successfully (ID: 1)

# Updating and deleting tasks
task update 1 "Buy groceries and cook dinner"
task delete 1

# Marking a task as in progress or done
task mark doing 1
task mark done 3

# Listing all tasks
task list

# Listing tasks by status
task list done
task list todo
task list doing
```
