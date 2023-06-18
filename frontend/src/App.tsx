import { useState, useEffect } from "react";

interface Task {
  id: number;
  name: string;
  completion: boolean;
}

export default function App() {
  const [tasks, setTasks] = useState<Task[]>([]);

  useEffect(() => {
    fetch("http://localhost:1239/tasks")
      .then((response) => response.json())
      .then((data) => setTasks(data));
  }, []);

  const handleCheckboxChange = (task: Task) => {
    // update task completion
    const updatedTask = { ...task, completion: !task.completion };
    setTasks((tasks) => tasks.map((t) => (t.id === task.id ? updatedTask : t)));

    // update task completion on server
    fetch(`http://localhost:1239/tasks/${task.id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(updatedTask),
    });
  };

  const handleAddTask = () => {
    // get task name from input field
    const name = (document.querySelector("input") as HTMLInputElement).value;
    if (!name) return alert("Task name cannot be empty!");

    // add task to server
    fetch("http://localhost:1239/tasks", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ id: tasks.length + 1, name, completion: false }),
    });
  };

  const handleDeleteTask = (task: Task) => {
    // delete task from server
    fetch(`http://localhost:1239/tasks/${task.id}`, {
      method: "DELETE",
    });
  };

  return (
    <div>
      {/* create an input field for task name */}
      <input type="text" />
      {/* create a button to add task */}
      <button onClick={handleAddTask}>Add Task</button>
      <h1>Tasks</h1>
      <ul>
        {tasks.map((task) => (
          // <li key={task.id}>{task.name}</li>
          <li key={task.id}>
            <input
              type="checkbox"
              checked={task.completion}
              onChange={() => handleCheckboxChange(task)}
            />
            {task.id} - {task.name}
            <button onClick={() => handleDeleteTask(task)}>Delete</button>
          </li>
        ))}
      </ul>
    </div>
  );
}
