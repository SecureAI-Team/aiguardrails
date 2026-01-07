<template>
  <div id="app">
    <header class="topbar">
      <h1>AI Security Console</h1>
      <nav>
        <RouterLink to="/">Tenants</RouterLink>
        <RouterLink to="/apps">Apps</RouterLink>
        <RouterLink to="/policies">Policies</RouterLink>
        <RouterLink to="/policy-history">Policy History</RouterLink>
        <RouterLink to="/capabilities">Capabilities</RouterLink>
        <RouterLink to="/rules">Rules</RouterLink>
        <RouterLink to="/users">Users</RouterLink>
        <RouterLink to="/logs">Logs</RouterLink>
        <RouterLink to="/login">Login</RouterLink>
      </nav>
    </header>
    <main class="content">
      <RouterView />
    </main>
  </div>
</template>

<style>
body {
  margin: 0;
  font-family: Arial, sans-serif;
  background: #f6f8fb;
}
.topbar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 20px;
  background: #0f172a;
  color: #fff;
}
.topbar nav a {
  color: #cbd5f5;
  margin-right: 12px;
  text-decoration: none;
}
.topbar nav a.router-link-active {
  color: #fff;
  font-weight: bold;
}
.content {
  padding: 20px;
}
.card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  padding: 16px;
  margin-bottom: 12px;
}
form {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}
input {
  padding: 6px 10px;
  border: 1px solid #cbd5e1;
  border-radius: 4px;
}
button {
  padding: 6px 12px;
  border: none;
  background: #2563eb;
  color: #fff;
  border-radius: 4px;
  cursor: pointer;
}
button:disabled {
  background: #94a3b8;
}
</style>

