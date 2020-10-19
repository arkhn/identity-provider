import React from "react";
import "./App.css";
import axios from "axios";

import Home from "./components/home";

import { refreshToken, removeTokens } from "./oauth/tokenManager";
import { ACCESS_TOKEN_STORAGE_KEY, TOKEN_URL } from "./constants";

// Set axios interceptor
axios.interceptors.request.use((config) => {
  const accessToken = localStorage.getItem(ACCESS_TOKEN_STORAGE_KEY);
  config.headers.Authorization = `Bearer ${accessToken}`;
  return config;
});

// Add an interceptor to refresh access token when needed
axios.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    const originalRequest = error.config;

    if (
      error.response.status === 401 &&
      originalRequest.url.startsWith(TOKEN_URL)
    ) {
      removeTokens();
      return Promise.reject(error);
    }

    if (error.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      const success = await refreshToken();
      if (!success) {
        removeTokens();
        return Promise.reject(error);
      }
      return axios(originalRequest);
    }
    return Promise.reject(error);
  }
);

function App() {
  return (
    <Home />
  );
}

export default App;
