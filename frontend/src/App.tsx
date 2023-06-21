import { useState, useEffect } from "react";

interface Task {
  id: number;
  name: string;
  completion: boolean;
}

export default function App() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [darkMode, setDarkMode] = useState<boolean>(false);

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
    }).then(() => {
      // refetch tasks from server
      fetch("http://localhost:1239/tasks")
        .then((response) => response.json())
        .then((data) => setTasks(data));
    });
  };

  const handleDeleteTask = (task: Task) => {
    // delete task from server
    fetch(`http://localhost:1239/tasks/${task.id}`, {
      method: "DELETE",
    }).then(() => {
      // refetch tasks from server
      fetch("http://localhost:1239/tasks")
        .then((response) => response.json())
        .then((data) => setTasks(data));
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
    }).then(() => {
      // refetch tasks from server
      fetch("http://localhost:1239/tasks")
        .then((response) => response.json())
        .then((data) => setTasks(data));
    });
  };

  return (
    <div
      className={`w-screen h-screen ${
        darkMode ? "bg-gray-900 text-gray-100" : "bg-gray-100 text-gray-800"
      }`}
    >
      <div className="max-w-2xl mx-auto py-10">
        <h1 className="text-4xl font-bold mb-8">Tasks</h1>
        <div className="flex mb-4">
          <input
            type="text"
            className={`flex-grow rounded border ${
              darkMode ? "border-gray-600" : "border-gray-300"
            } p-2 mr-2 ${
              darkMode ? "bg-gray-800 text-gray-100" : "bg-white text-gray-800"
            }`}
          />
          <button
            onClick={() => {
              handleAddTask();
              (document.querySelector("input") as HTMLInputElement).value = "";
            }}
            className={`${
              darkMode ? "bg-blue-400" : "bg-blue-500"
            } text-white rounded px-4 py-2 hover:${
              darkMode ? "bg-blue-500" : "bg-blue-600"
            }`}
          >
            Add Task
          </button>
        </div>
        <ul
          className={`divide-y ${
            darkMode ? "divide-gray-700" : "divide-gray-300"
          }`}
        >
          {tasks.map((task) => (
            <li key={task.id} className="py-4 flex items-center">
              <input
                type="checkbox"
                checked={task.completion}
                onChange={() => handleCheckboxChange(task)}
                className="mr-4"
              />
              <span
                className={`flex-grow text-lg ${
                  task.completion ? "line-through text-gray-500" : ""
                }`}
              >
                {task.name}
              </span>
              <button
                className={`${
                  darkMode
                    ? "bg-gray-800 text-gray-100"
                    : "bg-white text-gray-800"
                } ml-2 rounded p-2 hover:bg-red-500 hover:text-white`}
                onClick={() => handleDeleteTask(task)}
              >
                Delete
              </button>
              <button
                className={`${
                  darkMode
                    ? "bg-gray-800 text-gray-100"
                    : "bg-white text-gray-800"
                } ml-2 rounded p-2 hover:bg-yellow-500 hover:text-white`}
                onClick={() => handleEditTask(task)}
              >
                Edit
              </button>
            </li>
          ))}
        </ul>
      </div>
      <button
        className={`fixed bottom-4 left-4 ${
          darkMode ? "bg-gray-800 text-gray-100" : "bg-white text-gray-800"
        } rounded-full p-4 shadow-lg hover:${
          darkMode ? "bg-gray-700" : "bg-gray-100"
        }`}
        onClick={() => setDarkMode(!darkMode)}
      >
        {darkMode ? "Light" : "Dark"}
      </button>
    </div>
  );
}
