import React from 'react';
import BlockInfo from './BlockInfo';
import Login from './Login';
import { AuthProvider } from '../contexts/AuthContext';
import '../styles/App.css';

function App() {
  return (
    <AuthProvider>
      <div className="app">
        <header className="app-header">
          <h1>Ethereum Balance Tracker</h1>
        </header>
        
        <main className="app-content">
          <Login />
          <BlockInfo />
        </main>
        
        <footer className="app-footer">
          <p>© 2025 Ethereum Balance Tracker</p>
        </footer>
      </div>
    </AuthProvider>
  );
}

export default App;