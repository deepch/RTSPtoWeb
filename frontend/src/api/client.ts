import axios from 'axios';

const client = axios.create({
  baseURL: '/',
  withCredentials: true,
});

export default client;
