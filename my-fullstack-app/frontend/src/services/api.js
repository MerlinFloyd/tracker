/**
 * API Service for blockchain-related operations
 */

const API_URL = process.env.REACT_APP_API_URL || '/api';

/**
 * Fetches the latest Ethereum block number
 * @returns {Promise<Object>} The response with block number data
 * @throws {Error} If the API request fails
 */
export const fetchLatestBlockNumber = async () => {
  try {
    const response = await fetch(`${API_URL}/eth/block`);
    
    if (!response.ok) {
      throw new Error(`API error: ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error('Failed to fetch block number:', error);
    throw error;
  }
};

/**
 * Fetches balance for an Ethereum address
 * @param {string} address - The Ethereum address
 * @returns {Promise<Object>} The response with balance data
 * @throws {Error} If the API request fails
 */
export const fetchAddressBalance = async (address) => {
  try {
    const response = await fetch(`${API_URL}/eth/balance?address=${address}`);
    
    if (!response.ok) {
      throw new Error(`API error: ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error('Failed to fetch address balance:', error);
    throw error;
  }
};