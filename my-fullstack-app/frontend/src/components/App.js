import React from 'react';
import BlockInfo from './BlockInfo';
import Login from './Login';
import { AuthProvider, useAuth } from '../contexts/AuthContext';
import '../styles/App.css';

// Protected component that only renders for authenticated users
const ProtectedContent = () => {
  const { currentUser } = useAuth();
  
  return (
    <>
      {currentUser ? (
        <BlockInfo />
      ) : (
        <div className="authentication-required">
          <h2>Authentication Required</h2>
          <p>Please sign in with Google to view Ethereum blockchain data.</p>
        </div>
      )}
    </>
  );
};

function App() {
  return (
    <AuthProvider>
      <div className="app">
        <header className="app-header">
          <h1>Ethereum Balance Tracker</h1>
        </header>
        
        <main className="app-content">
          <Login />
          <ProtectedContent />
        </main>
        
        <footer className="app-footer">
          <p>© 2025 Ethereum Balance Tracker</p>
        </footer>
      </div>
    </AuthProvider>
  );
}

export default App;