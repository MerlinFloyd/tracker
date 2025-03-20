import React, { useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import '../styles/Login.css';

const Login = () => {
  const { currentUser, signInWithGoogle, logout } = useAuth();
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleGoogleSignIn = async () => {
    try {
      setError('');
      setLoading(true);
      await signInWithGoogle();
    } catch (error) {
      setError('Failed to sign in with Google: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = async () => {
    try {
      setError('');
      await logout();
    } catch (error) {
      setError('Failed to log out: ' + error.message);
    }
  };

  return (
    <div className="login-container">
      {currentUser ? (
        <div className="user-profile">
          <div className="profile-info">
            <img 
              src={currentUser.photoURL} 
              alt="Profile" 
              className="profile-image"
            />
            <div className="user-details">
              <h3>{currentUser.displayName}</h3>
              <p>{currentUser.email}</p>
            </div>
          </div>
          <button 
            onClick={handleLogout} 
            className="logout-button"
            disabled={loading}
          >
            Sign Out
          </button>
        </div>
      ) : (
        <div className="login-card">
          <h2>Sign In</h2>
          {error && <div className="error-message">{error}</div>}
          <button 
            onClick={handleGoogleSignIn} 
            className="google-button"
            disabled={loading}
          >
            <img 
              src="https://upload.wikimedia.org/wikipedia/commons/5/53/Google_%22G%22_Logo.svg" 
              alt="Google logo" 
            />
            Sign in with Google
          </button>
        </div>
      )}
    </div>
  );
};

export default Login;