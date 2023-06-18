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

  return (
    <ul>
      {tasks.map((task) => (
        // <li key={task.id}>{task.name}</li>
        <li key={task.id}>
          <input
            type="checkbox"
            checked={task.completion}
            onChange={() => handleCheckboxChange(task)}
          />
          {task.name}
        </li>
      ))}
    </ul>
  );
}
