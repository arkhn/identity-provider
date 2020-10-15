import React, { useEffect, useState } from "react";
import queryString from "query-string";
import { v4 as uuid } from "uuid";
import { Button, CircularProgress } from "@material-ui/core";

import {
  fetchTokens,
  getIdToken,
  revokeToken,
  removeTokens,
} from "../../oauth/tokenManager";
import { authClient } from "../../oauth/authClient";

import {
  ACCESS_TOKEN_STORAGE_KEY,
  LOGOUT_URL,
  LOGOUT_REDIRECT_URL,
  STATE_STORAGE_KEY,
} from "../../constants";

const startAuthentication = () => {
  const state = uuid();
  localStorage.setItem(STATE_STORAGE_KEY, state);
  const uri = authClient.code.getUri({
    state: state,
  });
  window.location.assign(uri);
};

const logout = () => {
  const idToken = getIdToken();
  if (!idToken) throw new Error("Can't logout, id token not found.");
  const logoutUrl = `${LOGOUT_URL}?id_token_hint=${idToken}&post_logout_redirect_uri=${LOGOUT_REDIRECT_URL}`;
  revokeToken();
  removeTokens();
  window.location.assign(logoutUrl);
};

const Home = () => {
  const params = queryString.parse(window.location.search);

  const [accessToken, setAccessToken] = useState(
    localStorage.getItem(ACCESS_TOKEN_STORAGE_KEY)
  );
  const storedState = localStorage.getItem(STATE_STORAGE_KEY);
  const stateMatch =
    "code" in params && "state" in params && params.state === storedState;

  const changeToken = async () => {
    await fetchTokens();
    setAccessToken(localStorage.getItem(ACCESS_TOKEN_STORAGE_KEY));
  };

  useEffect(() => {
    if (stateMatch) {
      changeToken();
      localStorage.removeItem(STATE_STORAGE_KEY);
    }
  }, [stateMatch]);

  if (!accessToken) {
    if (stateMatch) {
      // Wait for the code to be exchanged for a token
      return <CircularProgress />;
    } else {
      return (
        <div className="App">
          <header className="App-header">
            <p>Welcome to example-client-react</p>
            <Button onClick={startAuthentication}>Click here to login</Button>
          </header>
        </div>
      );
    }
  }

  return (
    <div className="App">
      <header className="App-header">
        <p>Connected with token starting with {accessToken.substring(0, 15)}</p>
        <Button onClick={logout}>Click here to logout</Button>
      </header>
    </div>
  );
};

export default Home;
