export const ACCESS_TOKEN_STORAGE_KEY = 'ARKHN_ACCESS_TOKEN';
export const REFRESH_TOKEN_STORAGE_KEY = 'ARKHN_REFRESH_TOKEN';
export const ID_TOKEN_STORAGE_KEY = 'ARKHN_ID_TOKEN';
export const TOKEN_DATA_STORAGE_KEY = 'ARKHN_TOKEN_DATA';
export const STATE_STORAGE_KEY = 'ARKHN_AUTH_STATE';
export const {
  REACT_APP_CLIENT_ID: CLIENT_ID,
  REACT_APP_CLIENT_SECRET: CLIENT_SECRET,
  REACT_APP_AUTH_URL: AUTH_URL,
  REACT_APP_TOKEN_URL: TOKEN_URL,
  REACT_APP_LOGOUT_URL: LOGOUT_URL,
  REACT_APP_REVOKE_URL: REVOKE_URL,
  REACT_APP_LOGIN_REDIRECT_URL: LOGIN_REDIRECT_URL,
  REACT_APP_LOGOUT_REDIRECT_URL: LOGOUT_REDIRECT_URL
} = process.env;
