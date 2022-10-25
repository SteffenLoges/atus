interface IAuthResponse {
  refreshToken: string;
}

interface IJWTContent {
  uid: string;
  name: string;
  exp: number;
}
