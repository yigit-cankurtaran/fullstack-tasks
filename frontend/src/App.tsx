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

  const handleEditTask = (task: Task) => {
    // prompt user for new task name
    const name = prompt("Enter new task name:");
    if (!name) return alert("Task name cannot be empty!");

    // update task name on server
    fetch(`http://localhost:1239/tasks/${task.id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ ...task, name }),
    });
  };

  return (
    <div className="w-screen h-screen bg-black text-white">
      {/* create an input field for task name */}
      <div className="flex-col justify-center items-center text-center">
        <input type="text" className="text-black" />
        {/* create a button to add task */}
        <button
          onClick={handleAddTask}
          className="ml-4 bg-white text-black hover:bg-green-500 hover:text-white"
        >
          Add Task
        </button>
        <h1 className="text-4xl mb-8">Tasks</h1>
        <ul>
          {tasks.map((task) => (
            // <li key={task.id}>{task.name}</li>
            <li key={task.id} className="mb-4 text-xl">
              <input
                type="checkbox"
                checked={task.completion}
                onChange={() => handleCheckboxChange(task)}
                className="mr-4"
              />
              <span className={task.completion ? "line-through" : ""}>
                {task.id} - {task.name}
              </span>
              <button
                className="bg-white text-black ml-4 mr-4 rounded-sm p-0.5 hover:bg-red-500 hover:text-white"
                onClick={() => handleDeleteTask(task)}
              >
                Delete
              </button>
              <button
                className="bg-white text-black ml-4 mr-4 rounded-sm p-0.5 hover:bg-yellow-500 hover:text-white"
                onClick={() => handleEditTask(task)}
              >
                Edit
              </button>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
