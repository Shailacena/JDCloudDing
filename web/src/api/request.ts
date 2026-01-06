import axios from 'axios';

export const request = axios.create({
  baseURL: `${window.location.origin}/web_api/`,
  timeout: 60000,
  headers: {
    'Content-Type': 'application/json'
  },
});