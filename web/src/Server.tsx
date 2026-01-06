import { AUTH_TYPE } from "./AppProvider";

/**
 * 模拟了Server API
 */

import { v4 as uuidv4 } from 'uuid';

function generateToken() {
  return uuidv4();
}

const fakeAppProvider = {
    isAuthenticated: false,
    token: '',
    account: '',
    signin(params: {account: string, password: string, userType: AUTH_TYPE, code: string}, callback: Function) {
      fakeAppProvider.isAuthenticated = true;
      fakeAppProvider.account = params.account;
      fakeAppProvider.token = generateToken();
      setTimeout(callback({userType:params.userType, account: params.account}, fakeAppProvider.token), 100); // fake async
    },
    signout(_: {account: string, userType: AUTH_TYPE}, callback: Function) {
      fakeAppProvider.isAuthenticated = false;
      setTimeout(callback(), 100);
    },
    checkToken(userType: AUTH_TYPE, _: string, callback: Function) {
      // if (fakeAppProvider.token != '' && fakeAppProvider.token == token) {
        fakeAppProvider.isAuthenticated = true;
        setTimeout(callback({userType:userType, account: '123'}), 100);
      // } else {
      //   fakeAppProvider.token = '';
      //   fakeAppProvider.isAuthenticated = false;
      //   setTimeout(callback({userType:userType}), 100);
      // }
    }
  };

  export { fakeAppProvider };
