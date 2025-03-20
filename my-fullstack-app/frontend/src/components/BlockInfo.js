import React, { useState, useEffect } from 'react';
import { fetchLatestBlockNumber } from '../services/api';
import { useAuth } from '../contexts/AuthContext';
import '../styles/BlockInfo.css';

const BlockInfo = () => {
  const [blockNumber, setBlockNumber] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [lastUpdated, setLastUpdated] = useState(null);
  const { currentUser } = useAuth();

  const getBlockNumber = async () => {
    setLoading(true);
    try {
      const response = await fetchLatestBlockNumber();
      setBlockNumber(response.data);
      setLastUpdated(new Date());
      setError(null);
    } catch (err) {
      setError('Failed to fetch block number. Please try again later.');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    getBlockNumber();
    
    // Set up polling every 15 seconds
    const intervalId = setInterval(getBlockNumber, 15000);
    
    // Clean up interval on component unmount
    return () => clearInterval(intervalId);
  }, []);

  return (
    <div className="block-info-container">
      <h2>Ethereum Network Status</h2>
      
      {currentUser && (
        <p className="welcome-message">Welcome, {currentUser.displayName}!</p>
      )}
      
      {loading && blockNumber === null ? (
        <div className="loading">Loading latest block number...</div>
      ) : error ? (
        <div className="error">
          <p>{error}</p>
          <button onClick={getBlockNumber}>Try Again</button>
        </div>
      ) : (
        <div className="block-data">
          <div className="block-number">
            <span className="label">Latest Block:</span>
            <span className="value">{blockNumber.toLocaleString()}</span>
          </div>
          
          {lastUpdated && (
            <div className="update-time">
              Last updated: {lastUpdated.toLocaleTimeString()}
            </div>
          )}
          
          <button className="refresh-button" onClick={getBlockNumber}>
            Refresh
          </button>
        </div>
      )}
    </div>
  );
};

export default BlockInfo;