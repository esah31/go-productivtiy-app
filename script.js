// DOM elements
const todoInput = document.getElementById("todo-input");
const todoList = document.getElementById("todo-list");

// Function to fetch tasks from server
async function fetchTasks() {
  try {
    const response = await fetch("/api/tasks");
    const tasks = await response.json();
    displayTasks(tasks);
  } catch (error) {
    console.error("Error fetching tasks:", error);
  }
}

// Function to display tasks
function displayTasks(tasks) {
  todoList.innerHTML = "";
  tasks.forEach((task) => {
    const listItem = document.createElement("li");
    listItem.innerHTML = `
            <span class="${task.completed ? "completed" : ""}">${
      task.text
    }</span>
            <button class="delete-btn" data-id="${
              task.id
            }"><i class="fas fa-trash"></i></button>
        `;
    if (!task.completed) {
      listItem.innerHTML += `<button class="complete-btn" data-id="${task.id}"><i class="fas fa-check"></i></button>`;
    }
    todoList.appendChild(listItem);
  });
}

// Function to add a new task
async function addTask(text) {
  try {
    const response = await fetch("/api/tasks", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ text }),
    });
    if (response.ok) {
      fetchTasks();
    } else {
      console.error("Error adding task:", response.statusText);
    }
  } catch (error) {
    console.error("Error adding task:", error);
  }
}

// Function to delete a task
async function deleteTask(id) {
  try {
    const response = await fetch(`/api/tasks/${id}`, {
      method: "DELETE",
    });
    if (response.ok) {
      fetchTasks();
    } else {
      console.error("Error deleting task:", response.statusText);
    }
  } catch (error) {
    console.error("Error deleting task:", error);
  }
}

// Function to complete a task
async function completeTask(id) {
  try {
    const response = await fetch(`/api/tasks/${id}/complete`, {
      method: "PUT",
    });
    if (response.ok) {
      fetchTasks();
    } else {
      console.error("Error completing task:", response.statusText);
    }
  } catch (error) {
    console.error("Error completing task:", error);
  }
}

// Event listener for adding a new task
document.getElementById("add-btn").addEventListener("click", function () {
  const text = todoInput.value.trim();
  if (text !== "") {
    addTask(text);
    todoInput.value = "";
  }
});

// Event delegation for deleting and completing tasks
todoList.addEventListener("click", function (event) {
  if (event.target.classList.contains("delete-btn")) {
    const id = event.target.dataset.id;
    deleteTask(id);
  } else if (event.target.classList.contains("complete-btn")) {
    const id = event.target.dataset.id;
    completeTask(id);
  }
});

// Fetch tasks when the page loads
fetchTasks();
